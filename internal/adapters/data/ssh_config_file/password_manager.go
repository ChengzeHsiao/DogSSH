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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Adembc/lazyssh/internal/core/domain"
	"go.uber.org/zap"
)

// PasswordManager handles encrypted storage and retrieval of server passwords
type PasswordManager struct {
	filePath string
	logger   *zap.SugaredLogger
}

// NewPasswordManager creates a new password manager instance
func NewPasswordManager(filePath string, logger *zap.SugaredLogger) *PasswordManager {
	return &PasswordManager{filePath: filePath, logger: logger}
}

// loadPasswords 从文件中加载所有密码
func (p *PasswordManager) loadPasswords() (map[string]string, error) {
	passwords := make(map[string]string)

	if _, err := os.Stat(p.filePath); os.IsNotExist(err) {
		return passwords, nil
	}

	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("read passwords '%s': %w", p.filePath, err)
	}

	if len(data) == 0 {
		return passwords, nil
	}

	if err := json.Unmarshal(data, &passwords); err != nil {
		return nil, fmt.Errorf("parse passwords JSON '%s': %w", p.filePath, err)
	}

	return passwords, nil
}

// savePasswords 将所有密码保存到文件
func (p *PasswordManager) savePasswords(passwords map[string]string) error {
	if err := p.ensureDirectory(); err != nil {
		p.logger.Errorw("failed to ensure passwords directory", "path", p.filePath, "error", err)
		return fmt.Errorf("ensure passwords directory for '%s': %w", p.filePath, err)
	}

	data, err := json.MarshalIndent(passwords, "", "  ")
	if err != nil {
		p.logger.Errorw("failed to marshal passwords", "path", p.filePath, "error", err)
		return fmt.Errorf("marshal passwords for '%s': %w", p.filePath, err)
	}

	// Save password file with 0600 permissions (owner read/write only)
	if err := os.WriteFile(p.filePath, data, 0o600); err != nil {
		p.logger.Errorw("failed to write passwords file", "path", p.filePath, "error", err)
		return fmt.Errorf("write passwords '%s': %w", p.filePath, err)
	}
	return nil
}

// getEncryptionKey derives a consistent key from a password for AES encryption
func (p *PasswordManager) getEncryptionKey() []byte {
	// Use a combination of file path and a static string to create a consistent key
	// In production, this should use a proper key derivation function with salt
	keyMaterial := p.filePath + "lazyssh-password-encryption-key"
	hash := sha256.Sum256([]byte(keyMaterial))
	return hash[:]
}

// EncryptPassword encrypts a password using AES
func (p *PasswordManager) EncryptPassword(password string) (string, error) {
	key := p.getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a random nonce
	nonce := make([]byte, 12) // GCM nonce size
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(password), nil)
	// Prepend nonce to ciphertext
	encrypted := make([]byte, 0, len(nonce)+len(ciphertext))
	encrypted = append(encrypted, nonce...)
	encrypted = append(encrypted, ciphertext...)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptPassword decrypts an encrypted password
func (p *PasswordManager) DecryptPassword(encryptedPassword string) (string, error) {
	key := p.getEncryptionKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	encrypted, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	if len(encrypted) < 12 {
		return "", fmt.Errorf("encrypted password too short")
	}

	nonce := encrypted[:12]
	ciphertext := encrypted[12:]

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// UpdateServerPassword updates server password
func (p *PasswordManager) UpdateServerPassword(server domain.Server, newPassword string) error {
	// If password is empty, don't update
	if newPassword == "" {
		return nil
	}

	passwords, err := p.loadPasswords()
	if err != nil {
		p.logger.Errorw("failed to load passwords in UpdateServerPassword", "path", p.filePath, "alias", server.Alias, "error", err)
		return fmt.Errorf("load passwords: %w", err)
	}

	// Encrypt the new password
	encryptedPassword, err := p.EncryptPassword(newPassword)
	if err != nil {
		p.logger.Errorw("failed to encrypt password in UpdateServerPassword", "alias", server.Alias, "error", err)
		return fmt.Errorf("encrypt password: %w", err)
	}

	// Save the encrypted password
	passwords[server.Alias] = encryptedPassword
	return p.savePasswords(passwords)
}

// GetServerPassword 获取服务器的密码哈希值
func (p *PasswordManager) GetServerPassword(alias string) (string, error) {
	passwords, err := p.loadPasswords()
	if err != nil {
		p.logger.Errorw("failed to load passwords in GetServerPassword", "path", p.filePath, "alias", alias, "error", err)
		return "", fmt.Errorf("load passwords: %w", err)
	}

	hashedPassword, exists := passwords[alias]
	if !exists {
		return "", fmt.Errorf("password for server '%s' not found", alias)
	}

	return hashedPassword, nil
}

// DeleteServerPassword 删除服务器的密码
func (p *PasswordManager) DeleteServerPassword(alias string) error {
	passwords, err := p.loadPasswords()
	if err != nil {
		p.logger.Errorw("failed to load passwords in DeleteServerPassword", "path", p.filePath, "alias", alias, "error", err)
		return fmt.Errorf("load passwords: %w", err)
	}

	delete(passwords, alias)
	return p.savePasswords(passwords)
}

// ensureDirectory ensures the directory for storing passwords exists
func (p *PasswordManager) ensureDirectory() error {
	dir := filepath.Dir(p.filePath)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("mkdir '%s': %w", dir, err)
	}
	return nil
}
