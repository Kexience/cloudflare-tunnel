package repo

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/internal/infra/logger"
)

type Repo struct {
	Client *ent.Client
	Log    logger.Logger
}

type UserRepo interface {
	GetUserByID(id string) (*ent.User, error)
	GetUserByEmail(email string) (*ent.User, error)
	GetUserByUsername(username string) (*ent.User, error)
	GetUserCount() (int, error)
	CreateUser(nickname, username, password, email string) (*ent.User, error)
}
