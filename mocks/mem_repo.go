package mocks

import (
	"encoding/json"
	"reflect"

	"github.com/google/uuid"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

type MemRepo struct {
	Data map[string][]map[string]interface{}
}

func (repo *MemRepo) List(collection string, results interface{}, skip int, limit int, filters ...ports.Filter) error {
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
		Id:    id,
		Model: collection,
	}
}

func (repo *MemRepo) Create(collection string, entity interface{}) (string, error) {
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

func (repo *MemRepo) Update(collection string, id string, entity interface{}) error {
	colData := repo.Data[collection]
	for i, item := range colData {
		if item["id"] == id {
			var jsonMap map[string]interface{}
			doc, err := json.Marshal(entity)

			if err != nil {
				return err
			}

			json.Unmarshal(doc, &jsonMap)
			colData[i] = jsonMap
		}
	}
	return nil
}

func (repo *MemRepo) Delete(collection string, id string) error {
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
