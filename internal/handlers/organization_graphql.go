package handlers

import (
	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/service"
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
	validatedLogo := nilCoalescing(logo, "")

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
		Name:        nilCoalescing(name, current.Name),
		Description: nilCoalescing(description, current.Description),
		Logo:        nilCoalescing(logo, current.Logo),
	}

	output, err := handler.service.Update(new)

	if err != nil {
		return nil, err
	}

	return domainToGraphQLModel(&output), nil
}

func (handler *OrganizationGraphqlHandler) Query() ([]*model.Organization, error) {
	all, err := handler.service.All()

	if err != nil {
		return []*model.Organization{}, err
	}

	out := make([]*model.Organization, len(all))

	for index, elem := range all {
		out[index] = domainToGraphQLModel(&elem)
	}

	return out, nil
}

func (handler *OrganizationGraphqlHandler) QueryById(id string) (*model.Organization, error) {
	out, err := handler.service.Get(id)

	if err != nil {
		return nil, err
	}

	return domainToGraphQLModel(&out), nil
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

	return domainToGraphQLModel(&out), nil
}

func domainToGraphQLModel(source *domain.Organization) *model.Organization {
	return &model.Organization{
		ID:          source.Id,
		Name:        source.Name,
		Description: source.Description,
		Logo:        &source.Logo,
	}
}

func nilCoalescing(value *string, fallback string) string {
	if value != nil {
		return *value
	} else {
		return fallback
	}
}
