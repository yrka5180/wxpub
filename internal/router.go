package internal

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/controller"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/middleware"
	"github.com/gin-gonic/gin"
)

var (
	wx          *controller.WX
	accessToken *controller.AccessToken
	user        *controller.User
	msg         *controller.Message
)

var (
	auth *middleware.Passport
)

func registerController() {
	wx = controller.NewWXController(
		repository.DefaultWXRepository())
	accessToken = controller.NewAccessTokenController(
		repository.DefaultAccessTokenRepository())
	user = controller.NewUserController(
		repository.DefaultUserRepository())
	msg = controller.NewMessageController(
		repository.DefaultMessageRepository())
}

func registerMiddleware() {
	auth = middleware.NewPassportMiddleware(
		repository.DefaultPassportRepository())
}

func Run() *gin.Engine {
	registerController()
	registerMiddleware()
	engine := gin.Default()
	initRouter(engine)
	return engine
}

func initRouter(router *gin.Engine) {
	open := router.Group("")
	// wx api
	routerWX(open)
	// user info verify and binding
	routerVerify(open)

	router.Use(middleware.NovaContext)
	interval := router.Group("/interval/v1", auth.VerifyToken)

	// access token
	routerAccessToken(interval)

	// 获取wx user info
	routerUser(interval)

	// 消息推送
	routerMsgPush(interval)
}

func routerWX(router *gin.RouterGroup) {
	wxGroup := router.Group("/")
	{
		// wx开放平台接入测试接口
		wxGroup.GET("", wx.GetWXCheckSign)
		// todo: 暂时先用明文传输，后续补充aes加密传输
		// wx开放平台事件接收
		wxGroup.POST("", wx.GetEventXML)
	}
}

func routerVerify(router *gin.RouterGroup) {
	smsProfileGroup := router.Group("/user")
	{
		smsProfileGroup.GET("/send-sms", user.SendSms)
		smsProfileGroup.POST("/verify-sms", user.VerifyAndUpdatePhone)
		smsProfileGroup.GET("/captcha", user.GenCaptcha)
	}
}

func routerAccessToken(router *gin.RouterGroup) {
	akGroup := router.Group("/access_token")
	{
		// 获取wx access token
		akGroup.GET("", accessToken.GetAccessToken)
		// 刷新wx access token
		// todo:接口限频，微信日调用次数2000次，如果access token缓存值没失效则被视为有效调用（获取ak时ak不存在也会调用），调用次数记录，1分钟1次
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

func routerMsgPush(router *gin.RouterGroup) {
	msgPushGroup := router.Group("/message/push")
	{
		// 模板消息推送
		tmplSubGroup := msgPushGroup.Group("/tmpl")
		{
			tmplSubGroup.POST("", msg.SendTmplMessage)
		}
	}
}
