package vault

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVaultRoundTrip(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	v, err := Create("alice", "s3cret-master")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if _, err := v.Add("github", "alice", "gh-token", "work"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	if _, err := v.Add("email", "alice@example.com", "mail-pass", ""); err != nil {
		t.Fatalf("Add email: %v", err)
	}

	v.Lock()

	unlocked, err := Unlock("alice", "s3cret-master")
	if err != nil {
		t.Fatalf("Unlock: %v", err)
	}

	entries := unlocked.List()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	got, err := unlocked.Get("GitHub")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Password != "gh-token" || got.Username != "alice" {
		t.Fatalf("unexpected entry: %+v", got)
	}

	if err := unlocked.Delete("email"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if len(unlocked.List()) != 1 {
		t.Fatalf("expected 1 entry after delete")
	}

	if _, err := Unlock("alice", "wrong"); err != ErrWrongPassword {
		t.Fatalf("expected ErrWrongPassword, got %v", err)
	}

	path, err := PathFor("alice")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Dir(path) != filepath.Join(home, ".pass") {
		t.Fatalf("unexpected vault path: %s", path)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		t.Fatalf("vault file too permissive: %v", info.Mode())
	}
}

func TestCreateDuplicate(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if _, err := Create("bob", "pw"); err != nil {
		t.Fatal(err)
	}
	if _, err := Create("bob", "pw"); err != ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}
