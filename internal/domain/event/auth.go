package event

// UserEventData - данные пользователя в событии
type UserEventData struct {
	Uuid  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
