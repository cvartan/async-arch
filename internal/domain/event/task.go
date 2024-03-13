package event

// Данные задачи для бизнес-события
type TaskEventData struct {
	Uuid string `json:"uuid"`
	// Расширили модель двумя новыми атрибутами
	Title            string `json:"title"`
	JiraId           string `json:"jira-id"`
	Description      string `json:"description"`
	AssignedUserUuid string `json:"assignedUserUuid"`
}

// Данные задачи для CUD-события
type TaskStreamData struct {
	TaskEventData
	State string `json:"state"`
}
