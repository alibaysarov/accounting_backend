from fastapi import FastAPI
from contextlib import asynccontextmanager

from app.routers.api import api_router

import asyncio

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("App is starting...")

    # например: подключение к БД
    # app.state.db = await create_db_pool()

    yield  # приложение работает

    # Shutdown (graceful shutdown)
    print("App is shutting down...")

    # закрытие ресурсов
    # await app.state.db.close()

    await asyncio.sleep(0.5)  # имитация завершения фоновых задач

app = FastAPI()

app.include_router(api_router)

@app.get("/")
async def root():
    return {"status": "ok"}