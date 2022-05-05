package middleware

// func makeSignature(token string, timestamp, nonce string) string {
// 	// 本地计算signature
// 	si := []string{token, timestamp, nonce}
// 	sort.Strings(si)            // 字典序排序
// 	str := strings.Join(si, "") // 组合字符串
// 	s := sha1.New()             // 返回一个新的使用SHA1校验的hash.Hash接口
// 	io.WriteString(s, str)      // WriteString函数将字符串数组str中的内容写入到s中
// 	return fmt.Sprintf("%x", s.Sum(nil))
// }
//
// func ValidateUrl(ctx *gin.Context) {
// 	timestamp := strings.Join(ctx.Request.Form["timestamp"], "")
// 	nonce := strings.Join(ctx.Request.Form["nonce"], "")
// 	signature := strings.Join(ctx.Request.PostForm["signature"], "")
// 	echostr := strings.Join(ctx.Request.Form["echostr"], "")
// 	signatureGen := makeSignature(consts.Token, timestamp, nonce)
//
// 	if signatureGen != signature {
// 		httputil.Abort(ctx, errors.CodeResourcesConflict)
// 	}
// 	fmt.Fprintf(ctx.Writer, echostr) // 原样返回eechostr给微信服务器
// 	return true
// }
