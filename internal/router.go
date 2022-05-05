package internal

import (
	"public-platform-manager/internal/domain/repository"
	"public-platform-manager/internal/infrastructure/persistence"
	"public-platform-manager/internal/interfaces/controller"
	"public-platform-manager/internal/interfaces/middleware"

	"github.com/gin-gonic/gin"
)

var (
	wxChecker   *controller.WXCheckSignature
	accessToken *controller.AccessToken
	user        *controller.User
)

func registerController() {
	wxChecker = controller.NewWXCheckSignController(repository.NewWXCheckSignRepository())
	accessToken = controller.NewAccessTokenController(
		repository.NewAccessTokenRepository(
			persistence.NewAkRepo()))
	user = controller.NewUserController(repository.NewUserRepository())
}

func Run() *gin.Engine {
	registerController()
	engine := gin.Default()
	initRouter(engine)
	return engine
}

func initRouter(router *gin.Engine) {
	// wx开放平台接入测试接口
	router.GET("/", wxChecker.GetWXCheckSign)

	router.Use(middleware.NovaContext)
	// todo:鉴权认证
	interval := router.Group("/interval/v1")
	akGroup := interval.Group("/access_token")
	{
		// 获取wx access token
		akGroup.GET("", accessToken.GetAccessToken)
		// 刷新wx access token
		akGroup.GET("/fresh", accessToken.FreshAccessToken)
	}
	// 获取wx user info
	userGroup := interval.Group("/user")
	{
		userGroup.GET("", user.ListUser)
		userGroup.GET("/:id", user.GetUser)
	}
}
