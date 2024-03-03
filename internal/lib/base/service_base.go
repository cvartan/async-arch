package base

import (
	msg "async-arch/internal/lib/base/messages"
	repo "async-arch/internal/lib/base/repository"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/xlab/closer"
)

// serviceApplication - шаблон для сервиса
type serviceApplication struct {
	messageManagers map[string]msg.MessageManager
	httpServer      *http.Server
	httpClients     map[string]httpRequestProducer
	cancelFuncs     map[string]func()
	domains         map[string]repo.DomainRepositoryManager
}

// httpRequestProducer - шаблон http-запроса для клиента
type httpRequestProducer struct {
	serverAddr string
	methodURL  string
	methodType string
}

type HttpCookie struct {
	Name      string
	Value     string
	ExpiredAt time.Time
}

// RegisterMessageManager - метод регистрации нового менеджера очередей
func (app *serviceApplication) RegisterMessageManager(manager msg.MessageManager) {
	if _, ok := app.messageManagers[manager.ID()]; !ok {
		app.messageManagers[manager.ID()] = manager
	}
}

// RegisterMessageProducer - метод регистрации нового публикатора сообщений
func (app *serviceApplication) RegisterMessageProducer(managerId, producerId string, queueName string) (msg.MessageProducer, error) {

	if managerId == "" {
		return nil, errors.New("необходимо указать имя используемого менеджера очередей")
	}

	manager, ok := app.messageManagers[managerId]

	if !ok {
		return nil, errors.New("менеджер с таким Id не найден")
	}

	if producerId == "" {
		return nil, errors.New("необходимо указать имя-ключ для публикатора сообщений")
	}

	if queueName == "" {
		return nil, errors.New("необходимо указать имя очереди для публикации")
	}

	if _, ok := manager.GetProducer(producerId); ok {
		return nil, errors.New("публикатор с таким именем-ключом уже добавлен")
	}

	if err := manager.AddProducer(producerId, queueName); err != nil {
		return nil, err
	} else {
		producer, _ := manager.GetProducer(producerId)
		return producer, nil
	}
}

// SendMsg - метод отправки сообщений через укзанного публикатора сообщений
func (app *serviceApplication) ProduceMessage(managerId, producerId string, key []byte, value []byte, headers map[string]interface{}) error {
	if managerId == "" {
		return errors.New("необходимо указать имя используемого менеджера очередей")
	}

	manager, ok := app.messageManagers[managerId]

	if !ok {
		return errors.New("менеджер с таким Id не найден")
	}

	if producerId == "" {
		return errors.New("необходимо указать имя-ключ для публикатора")
	}

	if len(key) == 0 {
		return errors.New("необходимо указать ключ для сообщения")
	}

	if producer, ok := manager.GetProducer(producerId); !ok {
		return errors.New("публикатор с таким именем не найден")
	} else {
		return producer.ProduceMessage(key, value, headers)
	}
}

// SendStrMsg - метод отправки сообщений в строковом формате через указанного публикатора сообщений
func (app *serviceApplication) ProduceStringMessage(managerId, producerId string, key string, value string, headers map[string]interface{}) error {
	return app.ProduceMessage(managerId, producerId, []byte(key), []byte(value), headers)
}

// Consume - метод прослушивания очереди
func (app *serviceApplication) Consume(managerId, queueName string, handler msg.MessageHandler) error {
	if managerId == "" {
		return errors.New("необходимо указать имя используемого менеджера очередей")
	}

	manager, ok := app.messageManagers[managerId]
	if !ok {
		return errors.New("менеджер с таким Id не найден")
	}

	return manager.Consume(queueName, handler)

}

