package domain

import (
	"time"
)

const USER_COL_NAME = "users"

// An user authenticated into minerva platform using one OAuth2 provider
type User struct {
	Id         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Username   string    `bson:"username,omitempty" json:"username,omitempty"`
	Name       string    `bson:"name,omitempty" json:"name,omitempty"`
	Picture    string    `bson:"picture,omitempty" json:"picture,omitempty"`
	Role       string    `bson:"role,omitempty" json:"role,omitempty"`
	Provider   string    `bson:"provider,omitempty" json:"provider,omitempty"`
	TokenID    string    `bson:"tokenID,omitempty" json:"tokenID,omitempty"`
	CreateDate time.Time `bson:"createDate,omitempty" json:"createDate,omitempty"`
	UpdateDate time.Time `bson:"updateDate,omitempty" json:"updateDate,omitempty"`
	Status     string    `bson:"status,omitempty" json:"status,omitempty"`
}
