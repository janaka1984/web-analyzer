package domain

import "time"

type Analysis struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	URL           string    `gorm:"index;not null" json:"url"`
	HTMLVersion   string    `json:"html_version"`
	Title         string    `json:"title"`
	H1            int       `json:"h1"`
	H2            int       `json:"h2"`
	H3            int       `json:"h3"`
	H4            int       `json:"h4"`
	H5            int       `json:"h5"`
	H6            int       `json:"h6"`
	InternalLinks int       `json:"internal_links"`
	ExternalLinks int       `json:"external_links"`
	BrokenLinks   int       `json:"broken_links"`
	HasLoginForm  bool      `json:"has_login_form"`
	HTTPStatus    int       `json:"http_status"`
	ErrorMessage  string    `json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
}
