package model

type PaginationParams struct {
	Page  int
	Limit int
}

func (p *PaginationParams) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		p.Limit = 10
	}

	if p.Limit > 100 {
		p.Limit = 100
	}
}

func (p *PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Data       []SubscriptionDB `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}
