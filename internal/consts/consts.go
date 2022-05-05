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

	ContextAccessKey = contextKey(AccessKeyHeader)
	ContextAppID     = contextKey(AppIDHeader)

	SignTimestampHeader = "x-nova-sign-timestamp"
	SignExpireHeader    = "x-nova-sign-expire"
	SignDebugHeader     = "x-nova-sign-debug"

	AppIDHeader     = "x-nova-app-id"
	AccessKeyHeader = "x-nova-access-key"

	Authorization          = "Authorization"
	InternalAPITokenHeader = "x-auth-token"

	DefaultPage     = 1
	DefaultPageSize = 20
	MaxLimitSize    = 100 // 最大只能查询100，且默认为100

	// Token wx 公众号token
	Token = "nova"
)

const (
	Module               = "public-platform-manager"
	DLockPrefix          = "__dlock-"
	RedisKeyAccessToken  = Module + "-access_token"
	RedisLockAccessToken = DLockPrefix + Module + "-access_token"
)

const (
	Credential = "client_credential"
)
