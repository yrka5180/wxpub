package utils

import (
	"context"
	"crypto/sha1"
	"encoding/hex"

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
		log.Errorln("Could not get trace id from context")
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
