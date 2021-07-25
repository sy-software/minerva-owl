package handlers

import (
	"regexp"
	"testing"
	"time"

	"github.com/sy-software/minerva-owl/cmd/graphql/graph/model"
	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/core/service"
	"github.com/sy-software/minerva-owl/internal/utils"
	"github.com/sy-software/minerva-owl/mocks"
)

const authKey = "2b7e151628aed2a6abf71589a12b4da32"

func TestUserCreateOperation(t *testing.T) {
	t.Run("Create an User", func(t *testing.T) {
		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		now := utils.UnixUTCNow()
		input := model.NewUser{
			Name:     "Tony Stark",
			Username: "IronMan",
			Role:     "hero",
			Provider: "avengers",
			TokenID:  "mytoken",
			Status:   "deceased",
		}
		got, err := handlerInstance.Create(input)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		match, err := regexp.MatchString(mocks.ID_REGEX, got.ID)
		if !match || err != nil {
			t.Errorf("ID is not V4 UUID got: %q with error: %v", got.ID, err)
		}

		if got.Name != input.Name {
			t.Errorf("Expected Name to be: %q got: %q", input.Name, got.Name)
		}

		if got.Username != input.Username {
			t.Errorf("Expected Username to be: %q got: %q", input.Username, got.Username)
		}

		if *got.Picture != "" {
			t.Errorf("Expected.Picture to be \"\" got: %q", *got.Picture)
		}

		decrypted, err := utils.AES256Decrypt(config.Keys.Auth, got.TokenID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if decrypted != input.TokenID {
			t.Errorf("Expected TokenID to be: %q got: %q", input.TokenID, decrypted)
		}

		if got.CreateDate.After(now) {
			t.Errorf("Expected CreatedDate to be approximately equals: %q got: %q", now, got.CreateDate)
		}

		if got.UpdateDate.After(now) {
			t.Errorf("Expected UpdateDate to be approximately equals: %q got: %q", now, got.CreateDate)
		}
	})
}

func TestReadOperations(t *testing.T) {
	t.Run("List Users", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":       "1",
				"username": "CapAmerica",
			},
			{
				"id":       "2",
				"username": "IronMan",
			},
		}

		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		page := 1
		size := 10
		got, err := handlerInstance.Query(nil, &page, &size)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(got) != 2 {
			t.Errorf("Expected %d results got %d", 2, len(got))
		}
	})

	t.Run("List Users By Role", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":       "1",
				"username": "CapAmerica",
			},
			{
				"id":       "2",
				"username": "IronMan",
			},
		}

		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: dummyData,
		}

		called := false
		repo := mocks.MemRepo{
			Data: data,
			ListInterceptor: func(collection string, results interface{}, skip, limit int, filters ...ports.Filter) error {
				called = true

				if len(filters) != 1 {
					t.Errorf("Expected 1 filter got %d", len(filters))
				}

				if filters[0].Name != "role" {
					t.Errorf("Expected filter by key to be \"role\" got %q", filters[0].Name)
				}

				if filters[0].Value != "genius" {
					t.Errorf("Expected filter by value to be \"genius\" got %q", filters[0].Value)
				}

				return nil
			},
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		page := 1
		size := 10
		role := "genius"
		_, err := handlerInstance.Query(&role, &page, &size)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !called {
			t.Errorf("Expected List to be called")
		}
	})

	t.Run("Get User By Id", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":       "1",
				"username": "CapAmerica",
			},
			{
				"id":       "2",
				"username": "IronMan",
			},
		}

		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		got, err := handlerInstance.QueryById("1")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.ID != "1" {
			t.Errorf("Expected ID to be: %q got: %q", "1", got.ID)
		}
	})

	t.Run("Get User By Username", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":       "1",
				"username": "CapAmerica",
			},
			{
				"id":       "2",
				"username": "IronMan",
			},
		}

		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: dummyData,
		}

		called := false
		repo := mocks.MemRepo{
			Data: data,
			GetOneInterceptor: func(collection string, result interface{}, filters ...ports.Filter) error {
				called = true
				if len(filters) != 1 {
					t.Errorf("Expected 1 filter got %d", len(filters))
				}

				if filters[0].Name != "username" {
					t.Errorf("Expected filter by key to be \"username\" got %q", filters[0].Name)
				}

				if filters[0].Value != "IronMan" {
					t.Errorf("Expected filter by value to be \"IronMan\" got %q", filters[0].Value)
				}
				return nil
			},
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		_, err := handlerInstance.QueryByUsername("IronMan")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !called {
			t.Errorf("Expected GetOne to be called")
		}
	})
}

