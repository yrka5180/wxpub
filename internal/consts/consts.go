package consts

import "github.com/gin-gonic/gin"

type ServerMode string

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

const (
	ServerModeDebug   ServerMode = gin.DebugMode
	ServerModeRelease ServerMode = gin.ReleaseMode

	HTTPTraceIDHeader  = "x-nova-trace-id"
	HTTPTimeoutHeader  = "x-nova-timeout"
	DefaultHTTPTimeOut = 60

	// GinContextContext 存在 gin context 中的标准库 context 实例的 key
	GinContextContext = "context"
	ContextTraceID    = contextKey(HTTPTraceIDHeader)

	Authorization          = "Authorization"
	InternalAPITokenHeader = "x-auth-token"

	DefaultPage     = 1
	DefaultPageSize = 20
	MaxLimitSize    = 100 // 最大只能查询100，且默认为100

	// Token wx 公众号token
	Token = "nova"

	SmsSender = "nova-wxpublic-proxy"
)

const (
	Module               = "git.nova.net.cn/nova/misc/wx-public/proxy"
	DLockPrefix          = "__dlock-"
	RedisKeyAccessToken  = Module + "-access_token"
	RedisLockAccessToken = DLockPrefix + RedisKeyAccessToken
	RedisKeyMsgID        = Module + "-msg_id-"
	RedisKeyAuthN        = Module + "-authN_"

	RedisKeyVerifyCodeSmsID = "sms"
	RedisKeyPrefixChallenge = Module + "-challenge_"
	RedisKeyPrefixSms       = Module + "-sms_"
)

const (
	RedisMsgIDTTL             = 30
	RedisAuthTTL              = 300
	VerifyCodeSmsChallengeTTL = 1800 // 验证短信时候期限设置为30分钟，期间可以重发短信
)

const (
	Credential = "client_credential"
)

const (
	SubscribeEvent             = "subscribe"
	UnsubscribeEvent           = "unsubscribe"
	TEMPLATESENDJOBFINISHEvent = "TEMPLATESENDJOBFINISH"
)

const (
	SubscribeRespContent             = "您好！欢迎您关注【南凌科技】，南凌科技NOVAnet以信息网络服务，构建企业核心竞争力。该公众号用于设备告警信息推送，若贵客户与我司有相关业务联系需推送业务设备告警信息，请点击绑定信息链接："
	UnSubscribeRespContent           = ""
	TEMPLATESENDJOBFINISHRespContent = ""
)

const (
	TemplateSendSuccessStatus   = "success"
	TemplateSendUserBlockStatus = "failed:user block"
	TemplateSendFailedStatus    = "failed: system failed"
)

const (
	MaxRetryCount = 3        // 消息最大失败重试次数，实际调用接口次数为3*3(http client repeated count)=9
	MaxExpireTime = 3 * 3600 // 消息最大过期时间
)

const (
	SendSuccess = iota + 1
	SendRetry
	SendFailure
)

const (
	SendMaxExpireFailureCause = "超过最大重试时间，请人工确认该消息!"
	SendMaxRetryFailureCause  = "超过最大重试次数，请人工确认该消息!"
)
