package application

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
)

type templateApp struct {
	template repository.TemplateRepository
}

// templateApp implements the TemplateInterface
var _ TemplateInterface = &templateApp{}

type TemplateInterface interface {
	ListTemplate(ctx context.Context, param entity.ListTemplateReq) (entity.ListTemplateResp, error)
}

func (t *templateApp) ListTemplate(ctx context.Context, param entity.ListTemplateReq) (entity.ListTemplateResp, error) {
	return t.template.ListTemplate(ctx, param)
}
