package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/sy-software/minerva-owl/internal/core/domain"
	"github.com/sy-software/minerva-owl/internal/core/ports"
	"github.com/sy-software/minerva-owl/internal/utils"
	"github.com/sy-software/minerva-owl/mocks"
)

const authKey = "2b7e151628aed2a6abf71589a12b4da32"

func TestCreateOperations(t *testing.T) {
	config := domain.DefaultConfig()
	config.Keys = domain.KeyList{
		Auth: authKey,
	}

	t.Run("Test User is created", func(t *testing.T) {
		now := utils.UnixUTCNow()

		expected := domain.User{
			Name:       "Tony Stark",
			Username:   "ironman",
			Picture:    "https://mypicture/ironman.png",
			Role:       "hero",
			Provider:   "marvel",
			TokenID:    "myToken",
			CreateDate: now,
			UpdateDate: now,
			Status:     "deceased",
		}

		data := map[string][]map[string]interface{}{
			"users": {},
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		var service ports.UserService
		service = NewUserService(&repo, config)

		created, err := service.Create(
			expected.Name,
			expected.Username,
			expected.Picture,
			expected.Role,
			expected.Provider,
			expected.TokenID,
			expected.Status,
		)

		if err != nil {
			t.Errorf("Item should be created without errors: %v", err)
		}

		if created.Name != expected.Name {
			t.Errorf(
				"Name was not assigned expected: %q got: %q",
				expected.Name,
				created.Name,
			)
		}

		if created.Username != expected.Username {
			t.Errorf(
				"Username was not assigned expected: %q got: %q",
				expected.Username,
				created.Username,
			)
		}

		if created.Picture != expected.Picture {
			t.Errorf(
				"Picture was not assigned expected: %q got: %q",
				expected.Picture,
				created.Picture,
			)
		}

		if created.Role != expected.Role {
			t.Errorf(
				"Role was not assigned expected: %q got: %q",
				expected.Role,
				created.Role,
			)
		}

		if created.Provider != expected.Provider {
			t.Errorf(
				"Provider was not assigned expected: %q got: %q",
				expected.Provider,
				created.Provider,
			)
		}

		decryptedToken, err := utils.AES256Decrypt(config.Keys.Auth, created.TokenID)
		if err != nil {
			t.Errorf("Expected token id to be decrypted without errors: %v", err)
		}

		if decryptedToken != expected.TokenID {
			t.Errorf(
				"TokenID was not encrypted as expected: %q got: %q",
				expected.TokenID,
				decryptedToken,
			)
		}

		if !created.CreateDate.After(now) && !created.CreateDate.Equal(now) {
			t.Errorf(
				"CreateDate must be after or equal to baseline now: %d got %d",
				now.UnixNano(),
				created.CreateDate.UnixNano(),
			)
		}

		if !created.UpdateDate.After(now) && !created.UpdateDate.Equal(now) {
			t.Errorf(
				"UpdateDate must be after or equal to baseline now: %d got %d",
				now.UnixNano(),
				created.UpdateDate.UnixNano(),
			)
		}

		if created.Status != expected.Status {
			t.Errorf(
				"Status was not assigned expected: %q got: %q",
				expected.Status,
				created.Status,
			)
		}
	})

	t.Run("Test User with duplicated Username can't be created", func(t *testing.T) {
		now := utils.UnixUTCNow()

		expected := domain.User{
			Name:       "Steve Rogers",
			Username:   "CaptainAmerica",
			Picture:    "https://mypicture/cap.png",
			Role:       "hero",
			Provider:   "marvel",
			TokenID:    "myToken",
			CreateDate: now,
			UpdateDate: now,
			Status:     "inactive",
		}

		data := map[string][]map[string]interface{}{
			"users": {},
		}

		repo := mocks.MemRepo{
			Data: data,
			GetOneInterceptor: func(collection string, result interface{}, filters ...ports.Filter) error {
				elementPtr := reflect.ValueOf(result)
				elementVal := elementPtr.Elem()

				newElement := reflect.ValueOf(expected)
				elementVal.Set(newElement)
				return nil
			},
		}

		var service ports.UserService
		service = NewUserService(&repo, config)

		_, err := service.Create(
			expected.Name,
			expected.Username,
			expected.Picture,
			expected.Role,
			expected.Provider,
			expected.TokenID,
			expected.Status,
		)

		if err == nil {
			t.Error("Item should not be created and return an error")
		}

		expectedError := "duplicated Username: CaptainAmerica"
		if err.Error() != expectedError {
			t.Errorf("Expected error: %q got: %q", expectedError, err.Error())
		}
	})
}

func TestReadOperations(t *testing.T) {
	config := domain.DefaultConfig()
	config.Keys = domain.KeyList{
		Auth: authKey,
	}

	t.Run("Test list users", func(t *testing.T) {
		expected := []domain.User{
			{
				Id: "1",
			},
			{
				Id: "2",
			},
		}

		dummyData := []map[string]interface{}{
			{
				"id": "1",
			},
			{
				"id": "2",
			},
		}

		data := map[string][]map[string]interface{}{
			"users": dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		service := NewUserService(&repo, config)

		got, err := service.List(nil, nil)

		if err != nil {
			t.Errorf("Got error while getting all users: %v", err)
		}

		if len(got) != len(expected) {
			t.Errorf("Expected %d elements got %d", len(expected), len(got))
		}

		for i := 0; i < len(got); i++ {
			if got[i].Id != expected[i].Id {
				t.Errorf("Expected first item id to be: %v got: %v", expected[0].Id, got[0].Id)
			}
		}
	})

	t.Run("Test getting user by id", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id": "1",
			},
			{
				"id": "2",
			},
		}

		data := map[string][]map[string]interface{}{
			"users": dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		service := NewUserService(&repo, config)

		got, err := service.Get("1")

		if err != nil {
			t.Errorf("Got error while getting user by id: %v", err)
		}

		if got.Id != "1" {
			t.Errorf("Expected id: %q got: %q", "1", got.Id)
		}
	})

	t.Run("Test getting user by username", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":       "1",
				"username": "IronMan",
			},
			{
				"id":       "2",
				"username": "CapAmerica",
			},
		}

		data := map[string][]map[string]interface{}{
			"users": dummyData,
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

		service := NewUserService(&repo, config)

		_, err := service.GetByUsername("IronMan")

		if err != nil {
			t.Errorf("Got error while getting user by id: %v", err)
		}

		if !called {
			t.Errorf("Expected GetOne to be called")
		}
	})

	t.Run("Test getting users by role", func(t *testing.T) {
		page := 1
		size := 10

		dummyData := []map[string]interface{}{
			{
				"id":       "1",
				"username": "IronMan",
				"role":     "genius",
			},
			{
				"id":       "2",
				"username": "CapAmerica",
				"role":     "leader",
			},
		}

		data := map[string][]map[string]interface{}{
			"users": dummyData,
		}

		called := false
		repo := mocks.MemRepo{
			Data: data,
			ListInterceptor: func(collection string, results interface{}, skip, limit int, filters ...ports.Filter) error {
				called = true
				if skip != 0 {
					t.Errorf("Expect skip to be 0 got %d", skip)
				}

				if limit != 10 {
					t.Errorf("Expect limit to be 10 got %d", limit)
				}

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

		service := NewUserService(&repo, config)

		_, err := service.ListByRole("genius", &page, &size)

		if err != nil {
			t.Errorf("Got error while getting user by id: %v", err)
		}

		if !called {
			t.Errorf("Expected List to be called")
		}
	})
}

func TestUpdateOperations(t *testing.T) {
	config := domain.DefaultConfig()
	config.Keys = domain.KeyList{
		Auth: authKey,
	}

	tokenId := "myTokenId"
	newTokenId := "newTokenId"
	encrypted, _ := utils.AES256Encrypt(authKey, tokenId)
	newTokenEncrypted, _ := utils.AES256Encrypt(authKey, newTokenId)
	now := utils.UnixUTCNow()
	yesterday := now.Add(-24 * time.Hour)
	t.Run("Test user is updated", func(t *testing.T) {
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
			"users": dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		service := NewUserService(&repo, config)

		expected := domain.User{
			Id:         "1",
			Username:   "CapAmerica",
			Name:       "Sam Wilson",
			Picture:    "",
			TokenID:    newTokenId,
			CreateDate: now,
			UpdateDate: now,
			Status:     "active",
		}

		_, err := service.Update(expected)

		if err != nil {
			t.Errorf("Item should be updated without errors: %v", err)
		}

		got, _ := service.Get("1")

		if got.Name != expected.Name {
			t.Errorf(
				"Name was not assigned expected: %q got: %q",
				expected.Name,
				got.Name,
			)
		}

		if got.Username != expected.Username {
			t.Errorf(
				"Username was not assigned expected: %q got: %q",
				expected.Username,
				got.Username,
			)
		}

		if got.Picture != expected.Picture {
			t.Errorf(
				"Picture was not assigned expected: %q got: %q",
				expected.Picture,
				got.Picture,
			)
		}

		if got.Provider != expected.Provider {
			t.Errorf(
				"Provider was not assigned expected: %q got: %q",
				expected.Provider,
				got.Provider,
			)
		}

		decryptedToken, _ := utils.AES256Decrypt(authKey, got.TokenID)
		if decryptedToken != newTokenId {
			t.Errorf(
				"TokenID was not assigned expected: %q got: %q",
				newTokenEncrypted,
				got.TokenID,
			)
		}

		if !got.CreateDate.Equal(yesterday) {
			t.Errorf(
				"CreateDate must be equals to yesterday: %q got %q",
				yesterday,
				got.CreateDate,
			)
		}

		if !got.UpdateDate.Equal(now) {
			t.Errorf(
				"UpdateDate must be after baseline now: %d got %d",
				now.UnixNano(),
				got.UpdateDate.UnixNano(),
			)
		}

		if got.Status != expected.Status {
			t.Errorf(
				"Status was not assigned expected: %q got: %q",
				expected.Status,
				got.Status,
			)
		}
	})

	t.Run("Test update a non-existing id", func(t *testing.T) {
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
			"users": dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		service := NewUserService(&repo, config)

		expected := domain.User{
			Id: "3",
		}

		_, err := service.Update(expected)

		if err == nil {
			t.Errorf("Expected error got nil")
		}

		_, ok := err.(ports.ErrItemNotFound)
		if !ok {
			t.Errorf("Expected error of type ErrItemNotFound got: %T", err)
		}
	})
}

func TestDeleteOperations(t *testing.T) {
	config := domain.DefaultConfig()
	config.Keys = domain.KeyList{
		Auth: authKey,
	}

	t.Run("Delete an item", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":         "1",
				"username":   "CapAmerica",
				"name":       "Steve Rogers",
				"picture":    "",
				"tokenID":    "",
				"createDate": utils.UnixUTCNow(),
				"updateDate": utils.UnixUTCNow(),
				"status":     "active",
			},
		}

		data := map[string][]map[string]interface{}{
			"users": dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		service := NewUserService(&repo, config)

		err := service.Delete("1", false)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
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

	t.Run("Delete a non-existing id", func(t *testing.T) {
		dummyData := []map[string]interface{}{
			{
				"id":         "1",
				"username":   "CapAmerica",
				"name":       "Steve Rogers",
				"picture":    "",
				"tokenID":    "",
				"createDate": utils.UnixUTCNow(),
				"updateDate": utils.UnixUTCNow(),
				"status":     "active",
			},
		}

		data := map[string][]map[string]interface{}{
			"users": dummyData,
		}

		repo := mocks.MemRepo{
			Data: data,
		}

		service := NewUserService(&repo, config)

		_ = service.Delete("3", false)

		got, err := service.Get("1")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if got.Username != "CapAmerica" {
			t.Errorf(
				"Expected username: %q got: %q",
				got.Username,
				"CapAmerica",
			)
		}
	})
}
