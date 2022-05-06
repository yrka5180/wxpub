package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"public-platform-manager/internal/config"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/interfaces/errors"
	"public-platform-manager/internal/interfaces/httputil"
	"public-platform-manager/internal/utils"

	log "github.com/sirupsen/logrus"
)

type TemplateRepo struct {
}

func NewTemplateRepo() *TemplateRepo {
	return &TemplateRepo{}
}

func (a *TemplateRepo) ListTemplateFromRequest(ctx context.Context, param entity.ListTemplateReq) (entity.ListTemplateResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("ListTemplateFromRequest traceID:%s", traceID)
	// 请求wx template
	requestProperty := httputil.GetRequestProperty(http.MethodGet, config.WXTemplateListURL+fmt.Sprintf("?access_token=%s", param.AccessToken),
		nil, make(map[string]string))
	statusCode, body, _, err := httputil.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("request wx template list failed, traceID:%s, error:%v", traceID, err)
		return entity.ListTemplateResp{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("request wx template list failed, statusCode:%d,traceID:%s, error:%v", statusCode, traceID, err)
		return entity.ListTemplateResp{}, err
	}
	var tmpResp entity.ListTemplateResp
	err = json.Unmarshal(body, &tmpResp)
	if err != nil {
		log.Errorf("get wx template list failed by unmarshal, resp:%s, traceID:%s, err:%v", string(body), traceID, err)
		return entity.ListTemplateResp{}, err
	}
	// token过期
	if tmpResp.ErrCode == errors.CodeRIDExpired {
		err = errors.NewCustomError(nil, errors.CodeForbidden, errors.GetErrorMessage(errors.CodeForbidden))
		return entity.ListTemplateResp{}, err
	}
	// 获取失败
	if tmpResp.ErrCode != errors.CodeOK {
		log.Errorf("get wx template list failed,resp:%s,traceID:%s,errMsg:%s", string(body), traceID, tmpResp.ErrMsg)
		return entity.ListTemplateResp{}, fmt.Errorf("get wx template list failed,errMsg:%s", tmpResp.ErrMsg)
	}
	return tmpResp, nil
}
