package svc

import (
	"cloudflared-tunnel/ent"
	v1 "cloudflared-tunnel/internal/module/user/ui/api/req/v1"
	"cloudflared-tunnel/pkg/errno"

	"golang.org/x/crypto/bcrypt"
)

func (s *svc) Register(nickname, username, password, email string) (*v1.UserVO, error) {
	count, err := s.repo.GetUserCount()
	if err != nil {
		return nil, errno.ErrDB
	}
	if count > 0 {
		return nil, errno.ErrUserExists
	}

	existingUser, _ := s.repo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errno.ErrUserExists.WithMessage("用户名已存在")
	}

	existingUser, _ = s.repo.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errno.ErrUserExists.WithMessage("邮箱已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("密码加密失败", "error", err)
		return nil, errno.ErrInternal
	}

	u, err := s.repo.CreateUser(nickname, username, string(hashedPassword), email)
	if err != nil {
		return nil, errno.ErrDB
	}

	return &v1.UserVO{
		ID:       u.ID,
		Nickname: u.Nickname,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (s *svc) GetUserByID(id string) (*v1.UserVO, error) {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDB
	}

	return &v1.UserVO{
		ID:       u.ID,
		Nickname: u.Nickname,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (s *svc) Login(username, password string) (*v1.LoginVO, error) {
	u, err := s.repo.GetUserByUsername(username)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDB
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errno.ErrParam.WithMessage("密码错误")
	}

	token, err := s.jwt.GenerateToken(u.ID, u.Username)
	if err != nil {
		s.log.Error("生成token失败", "error", err)
		return nil, errno.ErrInternal
	}

	return &v1.LoginVO{
		Token: token,
		User: v1.UserVO{
			ID:       u.ID,
			Nickname: u.Nickname,
			Username: u.Username,
			Email:    u.Email,
		},
	}, nil
}
