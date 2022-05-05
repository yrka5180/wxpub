package repository

import (
	"public-platform-manager/internal/utils"
	"sort"
	"strings"
)

type WXCheckSignRepository struct {
}

func NewWXCheckSignRepository() *WXCheckSignRepository {
	return &WXCheckSignRepository{}
}

func (a *WXCheckSignRepository) GetWXCheckSign(signature, timestamp, nonce, token string) bool {
	// 本地计算signature
	si := []string{token, timestamp, nonce}
	// 字典序排序
	sort.Strings(si)
	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for _, v := range si {
		b.WriteString(v)
	}
	return utils.Sha1(b.String()) == signature
}
