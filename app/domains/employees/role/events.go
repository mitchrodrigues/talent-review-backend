package role

import "github.com/google/uuid"

type Created struct {
	ID             uuid.UUID    `json:"id"`
	OrganizationID uuid.UUID    `json:"organizationID"`
	Title          string       `json:"title"`
	Level          int          `josn:"level"`
	Track          EmployeeType `json:"track"`
}

type TitleUpdated struct {
	Title string `json:"title"`
}

type LevelUpdated struct {
	Level int `json:"level"`
}

type TrackUpdated struct {
	Track EmployeeType `json:"track"`
}
