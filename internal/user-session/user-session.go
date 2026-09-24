package usersession

import "github.com/lemmyMwaura/pass/internal/vault"

// UserSession holds an unlocked vault for the interactive session.
type UserSession struct {
	Vault *vault.Vault
}

func NewUserSession(v *vault.Vault) *UserSession {
	return &UserSession{Vault: v}
}

func (s *UserSession) Close() {
	if s.Vault != nil {
		s.Vault.Lock()
		s.Vault = nil
	}
}
