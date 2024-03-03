package gorm

import (
	repo "async-arch/internal/lib/base/repository"
	"errors"
	"fmt"
	"reflect"

	dbdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DomainRepository struct {
	db    *gorm.DB
	repos map[reflect.Type]*ObjectRepository
}

func CreateDomainRepository(host, dbname, scheme, user, password string) (*DomainRepository, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := gorm.Open(dbdriver.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", scheme),
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &DomainRepository{
		db:    db,
		repos: make(map[reflect.Type]*ObjectRepository),
	}, nil
}

func (r *DomainRepository) CreateObjectRepository(object interface{}) (repo.ObjectRepositoryManager, error) {
	repository := &ObjectRepository{
		domainRepo: r,
	}
	err := r.db.AutoMigrate(object)
	if err != nil {
		return nil, err
	}

	r.repos[reflect.TypeOf(object)] = repository

	return repository, nil
}

func (r *DomainRepository) GetObjectRepository(object interface{}) (repo.ObjectRepositoryManager, error) {
	if repository, ok := r.repos[reflect.TypeOf(object)]; !ok {
		return nil, errors.New("не найден репозиторий для этого типа")
	} else {
		return repository, nil
	}
}

func (r *DomainRepository) Append(object interface{}) error {

	repository, ok := r.repos[reflect.TypeOf(object)]
	if !ok {
		return errors.New("отсутствует репозиторий для этого типа")
	}

	return repository.Append(object)
}

func (r *DomainRepository) Update(object interface{}) error {
	repository, ok := r.repos[reflect.TypeOf(object)]
	if !ok {
		return errors.New("отсутствует репозиторий для этого типа")
	}

	return repository.Update(object)
}

func (r *DomainRepository) Delete(object interface{}) error {
	repository, ok := r.repos[reflect.TypeOf(object)]
	if !ok {
		return errors.New("отсутствует репозиторий для этого типа")
	}

	return repository.Delete(object)
}

func (r *DomainRepository) Get(object interface{}, filter interface{}) error {
	repository, ok := r.repos[reflect.TypeOf(object)]
	if !ok {
		return errors.New("отсутствует репозиторий для этого типа")
	}

	return repository.Get(object, filter)
}

func (r *DomainRepository) RawQuery(queryText string, args ...any) (interface{}, error) {
	result := r.db.Raw(queryText, args...)
	if result.Error != nil {
		return nil, result.Error
	}
	return result.Rows()
}

func (r *DomainRepository) Close() {
	// Заглушка для интерфейса
}

type ObjectRepository struct {
	domainRepo *DomainRepository
}

func (r *ObjectRepository) Append(object interface{}) error {
	if result := r.domainRepo.db.Create(object); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ObjectRepository) Update(object interface{}) error {
	if result := r.domainRepo.db.Save(object); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ObjectRepository) Delete(object interface{}) error {
	if result := r.domainRepo.db.Delete(object); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ObjectRepository) Get(object interface{}, filter interface{}) error {
	if filter == nil {
		if result := r.domainRepo.db.First(object); result.Error != nil {
			return result.Error
		}

		return nil
	}

	filterMap, ok := filter.(map[string]interface{})
	if !ok {
		return errors.New("некорректный формат фильтра")
	}

	if result := r.domainRepo.db.Where(filterMap).First(object); result.Error != nil {
		return result.Error
	}

	return nil
}
