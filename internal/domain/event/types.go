package event

type EventType string

const (
	// Называем английским названием бизнес-события
	TASK_BE_TASK_CREATED   EventType = "TASK_CREATED"
	TASK_BE_TASK_ASSIGNED  EventType = "TASK_ASSIGNED"
	TASK_BE_TASK_COMPLETED EventType = "TASK_COMPLETED"
	ACC_BE_DEBITING        EventType = "DEBITING"
	ACC_BE_VALUE           EventType = "VALUE"
	ACC_BE_PAYOFF          EventType = "PAYOFF"

	// Называем по шаблону ДОМЕН.СУЩНОСТЬ.ВЫПОЛНЕННОЕ ДЕЙСТВИЕ
	AUTH_CUD_USER_CREATED EventType = "AUTH.USER.CREATED"
	TASK_CUD_TASK_CREATED EventType = "TASK.TASK.CREATED"
	TASK_CUD_TASK_UPDATED EventType = "TASK.TASK.UPDATED"
	ACC_CUD_TASK_PRICED   EventType = "ACC.TASK.UPDATED"
	ACC_CUD_TRX_CREATED   EventType = "ACC.TRANSACTION.CREATED"
)
