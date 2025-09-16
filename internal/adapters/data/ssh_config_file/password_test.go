// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh_config_file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Adembc/lazyssh/internal/core/domain"
	"go.uber.org/zap"
)

func TestPasswordSaving(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "lazyssh_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	// Create logger
	logger := zap.NewNop().Sugar()

	// Create password manager with temp file
	passwordFile := filepath.Join(tempDir, "passwords.json")
	pm := NewPasswordManager(passwordFile, logger)

	// Test server
	server := domain.Server{
		Alias: "test-server",
		Host:  "example.com",
		User:  "root",
		Port:  22,
	}

	// Test password saving
	testPassword := "my-secret-password"
	err = pm.UpdateServerPassword(server, testPassword)
	if err != nil {
		t.Fatalf("Failed to save password: %v", err)
	}

	// Test password retrieval and encryption
	encryptedPassword, err := pm.GetServerPassword(server.Alias)
	if err != nil {
		t.Fatalf("Failed to get encrypted password: %v", err)
	}

	// Test password decryption
	decryptedPassword, err := pm.DecryptPassword(encryptedPassword)
	if err != nil {
		t.Fatalf("Failed to decrypt password: %v", err)
	}

	// Verify decrypted password matches original
	if decryptedPassword != testPassword {
		t.Fatalf("Decrypted password doesn't match original. Got: %s, Expected: %s", decryptedPassword, testPassword)
	}

	t.Logf("Password saving and verification working correctly")
}
