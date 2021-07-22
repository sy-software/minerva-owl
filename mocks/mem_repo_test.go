package mocks

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sy-software/minerva-owl/internal/core/ports"
)

type Companion interface{}

type Pokemon struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Generation int    `json:"generation,omitempty"`
}

// This is a little redundant but I need a way to check
// if my generic MemRepo that can store any model works
// as expected.
func TestMemRepo(t *testing.T) {

	t.Run("Test list action", func(t *testing.T) {
		expected := []Pokemon{
			{
				Id:   "1",
				Name: "Bulbasaur",
			},
			{
				Id:   "2",
				Name: "Ivysaur",
			},
			{
				Id:   "3",
				Name: "Venosaur",
			},
		}

		pokemons := []map[string]interface{}{
			{
				"id":   "1",
				"name": "Bulbasaur",
			},
			{
				"id":   "2",
				"name": "Ivysaur",
			},
			{
				"id":   "3",
				"name": "Venosaur",
			},
		}

		data := map[string][]map[string]interface{}{
			"pokemons": pokemons,
		}

		repo := MemRepo{
			Data: data,
		}

		var got []Pokemon
		err := repo.List("pokemons", &got, 0, 10)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !cmp.Equal(got, expected) {
			t.Errorf("Expected result: %+v. Got: %+v", expected, got)
		}
	})

	t.Run("Test get action", func(t *testing.T) {
		expected := Pokemon{
			Id:   "2",
			Name: "Ivysaur",
		}

		pokemons := []map[string]interface{}{
			{
				"id":   "1",
				"name": "Bulbasaur",
			},
			{
				"id":   "2",
				"name": "Ivysaur",
			},
			{
				"id":   "3",
				"name": "Venosaur",
			},
		}

		data := map[string][]map[string]interface{}{
			"pokemons": pokemons,
		}

		repo := MemRepo{
			Data: data,
		}

		got := Pokemon{}
		err := repo.Get("pokemons", "2", &got)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !cmp.Equal(got, expected) {
			t.Errorf("Expected result: %+v. Got: %+v", expected, got)
		}
	})

	t.Run("Test create action", func(t *testing.T) {
		expected := Pokemon{
			Name: "Ivysaur",
		}

		pokemons := []map[string]interface{}{
			{
				"id":   "1",
				"name": "Bulbasaur",
			},
			{
				"id":   "3",
				"name": "Venosaur",
			},
		}

		data := map[string][]map[string]interface{}{
			"pokemons": pokemons,
		}

		repo := MemRepo{
			Data: data,
		}

		id, err := repo.Create("pokemons", expected)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		got := Pokemon{}
		err = repo.Get("pokemons", id, &got)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if expected.Name != got.Name {
			t.Errorf("Expected result: %+v. Got: %+v", expected.Name, got.Name)
		}
	})

	t.Run("Test update action", func(t *testing.T) {
		expected := Pokemon{
			Id:   "2",
			Name: "Ivysaur",
		}

		pokemons := []map[string]interface{}{
			{
				"id":   "1",
				"name": "Bulbasaur",
			},
			{
				"id":   "2",
				"name": "ivysaur",
			},
			{
				"id":   "3",
				"name": "Venosaur",
			},
		}

		data := map[string][]map[string]interface{}{
			"pokemons": pokemons,
		}

		repo := MemRepo{
			Data: data,
		}

		err := repo.Update("pokemons", "2", expected)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		got := Pokemon{}
		err = repo.Get("pokemons", "2", &got)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !cmp.Equal(got, expected) {
			t.Errorf("Expected result: %+v. Got: %+v", expected, got)
		}
	})

	t.Run("Test update omit fields", func(t *testing.T) {
		expected := Pokemon{
			Id:         "2",
			Name:       "Ivysaur",
			Generation: 10,
		}

		pokemons := []map[string]interface{}{
			{
				"id":   "1",
				"name": "Bulbasaur",
			},
			{
				"id":         "2",
				"name":       "ivysaur",
				"generation": 1,
			},
			{
				"id":   "3",
				"name": "Venosaur",
			},
		}

		data := map[string][]map[string]interface{}{
			"pokemons": pokemons,
		}

		repo := MemRepo{
			Data: data,
		}

		err := repo.Update("pokemons", "2", expected, "generation")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		got := Pokemon{}
		err = repo.Get("pokemons", "2", &got)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.Generation == expected.Generation {
			t.Errorf("Expected result: 1. Got: %+v", got)
		}
	})

	t.Run("Test delete action", func(t *testing.T) {
		pokemons := []map[string]interface{}{
			{
				"id":   "1",
				"name": "Bulbasaur",
			},
			{
				"id":   "2",
				"name": "Ivysaur",
			},
			{
				"id":   "3",
				"name": "Venosaur",
			},
			{
				"id":   "deleteMe",
				"name": "Not A Pokemon",
			},
		}

		data := map[string][]map[string]interface{}{
			"pokemons": pokemons,
		}

		repo := MemRepo{
			Data: data,
		}

		err := repo.Delete("pokemons", "deleteMe")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		got := Pokemon{}
		err = repo.Get("pokemons", "deleteMe", &got)

		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if _, ok := err.(ports.ErrItemNotFound); !ok {
			t.Errorf("Expected error or type ErrItemNotFound got: %v", err)
		}
	})
}

