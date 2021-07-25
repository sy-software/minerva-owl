package handlers

import (
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/utils"
)

type OrganizationGraphqlHandler struct {
	service service.OrganizationService
}

func NewOrgGraphqlHandler(service service.OrganizationService) *OrganizationGraphqlHandler {
	return &OrganizationGraphqlHandler{
		service: service,
	}
}

func (handler *OrganizationGraphqlHandler) Create(name string, description string, logo *string) (*model.Organization, error) {
	validatedLogo := utils.CoalesceStr(logo, "")

	org, err := handler.service.Create(name, description, validatedLogo)

	if err != nil {
		return nil, err
	}

	graphModel := model.Organization{
		ID:          org.Id,
		Name:        org.Name,
		Description: org.Description,
		Logo:        &org.Logo,
	}

	return &graphModel, err
}

func (handler *OrganizationGraphqlHandler) Update(id string, name *string, description *string, logo *string) (*model.Organization, error) {
	// TODO: Avoid get to save but for now is required to support PATCH
	current, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	new := domain.Organization{
		Id:          id,
		Name:        utils.CoalesceStr(name, current.Name),
		Description: utils.CoalesceStr(description, current.Description),
		Logo:        utils.CoalesceStr(logo, current.Logo),
	}

	output, err := handler.service.Update(new)

	if err != nil {
		return nil, err
	}

	return orgToGraphQLModel(&output), nil
}

func (handler *OrganizationGraphqlHandler) Query(page *int, pageSize *int) ([]*model.Organization, error) {
	all, err := handler.service.List(page, pageSize)

	if err != nil {
		return []*model.Organization{}, err
	}

	out := make([]*model.Organization, len(all))

	for index, elem := range all {
		out[index] = orgToGraphQLModel(&elem)
	}

	return out, nil
}

func (handler *OrganizationGraphqlHandler) QueryById(id string) (*model.Organization, error) {
	out, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	return orgToGraphQLModel(&out), nil
}

func (handler *OrganizationGraphqlHandler) Delete(id string) (*model.Organization, error) {
	out, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	err = handler.service.Delete(id, false)

	if err != nil {
		return nil, err
	}

	return orgToGraphQLModel(&out), nil
}

func orgToGraphQLModel(source *domain.Organization) *model.Organization {
	return &model.Organization{
		ID:          source.Id,
		Name:        source.Name,
		Description: source.Description,
		Logo:        &source.Logo,
	}
}
