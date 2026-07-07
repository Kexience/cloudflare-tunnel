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
		return fmt.Errorf("创建 Cloudflare 客户端失败")
	}

	ctx := context.Background()

	tokenInfo, err := client.VerifyAPIToken(ctx)
	if err != nil {
		return fmt.Errorf("API Token 无效")
	}

	if tokenInfo.Status != "active" {
		return fmt.Errorf("API Token 状态异常: %s", tokenInfo.Status)
	}

	return nil
}
