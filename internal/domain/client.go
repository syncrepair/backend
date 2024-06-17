package domain

type Client struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	CompanyID   string `json:"company_id"`
}
