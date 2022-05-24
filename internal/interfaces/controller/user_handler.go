package controller

import (
	"strconv"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/middleware"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type User struct {
	user application.UserInterface
}

func NewUserController(user application.UserInterface) *User {
	return &User{
		user: user,
	}
}

func (u *User) GenCaptcha(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)

	width := c.DefaultQuery("width", strconv.Itoa(consts.CaptchaDefaultWidth))
	w, err := strconv.ParseInt(width, 10, 64)
	if err != nil {
		log.Errorf("%s get width failed,err: %+v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, errors.GetErrorMessage(errors.CodeInternalServerError))
		return
	}
	height := c.DefaultQuery("height", strconv.Itoa(consts.CaptchaDefaultHeight))
	h, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		log.Errorf("%s get height failed,err: %+v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, errors.GetErrorMessage(errors.CodeInternalServerError))
		return
	}

	captchaID, captchaBase64Value, err := u.user.GenCaptcha(ctx, int32(w), int32(h))
	if err != nil {
		log.Errorf("%s get captchaID and captchaBase64Value failed,err: %+v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, "get captchaID and captchaBase64Value failed")
		return
	}

	CaptchaResp := entity.CaptchaResp{
		CaptchaID:          captchaID,
		CaptchaBase64Value: captchaBase64Value,
	}
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, CaptchaResp)
}

func (u *User) SendSms(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)

	var req entity.SendSmsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("SendSms ShouldBindJSON error: %+v, traceID:%s", err, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, errors.GetErrorMessage(errors.CodeInvalidParams))
		return
	}

	if utils.VerifyMobilePhoneFormat(req.Phone) {
		log.Errorf("invaild phone number: %s, traceID:%s", req.Phone, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, errors.GetErrorMessage(errors.CodeInvalidParams))
		return
	}

	ok, err := u.user.VerifyCaptcha(ctx, req.CaptchaID, req.CaptchaAnswer)
	if err != nil {
		log.Errorf("VerifyCaptcha failed, traceID:%s, err:%+v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, errors.GetErrorMessage(errors.CodeInternalServerError))
		return
	}
	if !ok {
		log.Errorf("wrong captcha answer: %s, traceID:%s", req.CaptchaAnswer, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeForbidden, "wrong captcha answer")
		return
	}

	err = u.user.SendSms(ctx, req)
	if err != nil {
		log.Errorf("SendSms failed, traceID:%s, err:%+v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, errors.GetErrorMessage(errors.CodeInternalServerError))
		return
	}

	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, nil)
}

func (u *User) VerifyAndUpdatePhone(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)

	var req entity.VerifyCodeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Errorf("SendSms ShouldBindJSON error: %+v, traceID:%s", err, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, errors.GetErrorMessage(errors.CodeInvalidParams))
		return
	}

	if utils.VerifyMobilePhoneFormat(req.Phone) {
		log.Errorf("invaild phone number: %s, traceID:%s", req.Phone, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, errors.GetErrorMessage(errors.CodeInvalidParams))
		return
	}

	ok, isExpire, err := u.user.VerifySmsCode(ctx, req)
	if err != nil {
		log.Errorf("user VerifySmsCode failed, traceID:%s, err:%+v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, errors.GetErrorMessage(errors.CodeInternalServerError))
		return
	}
	if !ok {
		log.Errorf("verify code is not correct, code: %s, traceID: %s", req.VerifyCode, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeForbidden, errors.GetErrorMessage(errors.CodeForbidden))
		return
	}
	if isExpire {
		log.Errorf("sms code is expired, code: %s, traceID: %s", req.VerifyCode, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeTokenExpire, errors.GetErrorMessage(errors.CodeTokenExpire))
		return
	}

	user, err := u.user.GetUserByOpenID(ctx, req.OpenID)
	if err != nil {
		log.Errorf("VerifyAndUpdatePhone get user by open_id error: %+v, traceID: %s", err, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, err.Error())
		return
	}

	user.Phone = req.Phone
	user.Name = req.Name
	err = u.user.SaveUser(ctx, user, true)
	if err != nil {
		log.Errorf("VerifyAndUpdatePhone update user error: %+v, traceID: %s", err, traceID)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, err.Error())
		return
	}

	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, nil)
}
