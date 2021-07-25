package graph

import (
	"github.com/sy-software/minerva-owl/internal/handlers"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	OrgHandler handlers.OrganizationGraphqlHandler
	UsrHandler handlers.UserGraphqlHandler
}
