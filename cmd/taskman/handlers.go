// Обработчики http-запросов к сервису упраления задачами

package main

import (
	authmodel "async-arch/internal/domain/auth"
	eventmodel "async-arch/internal/domain/event"
	model "async-arch/internal/domain/taskman"
	authtool "async-arch/internal/lib/auth"
	base "async-arch/internal/lib/base"
	"async-arch/internal/lib/httputils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

// Инициализация обработчиков
func initHandlers() {
	base.App.HandleFunc("POST /api/v1/tasks", authtool.WithAuth(handleCreateTask, nil))
	base.App.HandleFunc(
		"POST /api/v1/tasks/reassign",
		authtool.WithAuth(handleReassignTask,
			[]authmodel.UserRole{
				authmodel.ADMIN,
				authmodel.MANAGER,
			},
		),
	)
	base.App.HandleFunc("POST /api/v1/tasks/{id}/complete", authtool.WithAuth(handleCompleteTask, nil))
	base.App.HandleFunc("GET /api/v1/tasks", authtool.WithAuth(handleGetUserTasks, nil))
}

// Обработчик запроса создания задачи
func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	// Получаем тело запроса
	var taskRq CreateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&taskRq)
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	// Делаем проверку, что в заголовке задачи нет указания на JiraId - то есть остутствуют символы []

	if strings.Contains(taskRq.Title, "[") && strings.Contains(taskRq.Title, "]") {
		httputils.SetStatus500(w, errors.New("jiraId is not allowed in task title"))
		return
	}

	task := &model.Task{
		Uuid:        uuid.NewString(),
		Title:       taskRq.Title,
		JiraId:      taskRq.JiraId,
		Description: taskRq.Description,
		State:       model.TASK_ACTIVE,
	}

	// Присваиваем задаче произвольного пользователя
	randomUser := CreateUserRandomizer()
	newUser := randomUser.Uuid()
	if newUser == "" {
		httputils.SetStatus500(w, errors.New("нет пользователей для назначения задаче"))
		return
	}
	task.AssignedUserUuid = newUser

	// Сохраняем задачу в БД
	repo, _ := base.App.GetDomainRepository("task")
	err = repo.Append(task)
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	// Формируем данные ответа на запрос
	taskResp := &TaskResponse{
		ID:          task.ID,
		Uuid:        task.Uuid,
		Title:       task.Title,
		JiraId:      task.JiraId,
		Description: task.Description,
		UserUuid:    task.AssignedUserUuid,
		State:       task.State,
	}

	err = json.NewEncoder(w).Encode(taskResp)
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	w.WriteHeader(201)

	// Отправляем CUD событие добавления задачи в очередь CUD-событий
	eventData := eventmodel.TaskEventData{
		Uuid:             task.Uuid,
		Title:            task.Title,
		JiraId:           task.JiraId,
		Description:      task.Description,
		AssignedUserUuid: task.AssignedUserUuid,
	}

	eventStreamData := eventmodel.TaskStreamData{
		TaskEventData: eventData,
		State:         string(task.State),
	}

	_, err = eventProducerCUD.ProduceEventData(eventmodel.TASK_CUD_TASK_CREATED, task.Uuid, reflect.TypeOf(*task).String(), eventStreamData, "2", []string{"1"})
	if err != nil {
		log.Fatal(err)
	}

	// Отправляем BE событие добавления задачи в очередь BE-событий
	_, err = eventProducerBE.ProduceEventData(eventmodel.TASK_BE_TASK_CREATED, task.Uuid, reflect.TypeOf(*task).String(), eventData, "2", []string{"1"})
	if err != nil {
		log.Fatal(err)
	}

	// BE-событие назначения задачи при создании не отправляем - сервису аккаунтинга будет достаточно события создания
	/*
		_, err = eventProducerBE.ProduceEventData(eventmodel.TASK_BE_TASK_ASSIGNED, task.Uuid, reflect.TypeOf(*task).String(), eventData, "2", []string{"1"})
		if err != nil {
			log.Fatal(err)
		}
	*/
}

