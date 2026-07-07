package svc

import (
	v1 "cloudflared-tunnel/internal/module/tunnel/ui/api/req/v1"
	"cloudflared-tunnel/pkg/errno"

	cf "github.com/cloudflare/cloudflare-go"
)

// DNSsvc DNS 记录管理服务接口
type DNSsvc interface {
	// ListDNSRecords 查询 DNS 记录列表
	ListDNSRecords(userID, credentialID int64, zoneID, name, recordType string) ([]*v1.DNSRecordVO, error)
	// CreateDNSRecord 创建 DNS 记录
	CreateDNSRecord(userID, credentialID int64, zoneID, name, content string, proxied *bool, ttl int) (*v1.DNSRecordVO, error)
	// UpdateDNSRecord 更新 DNS 记录
	UpdateDNSRecord(userID, credentialID int64, zoneID, recordID, name, content string, proxied *bool, ttl int) (*v1.DNSRecordVO, error)
	// DeleteDNSRecord 删除 DNS 记录
	DeleteDNSRecord(userID, credentialID int64, zoneID, recordID string) error
}

func (s *svc) ListDNSRecords(userID, credentialID int64, zoneID, name, recordType string) ([]*v1.DNSRecordVO, error) {
	_, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	params := cf.ListDNSRecordsParams{}
	if name != "" {
		params.Name = name
	}
	if recordType != "" {
		params.Type = recordType
	}

	records, err := s.dnsClient.ListDNSRecords(token, zoneID, params)
	if err != nil {
		s.log.Error("查询 DNS 记录失败", "zoneID", zoneID, "error", err)
		return nil, errno.ErrDNSRecordNotFound
	}

	vos := make([]*v1.DNSRecordVO, len(records))
	for i, record := range records {
		vos[i] = s.toDNSRecordVO(record)
	}

	return vos, nil
}

func (s *svc) CreateDNSRecord(userID, credentialID int64, zoneID, name, content string, proxied *bool, ttl int) (*v1.DNSRecordVO, error) {
	_, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	params := cf.CreateDNSRecordParams{
		Type:    "CNAME",
		Name:    name,
		Content: content,
		TTL:     ttl,
		Proxied: proxied,
	}

	record, err := s.dnsClient.CreateDNSRecord(token, zoneID, params)
	if err != nil {
		s.log.Error("创建 DNS 记录失败", "name", name, "error", err)
		return nil, errno.ErrDNSRecordCreate
	}

	return s.toDNSRecordVO(record), nil
}

func (s *svc) UpdateDNSRecord(userID, credentialID int64, zoneID, recordID, name, content string, proxied *bool, ttl int) (*v1.DNSRecordVO, error) {
	_, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return nil, err
	}

	params := cf.UpdateDNSRecordParams{
		ID:      recordID,
		Type:    "CNAME",
		Name:    name,
		Content: content,
		TTL:     ttl,
		Proxied: proxied,
	}

	record, err := s.dnsClient.UpdateDNSRecord(token, zoneID, params)
	if err != nil {
		s.log.Error("更新 DNS 记录失败", "recordID", recordID, "error", err)
		return nil, errno.ErrDNSRecordUpdate
	}

	return s.toDNSRecordVO(record), nil
}

func (s *svc) DeleteDNSRecord(userID, credentialID int64, zoneID, recordID string) error {
	_, token, err := s.getCredentialAndToken(userID, credentialID)
	if err != nil {
		return err
	}

	if err := s.dnsClient.DeleteDNSRecord(token, zoneID, recordID); err != nil {
		s.log.Error("删除 DNS 记录失败", "recordID", recordID, "error", err)
		return errno.ErrDNSRecordDelete
	}

	return nil
}
