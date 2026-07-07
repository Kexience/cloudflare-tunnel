package cloudflare

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go"
)

// DNSClient Cloudflare DNS 记录管理 API 客户端封装
type DNSClient interface {
	// ListDNSRecords 查询 DNS 记录列表
	ListDNSRecords(apiToken, zoneID string, params cf.ListDNSRecordsParams) ([]cf.DNSRecord, error)
	// CreateDNSRecord 创建 DNS 记录
	CreateDNSRecord(apiToken, zoneID string, params cf.CreateDNSRecordParams) (cf.DNSRecord, error)
	// UpdateDNSRecord 更新 DNS 记录
	UpdateDNSRecord(apiToken, zoneID string, params cf.UpdateDNSRecordParams) (cf.DNSRecord, error)
	// DeleteDNSRecord 删除 DNS 记录
	DeleteDNSRecord(apiToken, zoneID, recordID string) error
}

type dnsClient struct{}

func NewDNSClient() DNSClient {
	return &dnsClient{}
}

func (c *dnsClient) ListDNSRecords(apiToken, zoneID string, params cf.ListDNSRecordsParams) ([]cf.DNSRecord, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.ZoneIdentifier(zoneID)

	records, _, err := client.ListDNSRecords(ctx, rc, params)
	if err != nil {
		return nil, fmt.Errorf("查询 DNS 记录失败: %w", err)
	}

	return records, nil
}

func (c *dnsClient) CreateDNSRecord(apiToken, zoneID string, params cf.CreateDNSRecordParams) (cf.DNSRecord, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return cf.DNSRecord{}, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.ZoneIdentifier(zoneID)

	record, err := client.CreateDNSRecord(ctx, rc, params)
	if err != nil {
		return cf.DNSRecord{}, fmt.Errorf("创建 DNS 记录失败: %w", err)
	}

	return record, nil
}

func (c *dnsClient) UpdateDNSRecord(apiToken, zoneID string, params cf.UpdateDNSRecordParams) (cf.DNSRecord, error) {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return cf.DNSRecord{}, fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.ZoneIdentifier(zoneID)

	record, err := client.UpdateDNSRecord(ctx, rc, params)
	if err != nil {
		return cf.DNSRecord{}, fmt.Errorf("更新 DNS 记录失败: %w", err)
	}

	return record, nil
}

func (c *dnsClient) DeleteDNSRecord(apiToken, zoneID, recordID string) error {
	client, err := cf.NewWithAPIToken(apiToken)
	if err != nil {
		return fmt.Errorf("创建 Cloudflare 客户端失败: %w", err)
	}

	ctx := context.Background()
	rc := cf.ZoneIdentifier(zoneID)

	if err := client.DeleteDNSRecord(ctx, rc, recordID); err != nil {
		return fmt.Errorf("删除 DNS 记录失败: %w", err)
	}

	return nil
}
