package domain

type Vehicle struct {
	ID          string `json:"id"`
	Make        string `json:"make"`
	Model       string `json:"model"`
	Year        uint   `json:"year"`
	VIN         string `json:"vin"`
	PlateNumber string `json:"plate_number"`
	ClientID    string `json:"client_id"`
}
