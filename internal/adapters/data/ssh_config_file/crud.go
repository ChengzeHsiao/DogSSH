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
	"fmt"
	"strings"

	"github.com/ChengzeHsiao/dogssh/internal/core/domain"
	"github.com/kevinburke/ssh_config"
)

const (
	MaxBackups         = 10
	TempSuffix         = ".tmp"
	BackupSuffix       = "dogssh.backup"
	SSHConfigPerms     = 0o600
	OriginalBackupName = "config.original.backup"
)

// filterServers filters servers based on the query string.
func (r *Repository) filterServers(servers []domain.Server, query string) []domain.Server {
	query = strings.ToLower(query)
	filtered := make([]domain.Server, 0)

	for _, server := range servers {
		if r.matchesQuery(server, query) {
			filtered = append(filtered, server)
		}
	}

	return filtered
}

// matchesQuery checks if any field of the server matches the query string.
func (r *Repository) matchesQuery(server domain.Server, query string) bool {
	fields := []string{
		strings.ToLower(server.Host),
		strings.ToLower(server.User),
	}
	for _, tag := range server.Tags {
		fields = append(fields, strings.ToLower(tag))
	}
	for _, alias := range server.Aliases {
		fields = append(fields, strings.ToLower(alias))
	}

	for _, field := range fields {
		if strings.Contains(field, query) {
			return true
		}
	}

	return false
}

// serverExists checks if a server with the given alias already exists in the config.
func (r *Repository) serverExists(cfg *ssh_config.Config, alias string) bool {
	return r.findHostByAlias(cfg, alias) != nil
}

// findHostByAlias finds a host by its alias in the SSH config.
func (r *Repository) findHostByAlias(cfg *ssh_config.Config, alias string) *ssh_config.Host {
	for _, host := range cfg.Hosts {
		if r.hostContainsPattern(host, alias) {
			return host
		}
	}
	return nil
}

// hostContainsPattern checks if a host contains a specific pattern.
func (r *Repository) hostContainsPattern(host *ssh_config.Host, target string) bool {
	for _, pattern := range host.Patterns {
		if pattern.String() == target {
			return true
		}
	}
	return false
}

// createHostFromServer creates a new ssh_config.Host from a domain.Server.
func (r *Repository) createHostFromServer(server domain.Server) *ssh_config.Host {
	host := &ssh_config.Host{
		Patterns: []*ssh_config.Pattern{
			{Str: server.Alias},
		},
		Nodes:              make([]ssh_config.Node, 0),
		LeadingSpace:       4,
		EOLComment:         "Added by dogssh",
		SpaceBeforeComment: strings.Repeat(" ", 4),
	}

	r.addKVNodeIfNotEmpty(host, "HostName", server.Host)
	r.addKVNodeIfNotEmpty(host, "User", server.User)
	r.addKVNodeIfNotEmpty(host, "Port", fmt.Sprintf("%d", server.Port))
	for _, identityFile := range server.IdentityFiles {
		r.addKVNodeIfNotEmpty(host, "IdentityFile", identityFile)
	}

	return host
}

// addKVNodeIfNotEmpty adds a key-value node to the host if the value is not empty.
func (r *Repository) addKVNodeIfNotEmpty(host *ssh_config.Host, key, value string) {
	if value == "" {
		return
	}

	kvNode := &ssh_config.KV{
		Key:          key,
		Value:        value,
		LeadingSpace: 4,
	}
	host.Nodes = append(host.Nodes, kvNode)
}

// updateHostNodes updates the nodes of an existing host with new server details.
func (r *Repository) updateHostNodes(host *ssh_config.Host, newServer domain.Server) {
	updates := map[string]string{
		"hostname": newServer.Host,
		"user":     newServer.User,
		"port":     fmt.Sprintf("%d", newServer.Port),
	}
	for key, value := range updates {
		if value != "" {
			r.updateOrAddKVNode(host, key, value)
		}
	}
	// Replace IdentityFile entries entirely to reflect the new state.
	// This ensures removing/clearing identity files works as expected.

	removeKey := func(nodes []ssh_config.Node, key string) []ssh_config.Node {
		filtered := make([]ssh_config.Node, 0, len(nodes))
		for _, node := range nodes {
			if kv, ok := node.(*ssh_config.KV); ok {
				if strings.EqualFold(kv.Key, key) {
					continue // skip existing IdentityFile
				}
			}
			filtered = append(filtered, node)
		}
		return filtered
	}
	host.Nodes = removeKey(host.Nodes, "IdentityFile")

	for _, identityFile := range newServer.IdentityFiles {
		r.addKVNodeIfNotEmpty(host, "IdentityFile", identityFile)
	}
}

// updateOrAddKVNode updates an existing key-value node or adds a new one if it doesn't exist.
func (r *Repository) updateOrAddKVNode(host *ssh_config.Host, key, newValue string) {
	keyLower := strings.ToLower(key)

	// Try to update existing node
	for _, node := range host.Nodes {
		kvNode, ok := node.(*ssh_config.KV)
		if ok && strings.EqualFold(kvNode.Key, keyLower) {
			kvNode.Value = newValue
			return
		}
	}

	// Add new node if not found
	kvNode := &ssh_config.KV{
		Key:          r.getProperKeyCase(key),
		Value:        newValue,
		LeadingSpace: 4,
	}
	host.Nodes = append(host.Nodes, kvNode)
}

// getProperKeyCase returns the proper case for known SSH config keys.
// Reference: https://www.ssh.com/academy/ssh/config
func (r *Repository) getProperKeyCase(key string) string {
	keyMap := map[string]string{
		"hostname":     "HostName",
		"user":         "User",
		"port":         "Port",
		"identityfile": "IdentityFile",
	}

	if properCase, exists := keyMap[strings.ToLower(key)]; exists {
		return properCase
	}
	return key
}

// removeHostByAlias removes a host by its alias from the list of hosts.
func (r *Repository) removeHostByAlias(hosts []*ssh_config.Host, alias string) []*ssh_config.Host {
	for i, host := range hosts {
		if r.hostContainsPattern(host, alias) {
			return append(hosts[:i], hosts[i+1:]...)
		}
	}
	return hosts
}