func TestInterceptors(t *testing.T) {
	t.Run("Test list interceptor", func(t *testing.T) {
		called := false
		interceptor := func(collection string, results interface{}, skip, limit int, filters ...ports.Filter) error {
			called = true
			return nil
		}

		repo := MemRepo{
			Data:            map[string][]map[string]interface{}{},
			ListInterceptor: interceptor,
		}

		var res interface{}
		repo.List("coll", res, 0, 0)

		if !called {
			t.Errorf("Expected ListInterceptor to be called")
		}
	})

	t.Run("Test get interceptor", func(t *testing.T) {
		called := false
		interceptor := func(collection, id string, result interface{}) error {
			called = true
			return nil
		}

		repo := MemRepo{
			Data:           map[string][]map[string]interface{}{},
			GetInterceptor: interceptor,
		}

		var res interface{}
		repo.Get("col", "id", res)

		if !called {
			t.Errorf("Expected GetInterceptor to be called")
		}
	})

	t.Run("Test get one interceptor", func(t *testing.T) {
		called := false
		interceptor := func(collection string, result interface{}, filters ...ports.Filter) error {
			called = true
			return nil
		}
		repo := MemRepo{
			Data:              map[string][]map[string]interface{}{},
			GetOneInterceptor: interceptor,
		}

		var res interface{}
		repo.GetOne("col", res)

		if !called {
			t.Errorf("Expected GetOneInterceptor to be called")
		}
	})

	t.Run("Test list interceptor", func(t *testing.T) {
		called := false
		interceptor := func(collection string, entity interface{}) (string, error) {
			called = true
			return "", nil
		}
		repo := MemRepo{
			Data:              map[string][]map[string]interface{}{},
			CreateInterceptor: interceptor,
		}

		var res interface{}
		repo.Create("col", res)

		if !called {
			t.Errorf("Expected CreateInterceptor to be called")
		}
	})

	t.Run("Test update interceptor", func(t *testing.T) {
		called := false
		interceptor := func(collection, id string, entity interface{}, omit ...string) error {
			called = true
			return nil
		}
		repo := MemRepo{
			Data:              map[string][]map[string]interface{}{},
			UpdateInterceptor: interceptor,
		}

		var res interface{}
		repo.Update("col", "id", res)

		if !called {
			t.Errorf("Expected UpdateInterceptor to be called")
		}
	})

	t.Run("Test delete interceptor", func(t *testing.T) {
		called := false
		interceptor := func(collection, id string) error {
			called = true
			return nil
		}
		repo := MemRepo{
			Data:              map[string][]map[string]interface{}{},
			DeleteInterceptor: interceptor,
		}

		repo.Delete("coll", "id")

		if !called {
			t.Errorf("Expected DeleteInterceptor to be called")
		}
	})
}
