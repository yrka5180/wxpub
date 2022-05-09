package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"public-platform-manager/internal"
	"public-platform-manager/internal/config"
	"public-platform-manager/internal/consts"
	"public-platform-manager/internal/g"
	"public-platform-manager/internal/infrastructure/persistence"
	"syscall"
	"time"

	extra "git.nova.net.cn/nova/go-common/logrus-extra"
	log "github.com/sirupsen/logrus"
)

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

func main() {
	config.Init()
	extra.Default(config.LogLevel)
	globalCtx, globalCancel = context.WithCancel(context.Background())
	// init
	InitService()

	engine := internal.Run()
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: engine,
	}
	startServer(srv)
	gracefulShutdown(srv)
	globalCancel()
	g.Wait()
}

func startServer(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("http server listen err: %+v", err)
		}
	}()
}

func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infoln("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Infoln("Server exiting")
}

func InitService() {
	debugMode := config.SMode == consts.ServerModeDebug
	dbConf := persistence.DBConfig{
		DBUser:      config.DBUser,
		DBPassword:  config.DBPassword,
		DBHost:      config.DBHost,
		DBName:      config.DBName,
		MaxIdleConn: config.DBMaxIdleConn,
		MaxOpenConn: config.DBMaxOpenConn,
	}
	err := persistence.NewDBRepositories(dbConf, debugMode)
	if err != nil {
		panic(err)
	}
	err = persistence.NewRedisRepositories(config.RedisAddresses)
	if err != nil {
		panic(err)
	}
}
