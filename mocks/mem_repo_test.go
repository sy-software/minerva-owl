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
