package domain

type Organization struct {
	Id          string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Logo        string `bson:"logo,omitempty" json:"logo,omitempty"`
}

func (org *Organization) GetId() string {
	return org.Id
}
