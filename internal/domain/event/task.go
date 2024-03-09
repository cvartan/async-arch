package event

// Данные задачи для бизнес-события
type TaskEventData struct {
	Uuid             string `json:"uuid"`
	Description      string `json:"description"`
	AssignedUserUuid string `json:"assignedUserUuid"`
}

// Данные задачи для CUD-события
type TaskStreamData struct {
	TaskEventData
	State string `json:"state"`
}
