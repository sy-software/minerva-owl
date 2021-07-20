package utils

import "testing"

func TestAES256(t *testing.T) {
	key := "2b7e151628aed2a6abf71589a12b4da32"
	input := "Hello World"

	encrypted, err := AES256Encrypt(key, input)

	if err != nil {
		t.Errorf("Item should be encrypted without errors: %v", err)
	}

	if encrypted == input {
		t.Errorf("Encrypted value %q can't be same as input %q", encrypted, input)
	}

	decrypted, err := AES256Decrypt(key, encrypted)

	if decrypted != input {
		t.Errorf("Expected %q Got %q after decrypt", input, decrypted)
	}
}
