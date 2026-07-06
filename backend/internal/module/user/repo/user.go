package repo

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/ent/user"
	"cloudflared-tunnel/internal/infra/logger"
	"context"
	"strconv"
)

func NewUserRepo(client *ent.Client, log logger.Logger) UserRepo {
	return &Repo{
		Client: client,
		Log:    log,
	}
}

// GetUserByEmail implements [UserRepo].
func (r *Repo) GetUserByEmail(email string) (*ent.User, error) {
	ctx := context.Background()
	u, err := r.Client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		r.Log.Error("根据邮箱查询用户失败", "email", email, "error", err)
		return nil, err
	}
	return u, nil
}

// GetUserByID implements [UserRepo].
func (r *Repo) GetUserByID(id string) (*ent.User, error) {
	ctx := context.Background()
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		r.Log.Error("根据ID查询用户失败：无效的ID", "id", id, "error", err)
		return nil, err
	}
	u, err := r.Client.User.Get(ctx, idInt)
	if err != nil {
		r.Log.Error("根据ID查询用户失败", "id", id, "error", err)
		return nil, err
	}
	return u, nil
}

// GetUserByUsername implements [UserRepo].
func (r *Repo) GetUserByUsername(username string) (*ent.User, error) {
	ctx := context.Background()
	u, err := r.Client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	if err != nil {
		r.Log.Error("根据用户名查询用户失败", "username", username, "error", err)
		return nil, err
	}
	return u, nil
}

// GetUserCount implements [UserRepo].
func (r *Repo) GetUserCount() (int, error) {
	ctx := context.Background()
	count, err := r.Client.User.Query().Count(ctx)
	if err != nil {
		r.Log.Error("查询用户数量失败", "error", err)
		return 0, err
	}
	return count, nil
}

// CreateUser implements [UserRepo].
func (r *Repo) CreateUser(nickname, username, password, email string) (*ent.User, error) {
	ctx := context.Background()
	u, err := r.Client.User.Create().
		SetNickname(nickname).
		SetUsername(username).
		SetPassword(password).
		SetEmail(email).
		Save(ctx)
	if err != nil {
		r.Log.Error("创建用户失败", "username", username, "email", email, "error", err)
		return nil, err
	}
	return u, nil
}
