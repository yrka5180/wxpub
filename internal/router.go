package internal

import (
	"fmt"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/controller"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/middleware"
	"google.golang.org/grpc"

	smsPb "git.nova.net.cn/nova/notify/sms-xuanwu/pkg/grpcIFace"
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
		repository.NewWXRepository(
			persistence.NewWxRepo(), persistence.NewUserRepo(), persistence.NewMessageRepo(config.KafkaTopics)))
	accessToken = controller.NewAccessTokenController(
		repository.NewAccessTokenRepository(
			persistence.NewAkRepo()))
	user = controller.NewUserController(
		repository.NewUserRepository(
			persistence.NewUserRepo()))
	msg = controller.NewMessageController(
		repository.NewMessageRepository(
			persistence.NewMessageRepo(config.KafkaTopics), persistence.NewUserRepo()))
}

func registerMiddleware() {
	auth = middleware.NewPassportMiddleware(
		repository.NewPassportRepository(
			persistence.NewPassportRepo()))
}

func registerGRPCClient() {
	smsConn, err := grpc.Dial(config.SmsRPCAddr, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("failed to dial captcha grpc server: %v", err))
	}

	smsClient := smsPb.NewSenderClient(smsConn)
	repository.InitDefaultPhoneVerifyRepository(persistence.NewPhoneVerifyRepo(smsClient))
}

func Run() *gin.Engine {
	registerController()
	registerMiddleware()
	registerGRPCClient()
	engine := gin.Default()
	initRouter(engine)
	return engine
}

func initRouter(router *gin.Engine) {
	open := router.Group("/")
	routerWX(open)

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
	wxGroup := router.Group("")
	{
		// wx开放平台接入测试接口
		wxGroup.GET("", wx.GetWXCheckSign)
		// todo: 暂时先用明文传输，后续补充aes加密传输
		// wx开放平台事件接收
		wxGroup.POST("", wx.GetEventXML)
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