// Обработка запроса пакетного переназначения пользователей у открытых задач
func handleReassignTask(w http.ResponseWriter, r *http.Request) {
	var tasks []*model.Task
	prevTaskUserAssignments := make(map[uint]string)

	// Получаем список пользователей для назначения
	randomizer := CreateUserRandomizer()
	if randomizer.Len() == 0 {
		httputils.SetStatus500(w, errors.New("нет пользователей для назначения"))
		return
	}

	// Получаем список открытых задач
	repo, _ := base.App.GetDomainRepository("task")
	result, err := repo.RawQuery("select id,uuid,description,assigned_user_uuid,state from task.task where state='ACTIVE'")
	if err != nil {
		log.Fatal(err)
	}
	rows := result.(*sql.Rows)
	defer rows.Close()

	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(&task.ID, &task.Uuid, &task.Description, &task.AssignedUserUuid, &task.State)
		if err != nil {
			log.Fatal(err)
		}
		// Сохпаняем текущего пользователя у заадчи (чтобы потом не отсылать события, если пользователь у задачи не поменялся)
		prevTaskUserAssignments[task.ID] = task.AssignedUserUuid

		// Устанавливаем нового пользовтаеля задаче
		task.AssignedUserUuid = randomizer.Uuid()

		tasks = append(tasks, task)
	}

	// Формируем ответ на запрос, отправляем CUD-событие изменения задачи, BE-событие переназначения пользователя у задачи
	var responseData ReassingTasksResponse

	for _, task := range tasks {
		if task.AssignedUserUuid != prevTaskUserAssignments[task.ID] {
			// Делаем, если у задачи поменялся пользователь
			// Обновляем задачу в БД
			err := repo.Update(task)
			if err != nil {
				log.Fatal(err)
			}
			// Теперь отсылаем CUD событие изменения задачи
			eventData := eventmodel.TaskEventData{
				Uuid:             task.Uuid,
				Title:            task.Title,
				JiraId:           task.JiraId,
				Description:      task.Description,
				AssignedUserUuid: task.AssignedUserUuid,
			}

			eventStreamData := eventmodel.TaskStreamData{
				TaskEventData: eventData,
				State:         string(task.State),
			}

			_, err = eventProducerCUD.ProduceEventData(eventmodel.TASK_CUD_TASK_UPDATED, task.Uuid, reflect.TypeOf(*task).String(), eventStreamData, "2", []string{"1"})
			if err != nil {
				log.Fatal(err)
			}

			// Отправляем BEсобытие изменения пользователя задачи в очередь BE-событий
			_, err = eventProducerBE.ProduceEventData(eventmodel.TASK_BE_TASK_ASSIGNED, task.Uuid, reflect.TypeOf(*task).String(), eventData, "2", []string{"1"})
			if err != nil {
				log.Fatal(err)
			}

			// Добавляем задачу в ответ
			responseData = append(responseData, ReassignTasksResponseItem{
				ID:       task.ID,
				Uuid:     task.Uuid,
				UserUuid: task.AssignedUserUuid,
			})
		}
	}
	// Формируем ответ
	err = json.NewEncoder(w).Encode(responseData)
	if err != nil {
		log.Fatal(err)
	}
}

// Обработка запроса на закрытие задачи
func handleCompleteTask(w http.ResponseWriter, r *http.Request) {
	// Ищем задачу с указанным ID
	taskID := r.PathValue("id")
	repo, _ := base.App.GetDomainRepository("task")
	task := &model.Task{}
	err := repo.Get(task, map[string]interface{}{"id": taskID, "state": "ACTIVE"})
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	// Получаем идентификатор пользовтаеля по данным авторизации (из jwt-токена)
	// Сам заголовок добавила обертка над методом (WithAuth) после проверки jwt-токена
	userUuid := r.Header.Get("X-Auth-User-UUID")

	// Проверяем, что пользовтаель закрывает совю задачу
	if task.AssignedUserUuid != userUuid {
		httputils.SetStatus401(w, "Task assigned to another user")
		return
	}

	// Устанавливаем статус завершения
	task.State = model.TASK_COMPLETED
	// Обновляем задачу в БД
	repo.Update(task)

	// Формируем ответ
	taskResp := &TaskResponse{
		ID:          task.ID,
		Uuid:        task.Uuid,
		Title:       task.Title,
		JiraId:      task.JiraId,
		Description: task.Description,
		UserUuid:    task.AssignedUserUuid,
		State:       task.State,
	}

	err = json.NewEncoder(w).Encode(taskResp)
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}

	// Отправляем CUD событие изменения задачи в очередь CUD-событий
	eventData := eventmodel.TaskEventData{
		Uuid:             task.Uuid,
		Description:      task.Description,
		AssignedUserUuid: task.AssignedUserUuid,
	}

	eventStreamData := eventmodel.TaskStreamData{
		TaskEventData: eventData,
		State:         string(task.State),
	}

	_, err = eventProducerCUD.ProduceEventData(eventmodel.TASK_CUD_TASK_UPDATED, task.Uuid, reflect.TypeOf(*task).String(), eventStreamData, "2", []string{"1"})
	if err != nil {
		log.Fatal(err)
	}

	// Отправляем BEсобытие закрытия задачи в очередь BE-событий
	_, err = eventProducerBE.ProduceEventData(eventmodel.TASK_BE_TASK_COMPLETED, task.Uuid, reflect.TypeOf(*task).String(), eventData, "2", []string{"1"})
	if err != nil {
		log.Fatal(err)
	}
}

// Обработка запроса на получение списка задач текущего пользователя
func handleGetUserTasks(w http.ResponseWriter, r *http.Request) {

	var tasks UserTasksList

	repo, _ := base.App.GetDomainRepository("task")
	result, err := repo.RawQuery("select id,uuid,description,assigned_user_uuid,state from task.task where state='ACTIVE' and assigned_user_uuid =?", r.Header.Get("X-Auth-User-UUID"))
	if err != nil {
		httputils.SetStatus500(w, err)
		return
	}
	rows := result.(*sql.Rows)
	defer rows.Close()

	for rows.Next() {
		task := &model.Task{}
		err := rows.Scan(&task.ID, &task.Uuid, &task.Description, &task.AssignedUserUuid, &task.State)
		if err != nil {
			log.Fatal(err)
		}

		tasks = append(tasks, TaskResponse{
			ID:          task.ID,
			Uuid:        task.Uuid,
			Title:       task.Title,
			JiraId:      task.JiraId,
			Description: task.Description,
			UserUuid:    task.AssignedUserUuid,
			State:       task.State,
		})
	}

	err = json.NewEncoder(w).Encode(&tasks)
	if err != nil {
		httputils.SetStatus500(w, err)
	}
}
