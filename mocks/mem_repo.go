package mocks

import (
	"encoding/json"
	"reflect"

	"github.com/google/uuid"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

const ID_REGEX = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"

type MemRepo struct {
	Data              map[string][]map[string]interface{}
	ListInterceptor   func(collection string, results interface{}, skip int, limit int, filters ...ports.Filter) error
	GetInterceptor    func(collection string, id string, result interface{}) error
	GetOneInterceptor func(collection string, result interface{}, filters ...ports.Filter) error
	CreateInterceptor func(collection string, entity interface{}) (string, error)
	UpdateInterceptor func(collection string, id string, entity interface{}, omit ...string) error
	DeleteInterceptor func(collection string, id string) error
}

func (repo *MemRepo) List(collection string, results interface{}, skip int, limit int, filters ...ports.Filter) error {
	if repo.ListInterceptor != nil {
		return repo.ListInterceptor(collection, results, skip, limit, filters...)
	}

	colData := repo.Data[collection]

	if skip >= len(colData) {
		return nil
	}

	available := len(colData) - skip
	capLimit := limit
	if limit > available {
		capLimit = available
	}

	resultsRaw := colData[skip : skip+capLimit]

	resultsPtr := reflect.ValueOf(results)
	resultsVal := resultsPtr.Elem()
	elementType := resultsVal.Type().Elem()
	for _, r := range resultsRaw {
		newElement := reflect.New(elementType).Elem()
		jsonbody, err := json.Marshal(r)

		if err != nil {
			return err
		}

		parsed := newElement.Addr().Interface()
		err = json.Unmarshal(jsonbody, &parsed)

		if err != nil {
			return err
		}

		resultsVal.Set(reflect.Append(resultsVal, newElement))
	}

	return nil
}

func (repo *MemRepo) Get(collection string, id string, result interface{}) error {
	if repo.GetInterceptor != nil {
		return repo.GetInterceptor(collection, id, result)
	}

	colData := repo.Data[collection]

	elementPtr := reflect.ValueOf(result)
	elementVal := elementPtr.Elem()
	elementType := elementVal.Type()
	for _, item := range colData {
		if item["id"] == id {
			newElement := reflect.New(elementType).Elem()
			jsonbody, err := json.Marshal(item)

			if err != nil {
				return err
			}

			parsed := newElement.Addr().Interface()
			err = json.Unmarshal(jsonbody, &parsed)

			if err != nil {
				return err
			}

			elementVal.Set(newElement)
			return nil
		}
	}

	return ports.ErrItemNotFound{
		Id:    &id,
		Model: collection,
	}
}

func (repo *MemRepo) GetOne(collection string, result interface{}, filters ...ports.Filter) error {
	if repo.GetOneInterceptor != nil {
		return repo.GetOneInterceptor(collection, result, filters...)
	}
	// TODO: Support queries in our mock db
	return nil
}

func (repo *MemRepo) Create(collection string, entity interface{}) (string, error) {
	if repo.CreateInterceptor != nil {
		return repo.CreateInterceptor(collection, entity)
	}
	var inInterface map[string]interface{}
	doc, err := json.Marshal(entity)
	json.Unmarshal(doc, &inInterface)

	newId := uuid.New().String()

	inInterface["id"] = newId

	items := repo.Data[collection]
	items = append(items, inInterface)
	repo.Data[collection] = items

	return newId, err
}

func (repo *MemRepo) Update(collection string, id string, entity interface{}, omit ...string) error {
	if repo.UpdateInterceptor != nil {
		return repo.UpdateInterceptor(collection, id, entity, omit...)
	}
	colData := repo.Data[collection]

	omitMap := map[string]bool{}
	for _, v := range omit {
		omitMap[v] = true
	}

	for _, item := range colData {
		if item["id"] == id {
			var jsonMap map[string]interface{}
			doc, err := json.Marshal(entity)

			if err != nil {
				return err
			}

			json.Unmarshal(doc, &jsonMap)

			for k, v := range jsonMap {
				_, shouldOmit := omitMap[k]
				if k != "id" && !shouldOmit {
					item[k] = v
				}

			}
			// colData[i] = jsonMap
		}
	}
	return nil
}

func (repo *MemRepo) Delete(collection string, id string) error {
	if repo.DeleteInterceptor != nil {
		return repo.DeleteInterceptor(collection, id)
	}
	colData := repo.Data[collection]
	newData := []map[string]interface{}{}
	for _, item := range colData {
		if item["id"] != id {
			newData = append(newData, item)
		}
	}

	repo.Data[collection] = newData
	return nil
}
