package cloudflare

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go"
)

type Validator interface {
	Validate(apiToken, accountID string) error
}

type validator struct{}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) Validate(apiToken, accountID string) error {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()

	_, err = client.UserDetails(ctx)
	if err != nil {
		return fmt.Errorf("API Token 无效: %w", err)
	}

	_, _, err = client.Account(ctx, accountID)
	if err != nil {
		return fmt.Errorf("账号 ID 无效或无权访问: %w", err)
	}

	return nil
}
