package models

type Ministry struct {
	ID                int    `json:"id,omitempty"`
	Name              string `json:"name"`
	Google_map_script string `json:"google_map_script"`
}

type Department struct {
	ID                int    `json:"id,omitempty"`
	Name              string `json:"name"`
	Google_map_script string `json:"google_map_script"`
	MinistryID        int    `json:"ministry_id"`
}

type MinistryWithDepartments struct {
	Ministry
	Departments []Department
}
