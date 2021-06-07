package domain

type Organization struct {
	Id          string
	Name        string
	Description string
	Logo        string
}

func (org *Organization) GetId() string {
	return org.Id
}
