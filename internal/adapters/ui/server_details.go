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

package ui

import (
	"fmt"
	"strings"

	"github.com/Adembc/dogssh/internal/core/domain"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ServerDetails struct {
	*tview.TextView
}

func NewServerDetails() *ServerDetails {
	details := &ServerDetails{
		TextView: tview.NewTextView(),
	}
	details.build()
	return details
}

func (sd *ServerDetails) build() {
	sd.TextView.SetDynamicColors(true).
		SetWrap(true).
		SetBorder(true).
		SetTitle("Details").
		SetBorderColor(tcell.Color238).
		SetTitleColor(tcell.Color250)
}

// renderTagChips builds colored tag chips for details view.
func renderTagChips(tags []string) string {
	if len(tags) == 0 {
		return "-"
	}
	chips := make([]string, 0, len(tags))
	for _, t := range tags {
		chips = append(chips, fmt.Sprintf("[black:#5FAFFF] %s [-:-:-]", t))
	}
	return strings.Join(chips, " ")
}

// UpdateServer updates the details view with the provided server information.
func (sd *ServerDetails) UpdateServer(server domain.Server) {
	lastSeen := server.LastSeen.Format("2006-01-02 15:04:05")
	if server.LastSeen.IsZero() {
		lastSeen = "Never"
	}
	serverKey := strings.Join(server.IdentityFiles, ", ")

	pinnedStr := "true"
	if server.PinnedAt.IsZero() {
		pinnedStr = "false"
	}
	tagsText := renderTagChips(server.Tags)

	// 显示密码状态而不是明文密码
	passwordStatus := "Not set"

	text := fmt.Sprintf(
		"[::b]%s[-]\n\nHost: [white]%s[-]\nUser: [white]%s[-]\nPort: [white]%d[-]\nKey:  [white]%s[-]\nPassword: [white]%s[-]\nTags: %s\nPinned: [white]%s[-]\nLast SSH: %s\nSSH Count: [white]%d[-]\n\n[::b]Commands:[-]\n  Enter: SSH connect\n  c: Copy SSH command\n  g: Ping server\n  r: Refresh list\n  a: Add new server\n  e: Edit entry\n  t: Edit tags\n  d: Delete entry\n  p: Pin/Unpin",
		strings.Join(server.Aliases, ", "), server.Host, server.User, server.Port,
		serverKey, passwordStatus, tagsText, pinnedStr,
		lastSeen, server.SSHCount)
	sd.TextView.SetText(text)
}

// UpdateServerWithPasswordCheck updates the details view with the provided server information.
// It also checks if a password is stored for the server and displays the appropriate status.
func (sd *ServerDetails) UpdateServerWithPasswordCheck(server domain.Server, hasPassword bool) {
	lastSeen := server.LastSeen.Format("2006-01-02 15:04:05")
	if server.LastSeen.IsZero() {
		lastSeen = "Never"
	}
	serverKey := strings.Join(server.IdentityFiles, ", ")

	pinnedStr := "true"
	if server.PinnedAt.IsZero() {
		pinnedStr = "false"
	}
	tagsText := renderTagChips(server.Tags)

	// 显示密码状态而不是明文密码
	passwordStatus := "Not set"
	if hasPassword {
		passwordStatus = "Set (hidden)"
	}

	text := fmt.Sprintf(
		"[::b]%s[-]\n\nHost: [white]%s[-]\nUser: [white]%s[-]\nPort: [white]%d[-]\nKey:  [white]%s[-]\nPassword: [white]%s[-]\nTags: %s\nPinned: [white]%s[-]\nLast SSH: %s\nSSH Count: [white]%d[-]\n\n[::b]Commands:[-]\n  Enter: SSH connect\n  c: Copy SSH command\n  g: Ping server\n  r: Refresh list\n  a: Add new server\n  e: Edit entry\n  t: Edit tags\n  d: Delete entry\n  p: Pin/Unpin",
		strings.Join(server.Aliases, ", "), server.Host, server.User, server.Port,
		serverKey, passwordStatus, tagsText, pinnedStr,
		lastSeen, server.SSHCount)
	sd.TextView.SetText(text)
}

func (sd *ServerDetails) ShowEmpty() {
	sd.TextView.SetText("No servers match the current filter.")
}
