package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/extra"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/g"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/tasks"

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
	tasks.ConsumerTask(globalCtx)

	engine := internal.Run()
	srv := &http.Server{
		Addr:    config.ListenAddr,
		Handler: engine,
	}
	startServer(srv)
	gracefulShutdown(srv)
}

func startServer(srv *http.Server) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("http server listen err: %+v", err)
		}
	}()
}

func gracefulShutdown(srv *http.Server) {
	code := 1
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
EXIT:
	for {
		sig := <-quit
		log.Infoln("received signal:", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			code = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
		break
	}
	log.Infoln("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %+v", err)
	}

	log.Infoln("Server exiting")
	go g.Wait()
	waitWithCtx(g.QuitC)
	os.Exit(code)
}

// waitWithCtx returns when timeout or when all goroutine is quited
func waitWithCtx(c chan struct{}) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	select {
	case <-c:
		log.Info("server exiting")
		return
	case <-ticker.C:
		globalCancel()
		log.Fatal("server timeout to force to shutdown")
	}
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
	err := persistence.NewRepositories(dbConf, config.RedisAddresses, config.SmsRPCAddr, config.CaptchaRPCAddr, debugMode)
	if err != nil {
		panic(err)
	}
	// repository init
	repository.NewWXRepository(
		persistence.DefaultWxRepo(), persistence.DefaultUserRepo(), persistence.DefaultMessageRepo())
	repository.NewAccessTokenRepository(
		persistence.DefaultAkRepo())
	repository.NewUserRepository(
		persistence.DefaultUserRepo(), persistence.DefaultPhoneVerifyRepo())
	repository.NewMessageRepository(
		persistence.DefaultMessageRepo(), persistence.DefaultUserRepo())
}
