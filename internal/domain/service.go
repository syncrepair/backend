package domain

type Service struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CompanyID   string  `json:"company_id"`
}
