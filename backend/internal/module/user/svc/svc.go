package svc

import (
	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/internal/module/user/repo"
	v1 "cloudflared-tunnel/internal/module/user/ui/api/req/v1"
	"cloudflared-tunnel/pkg/core"
)

type UserSvc interface {
	Register(nickname, username, password, email string) (*v1.UserVO, error)
	GetUserByID(id string) (*v1.UserVO, error)
	Login(username, password string) (*v1.LoginVO, error)
}

type svc struct {
	repo repo.UserRepo
	log  logger.Logger
	jwt  *core.JWT
}

func NewUserSvc(repo repo.UserRepo, log logger.Logger, jwt *core.JWT) UserSvc {
	return &svc{
		repo: repo,
		log:  log,
		jwt:  jwt,
	}
}
