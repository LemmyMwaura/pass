package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("entry not found")
	ErrAlreadyExists = errors.New("vault already exists")
	ErrNoVault       = errors.New("vault does not exist")
)

// Entry is a single stored credential.
type Entry struct {
	ID       string    `json:"id"`
	Service  string    `json:"service"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Notes    string    `json:"notes,omitempty"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type payload struct {
	Entries []Entry `json:"entries"`
}

// Vault is an unlocked in-memory password store.
type Vault struct {
	Username       string
	masterPassword string
	path           string
	entries        []Entry
}

func vaultDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".pass")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return dir, nil
}

func PathFor(username string) (string, error) {
	dir, err := vaultDir()
	if err != nil {
		return "", err
	}
	safe := strings.ReplaceAll(username, string(os.PathSeparator), "_")
	return filepath.Join(dir, safe+".vault"), nil
}

func Exists(username string) (bool, error) {
	path, err := PathFor(username)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return err == nil, err
}

// Create initializes a new empty vault for username.
func Create(username, masterPassword string) (*Vault, error) {
	exists, err := Exists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAlreadyExists
	}

	path, err := PathFor(username)
	if err != nil {
		return nil, err
	}

	v := &Vault{
		Username:       username,
		masterPassword: masterPassword,
		path:           path,
		entries:        []Entry{},
	}
	if err := v.Save(); err != nil {
		return nil, err
	}
	return v, nil
}

// Unlock decrypts an existing vault with the master password.
func Unlock(username, masterPassword string) (*Vault, error) {
	path, err := PathFor(username)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoVault
		}
		return nil, err
	}

	plaintext, err := unseal(data, masterPassword)
	if err != nil {
		return nil, err
	}

	var p payload
	if err := json.Unmarshal(plaintext, &p); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorruptVault, err)
	}

	return &Vault{
		Username:       username,
		masterPassword: masterPassword,
		path:           path,
		entries:        p.Entries,
	}, nil
}

func (v *Vault) Save() error {
	plaintext, err := json.Marshal(payload{Entries: v.entries})
	if err != nil {
		return err
	}

	sealed, err := seal(plaintext, v.masterPassword)
	if err != nil {
		return err
	}

	tmp := v.path + ".tmp"
	if err := os.WriteFile(tmp, sealed, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, v.path)
}

func (v *Vault) List() []Entry {
	out := make([]Entry, len(v.entries))
	copy(out, v.entries)
	return out
}

func (v *Vault) Get(service string) (*Entry, error) {
	service = strings.ToLower(strings.TrimSpace(service))
	for i := range v.entries {
		if strings.ToLower(v.entries[i].Service) == service {
			e := v.entries[i]
			return &e, nil
		}
	}
	return nil, ErrNotFound
}

func (v *Vault) Add(service, username, password, notes string) (Entry, error) {
	service = strings.TrimSpace(service)
	if service == "" {
		return Entry{}, errors.New("service name is required")
	}

	if _, err := v.Get(service); err == nil {
		return Entry{}, fmt.Errorf("entry for %q already exists", service)
	}

	now := time.Now().UTC()
	entry := Entry{
		ID:       uuid.NewString(),
		Service:  service,
		Username: username,
		Password: password,
		Notes:    notes,
		Created:  now,
		Updated:  now,
	}
	v.entries = append(v.entries, entry)
	if err := v.Save(); err != nil {
		v.entries = v.entries[:len(v.entries)-1]
		return Entry{}, err
	}
	return entry, nil
}

func (v *Vault) Delete(service string) error {
	service = strings.ToLower(strings.TrimSpace(service))
	for i := range v.entries {
		if strings.ToLower(v.entries[i].Service) == service {
			v.entries = append(v.entries[:i], v.entries[i+1:]...)
			return v.Save()
		}
	}
	return ErrNotFound
}

// Lock clears sensitive material from memory as best-effort.
func (v *Vault) Lock() {
	v.masterPassword = ""
	for i := range v.entries {
		v.entries[i].Password = ""
	}
	v.entries = nil
}
