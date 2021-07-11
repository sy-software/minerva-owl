package domain

// Organization is the main element in the data model
//
// An organization witholds: Areas, Teams, Users, Software components
// and all other entities for a single client inside minerva
type Organization struct {
	Id          string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Logo        string `bson:"logo,omitempty" json:"logo,omitempty"`
}

// Area represents a subdivision of an organization such as: Engineering, Design, etc.
type Area struct {
	Id           string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string `bson:"name,omitempty" json:"name,omitempty"`
	Description  string `bson:"description,omitempty" json:"description,omitempty"`
	Organization string `bson:"organization,omitempty" json:"organization,omitempty"`
	Color        string `bson:"color,omitempty" json:"color,omitempty"`
	Icon         string `bson:"icon,omitempty" json:"icon,omitempty"`
}

// Team represents a unit of people working on a commong goal
//
// A Team is managed by a "Leader". And it's composed of multiple
// team members
// TODO: Should a team belong to a single Area?
type Team struct {
	Id           string   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string   `bson:"name,omitempty" json:"name,omitempty"`
	Description  string   `bson:"description,omitempty" json:"description,omitempty"`
	Organization string   `bson:"organization,omitempty" json:"organization,omitempty"`
	Leader       string   `bson:"leader,omitempty" json:"leader,omitempty"`
	Color        string   `bson:"color,omitempty" json:"color,omitempty"`
	Icon         string   `bson:"icon,omitempty" json:"icon,omitempty"`
	Techs        []string `bson:"techs,omitempty" json:"techs,omitempty"`
}

// Tech is a definition of tools, languages, frameworks, etc. Used within an Organization
type Teach struct {
	Id           string `bson:"_id,omitempty" json:"_id,omitempty"`
	Name         string `bson:"name,omitempty" json:"name,omitempty"`
	Description  string `bson:"description,omitempty" json:"description,omitempty"`
	Organization string `bson:"organization,omitempty" json:"organization,omitempty"`
	Type         string `bson:"type,omitempty" json:"type,omitempty"`
}