// InitHTTPServer - инициализация http сервера (запуск сервера выполняется при старте приложения)
func (app *serviceApplication) InitHTTPServer(address string, port uint16) error {
	if app.httpServer != nil {
		return errors.New("http-сервер уже инициализирован")
	}
	if port == 0 {
		port = 80
	}
	completeAddress := fmt.Sprintf("%s:%d", address, port)

	mux := http.NewServeMux()

	app.httpServer = &http.Server{
		Addr:    completeAddress,
		Handler: mux,
	}

	log.Printf("Starting http server at %s\n", app.httpServer.Addr)
	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

// HandleFunc - назначение метода обработки шаблону пути (пути с параметрами требуют версии go от 1)
func (app *serviceApplication) HandleFunc(methodString string, handler http.HandlerFunc) error {
	if app.httpServer == nil {
		return errors.New("http-сервер должен быть сначала инициализрован")
	}
	if methodString == "" {
		return errors.New("путь к методу должен быть определен")
	}
	if handler == nil {
		return errors.New("метод обработки должен быть определен")
	}
	app.httpServer.Handler.(*http.ServeMux).HandleFunc(methodString, handler)
	return nil
}

// AddGetRequest - добавление шаблон для запроса GET
func (app *serviceApplication) AddGetRequest(requestId, serverAddr, methodUrl string) error {
	if requestId == "" {
		return errors.New("идентификатор запроса должен быть указан")
	}
	if serverAddr == "" {
		return errors.New("адрес сервера должен быть указан")
	}
	if methodUrl == "" {
		return errors.New("путь к методу должэен быть указан")
	}
	if _, ok := app.httpClients[requestId]; ok {
		return errors.New("запрос с таким идентификатором уже определен")
	}
	app.httpClients[requestId] = httpRequestProducer{
		serverAddr: serverAddr,
		methodURL:  methodUrl,
		methodType: http.MethodGet,
	}
	return nil
}

// AddPostRequest - добавление шаблон для запроса POST
func (app *serviceApplication) AddPostRequest(requestId, serverAddr, methodUrl string) error {
	if requestId == "" {
		return errors.New("идентификатор запроса должен быть указан")
	}
	if serverAddr == "" {
		return errors.New("адрес сервера должен быть указан")
	}
	if methodUrl == "" {
		return errors.New("путь к методу должэен быть указан")
	}
	if _, ok := app.httpClients[requestId]; ok {
		return errors.New("запрос с таким идентификатором уже определен")
	}
	app.httpClients[requestId] = httpRequestProducer{
		serverAddr: serverAddr,
		methodURL:  methodUrl,
		methodType: http.MethodPost,
	}
	return nil
}

// AddPutRequest - добавление шаблона для запроса PUT
func (app *serviceApplication) AddPutRequest(requestId, serverAddr, methodUrl string) error {
	if requestId == "" {
		return errors.New("идентификатор запроса должен быть указан")
	}
	if serverAddr == "" {
		return errors.New("адрес сервера должен быть указан")
	}
	if methodUrl == "" {
		return errors.New("путь к методу должэен быть указан")
	}
	if _, ok := app.httpClients[requestId]; ok {
		return errors.New("запрос с таким идентификатором уже определен")
	}
	app.httpClients[requestId] = httpRequestProducer{
		serverAddr: serverAddr,
		methodURL:  methodUrl,
		methodType: http.MethodPut,
	}
	return nil
}

// AddDeleteRequest - добавление шаблона для запроса DELETE
func (app *serviceApplication) AddDeleteRequest(requestId, serverAddr, methodUrl string) error {
	if requestId == "" {
		return errors.New("идентификатор запроса должен быть указан")
	}
	if serverAddr == "" {
		return errors.New("адрес сервера должен быть указан")
	}
	if methodUrl == "" {
		return errors.New("путь к методу должэен быть указан")
	}
	if _, ok := app.httpClients[requestId]; ok {
		return errors.New("шаблон запроса с таким идентификатором уже определен")
	}
	app.httpClients[requestId] = httpRequestProducer{
		serverAddr: serverAddr,
		methodURL:  methodUrl,
		methodType: http.MethodGet,
	}
	return nil
}

// Request - выполнение запроса HTTP по шаблону
func (app *serviceApplication) Request(requestId string, body []byte, params, query map[string]interface{}, headers map[string]string, cookies []HttpCookie) (*http.Response, error) {
	c, ok := app.httpClients[requestId]
	if !ok {
		return nil, errors.New("шаблон запроса с таким идентификатором не найден")
	}

	url := c.serverAddr + func(baseString string, changeMap map[string]interface{}) string {
		if changeMap == nil {
			return baseString
		}
		changedString := baseString
		for key, value := range changeMap {
			changedString = strings.ReplaceAll(changedString, "{"+key+"}", fmt.Sprintf("%v", value))
		}
		return changedString

	}(c.methodURL, params)

	if query != nil {
		queryStr := "?"
		len := len(query)
		i := 0
		for key, value := range query {
			queryStr = queryStr + key + "=" + fmt.Sprintf("%v", value)
			i++
			if i < len {
				queryStr = queryStr + "&"
			}
		}
		url = url + queryStr
	}

	client := http.Client{Timeout: time.Minute}

	req, err := http.NewRequest(c.methodType, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for header, value := range headers {
		req.Header.Add(header, value)
	}

	for _, cookie := range cookies {
		req.AddCookie(&http.Cookie{
			Name:    cookie.Name,
			Value:   cookie.Value,
			Expires: cookie.ExpiredAt,
		})
	}

	return client.Do(req)
}

// RegisterDomainRepository - регистрация нового репозитория для домена данных
func (app *serviceApplication) RegisterDomainRepository(repoId string, repository repo.DomainRepositoryManager) error {
	if repoId == "" {
		return errors.New("идентификатор домена должен быть указан")
	}
	if repository == nil {
		return errors.New("репозиторий домена должен быть указан")
	}
	if _, ok := app.domains[repoId]; ok {
		return errors.New("репозиторий с таким идентификатором уже добавлен")
	}
	app.domains[repoId] = repository
	return nil
}

func (app *serviceApplication) GetDomainRepository(repoId string) (repo.DomainRepositoryManager, error) {
	repository, ok := app.domains[repoId]
	if !ok {
		return nil, errors.New("репозиторий с таким идентфиикатором не найден")
	}

	return repository, nil
}

// RegisterCancelFunc - регистрация закрывающей функции (которая выполнится при выходе из приложения)
func (app *serviceApplication) RegisterCancelFunc(alias string, f func()) {
	app.cancelFuncs[alias] = f
}

// Hold - удержание приложения в работе
func (app *serviceApplication) Hold() {
	closer.Bind(app.Close)
	closer.Hold()
}

// Close - метод закрытия приложения
func (app *serviceApplication) Close() {
	log.Println("Starting closing process")

	// 1. Выполняем закрывающие операции
	for alias, f := range app.cancelFuncs {
		log.Printf("Runnig cancel func with id: %s\n", alias)
		f()
	}

	// 2. Завершаем работу http сервера
	if app.httpServer != nil {
		log.Println("Shutdown http server")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		app.httpServer.Shutdown(ctx)
	}

	// 3. Закрываем соединения для менеджеров очередей
	for id, manager := range app.messageManagers {
		log.Printf("Shutdown queue connection: %s\n", id)
		manager.Close()
	}

	// 4. Закрываем соединения с репозиториями
	for id, repository := range app.domains {
		log.Printf("Shutdown database connection for %s\n", id)
		repository.Close()
	}

	log.Println("Application has been closed")
}

var App serviceApplication

func init() {
	App = serviceApplication{
		messageManagers: make(map[string]msg.MessageManager),
		httpClients:     make(map[string]httpRequestProducer),
		domains:         make(map[string]repo.DomainRepositoryManager),
	}
	log.Println("Application started")
}
