import http from 'k6/http';
import { check, sleep,randomSeed } from 'k6';
import { Rate, Trend } from 'k6/metrics';
import { SharedArray } from 'k6/data';
import exec from 'k6/execution';

import sql from 'k6/x/sql';
import driver from 'k6/x/sql/driver/postgres';

const registerErrors = new Rate('register_errors');
const registerDuration = new Trend('register_duration', true);

const users = new SharedArray('users', function () {
  const arr = [];
  for (let i = 0; i < 10000; i++) {
    arr.push({
      full_name: `Test User ${i}`,
      email: `loadtest_${i}@example.com`,
      password: 'P@ssw0rd123!',
    });
  }
  return arr;
});

export const options = {
  scenarios: {
    register_load: {
      executor: 'ramping-arrival-rate',
      startRate: 10,           // стартовый RPS
      timeUnit: '1s',
      preAllocatedVUs: 100,    // сколько VU держим прогретыми заранее
      maxVUs: 1000,            // потолок VU, если ответы начнут тормозить
      stages: [
        { duration: '30s', target: 100 },  // разгон до 100 rps
        { duration: '30s', target: 500 },
        { duration: '1m', target: 500 },   // держим 500 rps
        { duration: '20s', target: 0 },    // спад
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
    register_errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080/api/v1';

export default function () {
  // Глобальный уникальный номер итерации по всему тесту
  const globalIdx = exec.scenario.iterationInTest;
  const user = users[globalIdx % users.length];

  const payload = JSON.stringify({
    full_name: user.full_name,
    email: `${globalIdx}_${user.email}`, // гарантия уникальности email
    password: user.password,
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(`${BASE_URL}/auth/register`, payload, params);

  registerDuration.add(res.timings.duration);

  const ok = check(res, {
    'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
    'response has body': (r) => r.body && r.body.length > 0,
  });

  if (!ok) {
    console.log(`FAILED: status=${res.status} body=${res.body}`);
  }

  registerErrors.add(!ok);
}


export function teardown() {
  const DB_URL = __ENV.DB_URL;

  if (!DB_URL) {
    console.log('DB_URL не задан — пропускаю очистку. Удали тестовых пользователей вручную.');
    return;
  }

  const db = sql.open(driver, `${DB_URL}?sslmode=disable`);

  try {
    const result = db.exec(`DELETE FROM users WHERE email LIKE '%loadtest%'`);

    console.log(`Удалено строк: ${result.rowsAffected()}`);
  } catch (e) {
    console.log(`Ошибка очистки БД: ${e}`);
  } finally {
    db.close();
  }
}