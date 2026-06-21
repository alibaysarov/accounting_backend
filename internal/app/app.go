package app

import (
	"acc_backend/internal/app/routers"
	"acc_backend/internal/container"
	"acc_backend/internal/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Application struct {
	appConfig *settings.AppConfig
}

func NewApp() *Application {
	cfg := &settings.AppConfig{}
	return &Application{cfg}
}

func (app *Application) Run() error {

	err := app.appConfig.Init()
	if err != nil {
		return err
	}
	gin.SetMode(app.appConfig.GinMode)

	db, err := app.dbInit()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	utils := &container.Utils{
		Config: app.appConfig,
	}
	c := container.NewContainer(db, utils)

	router := routers.NewRouter(c)
	if app.appConfig.Port == "" {
		return fmt.Errorf("No port for web server no app.appConfig.Port is given!")
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.appConfig.Port),
		Handler: router.Handler(),
	}

	go app.startServer(server)

	if app.appConfig.GinMode == "debug" {
		go startPprofServer("localhost:6060")
	}

	<-ctx.Done() // ждём сигнал
	fmt.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	app.shutDownServer(ctx, server)

	return nil
}

func (app *Application) startServer(server *http.Server) {

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	log.Println("server started", app.appConfig.Port)

}

func (app *Application) restrictConn(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	// Ограничиваем пул так, чтобы оставался запас под другие нужды
	// (миграции, админ-подключения, k6 cleanup и т.д.)
	sqlDB.SetMaxOpenConns(80)                  // максимум одновременных соединений
	sqlDB.SetMaxIdleConns(20)                  // сколько держим "прогретыми" в простое
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // пересоздавать соединения раз в N времени
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	return nil
}

func (app *Application) dbInit() (*gorm.DB, error) {
	url := app.appConfig.DbUrl

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = app.restrictConn(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (a Application) shutDownServer(ctx context.Context, server *http.Server) {

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
