package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/infra/logger"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	_ "github.com/lib/pq"
)

func NewLogger(t *testing.T) logger.Logger {
	t.Helper()
	log, err := logger.NewLoggerForTest()
	if err != nil {
		t.Fatalf("创建日志失败: %v", err)
	}
	return log
}

type PostgresContainer struct {
	Container testcontainers.Container
	Client    *ent.Client
}

func NewPostgresContainer(t *testing.T) *PostgresContainer {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("启动 PostgreSQL 容器失败: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("获取容器主机失败: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("获取容器端口失败: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=test password=test dbname=testdb sslmode=disable", host, port.Port())
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("创建 ent client 失败: %v", err)
	}

	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("运行数据库迁移失败: %v", err)
	}

	t.Cleanup(func() {
		client.Close()
		container.Terminate(context.Background())
	})

	return &PostgresContainer{
		Container: container,
		Client:    client,
	}
}

func (p *PostgresContainer) InsertUser(t *testing.T, nickname, username, password, email string) *ent.User {
	t.Helper()
	ctx := context.Background()
	u, err := p.Client.User.Create().
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

func (p *PostgresContainer) InsertCredential(t *testing.T, userID int64, name, apiToken, accountID string, isDefault bool) *ent.Credential {
	t.Helper()
	ctx := context.Background()
	c, err := p.Client.Credential.Create().
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
