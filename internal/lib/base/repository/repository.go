package repository

// RepositoryStorageManager - интерфейс для реализации менеджера репозиториев домена данных (по факту коннект к экземпляру БД, где содержаться данные объектов)
type DomainRepositoryManager interface {
	// CreateObjectRepository - создание репозитория для объектов
	CreateObjectRepository(interface{}) (ObjectRepositoryManager, error)
	// GetObjectRepository - возвращает репозиторий для определенного типа данных
	GetObjectRepository(interface{}) (ObjectRepositoryManager, error)
	Append(newObject interface{}) error                       // Добавление нового объекта в репозиторий
	Update(updatedObject interface{}) error                   // Обновление объекта в репозитории
	Delete(deletedObject interface{}) error                   // Удаление объекта из репозитория
	Get(selectedObject interface{}, filter interface{}) error // Получение объекта из репозитория по соответствующему условию
	// RawQuery - выполняет произвольный запрос и возвращает результат (теоретически в формате sql.Rows)
	RawQuery(string, ...any) (interface{}, error)
	// Close - закрытие репозитория
	Close()
}

// ObjectRepositoryManager - интерфейс для реализации мененеджера репозитория объекта
type ObjectRepositoryManager interface {
	Append(newObject interface{}) error                       // Добавление нового объекта в репозиторий
	Update(updatedObject interface{}) error                   // Обновление объекта в репозитории
	Delete(deletedObject interface{}) error                   // Удаление объекта из репозитория
	Get(selectedObject interface{}, filter interface{}) error // Получение объекта из репозитория по соответствующему условию
}
