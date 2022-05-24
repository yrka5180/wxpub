package utils

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"regexp"
	"strconv"

	"google.golang.org/grpc/metadata"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func GetUUID() (string, error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return u.String(), nil
}

func ShouldGetTraceID(c context.Context) (traceID string) {
	it := c.Value(consts.ContextTraceID)
	if it == nil {
		// log.Warnln("Could not get trace id from context")
		return
	}
	var ok bool
	if traceID, ok = it.(string); !ok {
		log.Errorf("Invalid trace id value in context: %v", it)
	}
	return
}

func Sha1(str string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenVerifySmsCode 生成随机短信验证码
func GenVerifySmsCode() (verifyCodeAnswer string) {
	verifyCodeAnswer = strconv.Itoa(rand.Intn(900000) + 1e5)
	return
}

// VerifyMobilePhoneFormat 手机号格式检验
func VerifyMobilePhoneFormat(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))/d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}

func ToOutGoingContext(c context.Context) (out context.Context) {
	data := make(map[string]string)
	data[consts.HTTPTraceIDHeader] = ShouldGetTraceID(c)

	md := metadata.New(data)
	out = metadata.NewOutgoingContext(c, md)
	return
}

func RemoveStringRepeated(str []string) (ret []string) {
	ret = make([]string, 0, len(str))
	temp := map[string]struct{}{}
	for _, item := range str {
		// 如果字典中找不到元素，ok=false，!ok为true，就往切片中append元素。
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			ret = append(ret, item)
		}
	}
	return ret
}
