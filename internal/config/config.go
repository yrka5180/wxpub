package config

import (
	"public-platform-manager/internal/consts"
	"strings"

	config "git.nova.net.cn/nova/go-common/config2"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel   logrus.Level
	ListenAddr = config.DefaultString("listen_addr", ":80")
	// SMode 服务端运行状态，已知影响 gin 框架日志输出等级
	SMode = consts.ServerMode(config.DefaultString("server_mode", string(consts.ServerModeRelease)))

	// DBHost        string
	// DBUser        string
	// DBPassword    string
	// DBName        = config.DefaultString("db_name", "oplog_mgr")
	// DBMaxIdleConn = config.DefaultInt("max_db_idle_conn", 1000)
	// DBMaxOpenConn = config.DefaultInt("max_db_open_conn", 1000)

	// InternalAPISecret string

	WXBaseURL         = config.DefaultString("wx_base_url", "https://api.weixin.qq.com")
	WXAccessTokenURL  = config.DefaultString("wx_access_token_url", WXBaseURL+"/cgi-bin/token")
	WXTemplateListURL = config.DefaultString("wx_template_list_url", WXBaseURL+"/cgi-bin/template/get_all_private_template")
	RedisAddresses    []string
	AppID             = config.MustString("app_id")
	AppSecret         = config.MustString("app_secret")
)

func Init() {
	// DBHost = config.MustString("db_host")
	// DBUser = config.MustString("db_user")
	// DBPassword = config.MustString("db_password")

	if SMode == consts.ServerModeDebug {
		LogLevel = logrus.DebugLevel
	} else {
		LogLevel = logrus.InfoLevel
	}

	// InternalAPISecret = config.MustString("internal_api_secret")
	RedisAddresses = strings.Split(config.MustString("redis_addresses"), ",")
}
