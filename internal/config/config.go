package config

import (
	"strings"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"

	config "git.nova.net.cn/nova/go-common/config2"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel   logrus.Level
	ListenAddr = config.DefaultString("listen_addr", ":80")
	// SMode 服务端运行状态，已知影响 gin 框架日志输出等级
	SMode = consts.ServerMode(config.DefaultString("server_mode", string(consts.ServerModeRelease)))

	KafkaBrokers []string
	KafkaTopics  []string
	KafkaVersion string
	KafkaTopic   = config.DefaultString("kafka_topic", "public-platform-msg")
	KafkaGroup   = config.DefaultString("kafka_group", "public-platform-msg")

	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        = config.DefaultString("db_name", "pub_platform_mgr")
	DBMaxIdleConn = config.DefaultInt("max_db_idle_conn", 1000)
	DBMaxOpenConn = config.DefaultInt("max_db_open_conn", 1000)

	PassportBaseURL           = config.DefaultString("passport_base_url", "https://passport.nova.net.cn")
	PassportOIDCIntrospectURL = config.DefaultString("passport_oidc_oauth2_introspect", PassportBaseURL+"/apis/v1/oidc/oauth2/introspect")

	SmsRPCAddr           = config.DefaultString("sms_rpc_addr", "sms-xuanwu.common:80")
	SmsContentTemplateCN = config.DefaultString("sms_content_template_cn", "南凌科技验证码：%s。尊敬的用户，您正在绑定手机号，切勿轻易将验证码告知他人！")

	WXBaseURL        = config.DefaultString("wx_base_url", "https://api.weixin.qq.com")
	WXAccessTokenURL = config.DefaultString("wx_access_token_url", WXBaseURL+"/cgi-bin/token")
	WXMsgTmplSendURL = config.DefaultString("wx_msg_tmpl_send_url", WXBaseURL+"/cgi-bin/message/template/send")
	RedisAddresses   []string
	AppID            = config.MustString("app_id")
	AppSecret        = config.MustString("app_secret")
	TmplMsgID        = config.MustString("tmpl_msg_id")
)

func Init() {
	DBHost = config.MustString("db_host")
	DBUser = config.MustString("db_user")
	DBPassword = config.MustString("db_password")

	if SMode == consts.ServerModeDebug {
		LogLevel = logrus.DebugLevel
	} else {
		LogLevel = logrus.InfoLevel
	}

	KafkaBrokers = strings.Split(config.MustString("kafka_brokers"), ",")
	KafkaTopics = strings.Split(KafkaTopic, ",")
	KafkaVersion = config.DefaultString("kafka_version", "1.1.1")

	// InternalAPISecret = config.MustString("internal_api_secret")
	RedisAddresses = strings.Split(config.MustString("redis_addresses"), ",")
}
