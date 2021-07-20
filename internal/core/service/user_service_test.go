package service

import (
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
		now := utils.UnixNow()

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
				now.Unix(),
				created.CreateDate.Unix(),
			)
		}

		if !created.UpdateDate.After(now) && !created.UpdateDate.Equal(now) {
			t.Errorf(
				"UpdateDate must be after or equal to baseline now: %d got %d",
				now.Unix(),
				created.UpdateDate.Unix(),
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

	// TODO: Implement and test unique username
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
}

func TestUpdateOperations(t *testing.T) {
	config := domain.DefaultConfig()
	config.Keys = domain.KeyList{
		Auth: authKey,
	}

	tokenId := "myTokenId"
	encrypted, _ := utils.AES256Encrypt(authKey, tokenId)
	now := time.Unix(time.Now().Unix(), 0)
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
			TokenID:    encrypted,
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

		if got.TokenID != expected.TokenID {
			t.Errorf(
				"TokenID was not assigned expected: %q got: %q",
				expected.Provider,
				got.Provider,
			)
		}

		if got.CreateDate != yesterday {
			t.Errorf(
				"CreateDate must be equals to now: %d got %d",
				yesterday.Unix(),
				got.CreateDate.Unix(),
			)
		}

		if !got.UpdateDate.Equal(now) {
			t.Errorf(
				"UpdateDate must be after baseline now: %d got %d",
				now.Unix(),
				got.UpdateDate.Unix(),
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
}
