package repo

import (
	"cloudflared-tunnel/ent"
	"cloudflared-tunnel/ent/credential"
	"cloudflared-tunnel/ent/credentialtestlog"
	"cloudflared-tunnel/ent/user"
	"context"
)

func (r *Repo) CreateCredential(userID int64, name, encryptedToken, accountID string, isDefault bool) (*ent.Credential, error) {
	ctx := context.Background()
	c, err := r.Client.Credential.Create().
		SetName(name).
		SetAPIToken(encryptedToken).
		SetAccountID(accountID).
		SetIsDefault(isDefault).
		SetOwnerID(userID).
		Save(ctx)
	if err != nil {
		r.Log.Error("创建凭证失败", "userID", userID, "name", name, "error", err)
		return nil, err
	}
	return c, nil
}

func (r *Repo) GetCredentialByID(id int64) (*ent.Credential, error) {
	ctx := context.Background()
	c, err := r.Client.Credential.Get(ctx, id)
	if err != nil {
		r.Log.Error("查询凭证失败", "id", id, "error", err)
		return nil, err
	}
	return c, nil
}

func (r *Repo) GetCredentialByIDAndUserID(id, userID int64) (*ent.Credential, error) {
	ctx := context.Background()
	c, err := r.Client.Credential.Query().
		Where(
			credential.ID(id),
			credential.HasOwnerWith(user.ID(userID)),
		).
		Only(ctx)
	if err != nil {
		r.Log.Error("查询凭证失败", "id", id, "userID", userID, "error", err)
		return nil, err
	}
	return c, nil
}

func (r *Repo) GetCredentialsByUserID(userID int64) ([]*ent.Credential, error) {
	ctx := context.Background()
	credentials, err := r.Client.Credential.Query().
		Where(credential.HasOwnerWith(user.ID(userID))).
		Order(ent.Asc(credential.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		r.Log.Error("查询用户凭证列表失败", "userID", userID, "error", err)
		return nil, err
	}
	return credentials, nil
}

func (r *Repo) GetDefaultCredentialByUserID(userID int64) (*ent.Credential, error) {
	ctx := context.Background()
	c, err := r.Client.Credential.Query().
		Where(
			credential.IsDefault(true),
			credential.HasOwnerWith(user.ID(userID)),
		).
		Only(ctx)
	if err != nil {
		r.Log.Error("查询默认凭证失败", "userID", userID, "error", err)
		return nil, err
	}
	return c, nil
}

func (r *Repo) UpdateCredential(id int64, name, encryptedToken, accountID string, isDefault bool) (*ent.Credential, error) {
	ctx := context.Background()
	update := r.Client.Credential.UpdateOneID(id).
		SetName(name).
		SetAPIToken(encryptedToken).
		SetAccountID(accountID).
		SetIsDefault(isDefault)

	c, err := update.Save(ctx)
	if err != nil {
		r.Log.Error("更新凭证失败", "id", id, "error", err)
		return nil, err
	}
	return c, nil
}

func (r *Repo) DeleteCredential(id, userID int64) error {
	ctx := context.Background()
	_, err := r.Client.Credential.Delete().
		Where(
			credential.ID(id),
			credential.HasOwnerWith(user.ID(userID)),
		).
		Exec(ctx)
	if err != nil {
		r.Log.Error("删除凭证失败", "id", id, "userID", userID, "error", err)
		return err
	}
	return nil
}

func (r *Repo) ClearDefaultByUserID(userID int64) error {
	ctx := context.Background()
	_, err := r.Client.Credential.Update().
		Where(
			credential.IsDefault(true),
			credential.HasOwnerWith(user.ID(userID)),
		).
		SetIsDefault(false).
		Save(ctx)
	if err != nil {
		r.Log.Error("清除默认凭证失败", "userID", userID, "error", err)
		return err
	}
	return nil
}

func (r *Repo) CreateTestLog(credentialID int64, status string, errMsg *string) (*ent.CredentialTestLog, error) {
	ctx := context.Background()
	create := r.Client.CredentialTestLog.Create().
		SetStatus(status).
		SetCredentialID(credentialID)
	if errMsg != nil {
		create = create.SetErrorMessage(*errMsg)
	}
	log, err := create.Save(ctx)
	if err != nil {
		r.Log.Error("创建测试日志失败", "credentialID", credentialID, "error", err)
		return nil, err
	}
	return log, nil
}

func (r *Repo) GetTestLogsByCredentialID(credentialID int64, limit int) ([]*ent.CredentialTestLog, error) {
	ctx := context.Background()
	logs, err := r.Client.CredentialTestLog.Query().
		Where(credentialtestlog.HasCredentialWith(credential.ID(credentialID))).
		Order(ent.Desc(credentialtestlog.FieldTestedAt)).
		Limit(limit).
		All(ctx)
	if err != nil {
		r.Log.Error("查询测试日志失败", "credentialID", credentialID, "error", err)
		return nil, err
	}
	return logs, nil
}
