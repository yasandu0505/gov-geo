package models

type Ministry struct {
	ID                int
	Name              string
	Google_map_script string
}

type Department struct {
	ID                int
	Name              string
	Google_map_script string
	MinistryID        int
}