func TestUserUpdateOperation(t *testing.T) {
	tokenId := "myTokenId"
	encrypted, _ := utils.AES256Encrypt(authKey, tokenId)
	now := utils.UnixUTCNow()
	yesterday := now.Add(-24 * time.Hour)
	t.Run("Update an User", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":         "1",
				"username":   "CapAmerica",
				"name":       "Steve Rogers",
				"picture":    "",
				"tokenID":    encrypted,
				"createDate": yesterday,
				"updateDate": yesterday,
				"status":     "active",
			},
			{
				"id":       "2",
				"username": "other",
			},
		}

		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		input := model.UpdateUser{
			ID:       "1",
			Name:     "Sam Wilson",
			Username: "CapAmerica",
			Role:     "hero",
			Provider: "avengers",
			TokenID:  "newTokenId",
			Status:   "active",
		}
		got, err := handlerInstance.Update(input)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.ID != input.ID {
			t.Errorf("Expected ID to be: %q got: %q", input.Name, got.Name)
		}

		if got.Name != input.Name {
			t.Errorf("Expected Name to be: %q got: %q", input.Username, got.Username)
		}

		if got.Username != input.Username {
			t.Errorf("Expected Username to be: %q got: %q", input.Username, got.Username)
		}

		if *got.Picture != "" {
			t.Errorf("Expected.Picture to be \"\" got: %q", *got.Picture)
		}

		decrypted, err := utils.AES256Decrypt(config.Keys.Auth, got.TokenID)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if decrypted != input.TokenID {
			t.Errorf("Expected TokenID to be: %q got: %q", input.TokenID, decrypted)
		}

		if got.CreateDate.Equal(yesterday) {
			t.Errorf("Expected CreatedDate to be approximately equals: %q got: %q", now, got.CreateDate)
		}

		if got.UpdateDate.After(now) {
			t.Errorf("Expected UpdateDate to be approximately equals: %q got: %q", now, got.CreateDate)
		}
	})
}

func TestUserDeleteOperation(t *testing.T) {
	tokenId := "myTokenId"
	encrypted, _ := utils.AES256Encrypt(authKey, tokenId)
	now := utils.UnixUTCNow()
	yesterday := now.Add(-24 * time.Hour)
	t.Run("Delete an User", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":         "1",
				"username":   "CapAmerica",
				"name":       "Steve Rogers",
				"picture":    "",
				"tokenID":    encrypted,
				"createDate": yesterday,
				"updateDate": yesterday,
				"status":     "active",
			},
			{
				"id":       "2",
				"username": "other",
			},
		}

		data := map[string][]map[string]interface{}{
			domain.USER_COL_NAME: dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}
		config := domain.DefaultConfig()
		config.Keys.Auth = authKey
		service := service.NewUserService(&repo, config)
		handlerInstance := NewUserGraphqlHandler(*service)

		got, err := handlerInstance.Delete("1")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.ID != "1" {
			t.Errorf("Expected ID to be: \"1\" got: %q", got.ID)
		}

		_, err = service.Get("1")

		if err == nil {
			t.Errorf("Expected error got nil")
		}

		_, ok := err.(ports.ErrItemNotFound)
		if !ok {
			t.Errorf("Expected error of type ErrItemNotFound got: %T", err)
		}
	})
}
