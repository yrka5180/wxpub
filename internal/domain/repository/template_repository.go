package repository

import (
	"context"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/infrastructure/persistence"
)

type TemplateRepositoryInterface interface {
	ListTemplateFromRequest(ctx context.Context, param entity.ListTemplateReq) (entity.ListTemplateResp, error)
}

type TemplateRepository struct {
	template *persistence.TemplateRepo
}

// persistence.TemplateRepo implements the TemplateRepositoryInterface
var _ TemplateRepositoryInterface = &persistence.TemplateRepo{}

func NewTemplateRepository(template *persistence.TemplateRepo) *TemplateRepository {
	return &TemplateRepository{
		template: template,
	}
}

func (t *TemplateRepository) ListTemplate(ctx context.Context, param entity.ListTemplateReq) (entity.ListTemplateResp, error) {
	return t.template.ListTemplateFromRequest(ctx, param)
}
