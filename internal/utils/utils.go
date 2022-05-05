package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"public-platform-manager/internal/consts"
	"regexp"
	"strings"
	"time"

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
		log.Errorln("Could not get trace id from context")
		return
	}
	var ok bool
	if traceID, ok = it.(string); !ok {
		log.Errorf("Invalid trace id value in context: %v", it)
	}
	return
}

// GetAppID return app_id or empty if not exists
func GetAppID(c context.Context) (appID string) {
	it := c.Value(consts.ContextAppID)
	if it == nil {
		log.Errorf("could not get app_id from context")
		return
	}
	var ok bool
	if appID, ok = it.(string); !ok {
		log.Errorf("Invalid app_id value in context: %v", it)
	}
	return
}

func CheckIP(ip string) bool {
	addr := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}

func GetMD5String(b []byte) (result string) {
	res := md5.Sum(b)
	result = hex.EncodeToString(res[:])
	return
}

func GetAccessKey() (string, error) {
	UUID, err := GetUUID()
	if err != nil {
		return "", err
	}
	return strings.Replace(UUID, "-", "", -1), nil
}

func GetSecretAccessKey(accessKey string) string {
	secretAccessString := fmt.Sprintf("%s@%d", accessKey, time.Now().UnixNano())
	secretAccessKey := GetMD5String([]byte(secretAccessString))
	return secretAccessKey
}

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
