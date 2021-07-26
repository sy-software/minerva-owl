package domain

import (
	"time"
)

const USER_COL_NAME = "users"

// An user authenticated into minerva platform using one OAuth2 provider
type User struct {
	Id string `bson:"_id,omitempty" json:"id,omitempty"`
	// User screen name, used for login
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	// User real name
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// Optional url of the user display image
	Picture string `bson:"picture,omitempty" json:"picture,omitempty"`
	// For RBAC operations
	Role string `bson:"role,omitempty" json:"role,omitempty"`
	// The OAuth2 provider used by this user
	Provider string `bson:"provider,omitempty" json:"provider,omitempty"`
	// The identifier connection this user with the OAuth provider
	TokenID    string    `bson:"tokenID,omitempty" json:"tokenID,omitempty"`
	CreateDate time.Time `bson:"createDate,omitempty" json:"createDate,omitempty"`
	UpdateDate time.Time `bson:"updateDate,omitempty" json:"updateDate,omitempty"`
	// Can be used to control the user status inside the platform
	Status string `bson:"status,omitempty" json:"status,omitempty"`
}
