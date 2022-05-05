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
	template    *controller.Template
)

func registerController() {
	wxChecker = controller.NewWXCheckSignController(repository.NewWXCheckSignRepository())
	accessToken = controller.NewAccessTokenController(
		repository.NewAccessTokenRepository(
			persistence.NewAkRepo()))
	user = controller.NewUserController(repository.NewUserRepository())
	template = controller.NewTemplateController(
		repository.NewTemplateRepository(
			persistence.NewTemplateRepo()))
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

	// access token
	routerAccessToken(interval)

	// 获取wx user info
	routerUser(interval)

	// 模板管理
	routerTemplate(interval)

	// 事件推送
	routerPush(interval)
}

func routerAccessToken(router *gin.RouterGroup) {
	akGroup := router.Group("/access_token")
	{
		// 获取wx access token
		akGroup.GET("", accessToken.GetAccessToken)
		// 刷新wx access token
		akGroup.GET("/fresh", accessToken.FreshAccessToken)
	}
}

func routerUser(router *gin.RouterGroup) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("", user.ListUser)
		userGroup.GET("/:id", user.GetUser)
	}
}

func routerTemplate(router *gin.RouterGroup) {
	templateGroup := router.Group("/template")
	{
		templateGroup.GET("", template.ListTemplate)
	}
}

func routerPush(router *gin.RouterGroup) {
	pushGroup := router.Group("/push")
	{
		// 告警事件推送
		alertSubGroup := pushGroup.Group("/alert")
		{
			alertSubGroup.GET("")
		}
	}
}
