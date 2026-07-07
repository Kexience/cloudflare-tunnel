package testutil

import (
	"context"
	"testing"

	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/ent/enttest"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteContainer struct {
	Client *ent.Client
}

func NewSQLiteContainer(t *testing.T) *SQLiteContainer {
	t.Helper()

	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	t.Cleanup(func() {
		client.Close()
	})

	return &SQLiteContainer{
		Client: client,
	}
}

func (s *SQLiteContainer) InsertUser(t *testing.T, nickname, username, password, email string) *ent.User {
	t.Helper()
	ctx := context.Background()
	u, err := s.Client.User.Create().
		SetNickname(nickname).
		SetUsername(username).
		SetPassword(password).
		SetEmail(email).
		Save(ctx)
	if err != nil {
		t.Fatalf("插入用户失败: %v", err)
	}
	return u
}

func (s *SQLiteContainer) InsertCredential(t *testing.T, userID int64, name, apiToken, accountID string, isDefault bool) *ent.Credential {
	t.Helper()
	ctx := context.Background()
	c, err := s.Client.Credential.Create().
		SetName(name).
		SetAPIToken(apiToken).
		SetAccountID(accountID).
		SetIsDefault(isDefault).
		SetOwnerID(userID).
		Save(ctx)
	if err != nil {
		t.Fatalf("插入凭证失败: %v", err)
	}
	return c
}
