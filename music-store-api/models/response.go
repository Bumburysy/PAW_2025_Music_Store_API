package models

// ErrorResponse reprezentuje odpowiedź błędu
// swagger:model ErrorResponse
type ErrorResponse struct {
	// Komunikat błędu
	Error string `json:"error"`
}

// SuccessResponse reprezentuje odpowiedź sukcesu
// swagger:model SuccessResponse
type SuccessResponse struct {
	// Wiadomość o sukcesie
	Message string `json:"message"`
}
