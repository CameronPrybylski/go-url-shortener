package model

type URL struct {

	ID int `json:"id"`
	ShortenedCode string `json:"short_code"`
	OriginalURL string `json:"original_url"`
	CreatedAt string `json:"created_at"`
	VisitCount int `json:"visit_count"`
	LastAccessed string `json:"last_accessed"`
}
