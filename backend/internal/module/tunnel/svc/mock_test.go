package svc_test

import (
	"cloudflared-tunnel/pkg/cloudflare"

	cf "github.com/cloudflare/cloudflare-go"
)

// mockTunnelClient 隧道客户端 mock
type mockTunnelClient struct {
	listTunnelsFn           func(apiToken, accountID string) ([]cf.Tunnel, error)
	getTunnelFn             func(apiToken, accountID, tunnelID string) (cf.Tunnel, error)
	createTunnelFn          func(apiToken, accountID, name string) (cf.Tunnel, error)
	deleteTunnelFn          func(apiToken, accountID, tunnelID string) error
	getTunnelTokenFn        func(apiToken, accountID, tunnelID string) (string, error)
	listTunnelConnectionsFn func(apiToken, accountID, tunnelID string) ([]cf.Connection, error)
	getTunnelConfigFn       func(apiToken, accountID, tunnelID string) (cf.TunnelConfigurationResult, error)
	updateTunnelConfigFn    func(apiToken, accountID, tunnelID string, config cf.TunnelConfiguration) (cf.TunnelConfigurationResult, error)
}

func (m *mockTunnelClient) ListTunnels(apiToken, accountID string) ([]cf.Tunnel, error) {
	return m.listTunnelsFn(apiToken, accountID)
}

func (m *mockTunnelClient) GetTunnel(apiToken, accountID, tunnelID string) (cf.Tunnel, error) {
	return m.getTunnelFn(apiToken, accountID, tunnelID)
}

func (m *mockTunnelClient) CreateTunnel(apiToken, accountID, name string) (cf.Tunnel, error) {
	return m.createTunnelFn(apiToken, accountID, name)
}

func (m *mockTunnelClient) DeleteTunnel(apiToken, accountID, tunnelID string) error {
	return m.deleteTunnelFn(apiToken, accountID, tunnelID)
}

func (m *mockTunnelClient) GetTunnelToken(apiToken, accountID, tunnelID string) (string, error) {
	return m.getTunnelTokenFn(apiToken, accountID, tunnelID)
}

func (m *mockTunnelClient) ListTunnelConnections(apiToken, accountID, tunnelID string) ([]cf.Connection, error) {
	return m.listTunnelConnectionsFn(apiToken, accountID, tunnelID)
}

func (m *mockTunnelClient) GetTunnelConfig(apiToken, accountID, tunnelID string) (cf.TunnelConfigurationResult, error) {
	return m.getTunnelConfigFn(apiToken, accountID, tunnelID)
}

func (m *mockTunnelClient) UpdateTunnelConfig(apiToken, accountID, tunnelID string, config cf.TunnelConfiguration) (cf.TunnelConfigurationResult, error) {
	return m.updateTunnelConfigFn(apiToken, accountID, tunnelID, config)
}

// 确保 mockTunnelClient 实现 TunnelClient 接口
var _ cloudflare.TunnelClient = (*mockTunnelClient)(nil)

// mockDNSClient DNS 客户端 mock
type mockDNSClient struct {
	listDNSRecordsFn  func(apiToken, zoneID string, params cf.ListDNSRecordsParams) ([]cf.DNSRecord, error)
	createDNSRecordFn func(apiToken, zoneID string, params cf.CreateDNSRecordParams) (cf.DNSRecord, error)
	updateDNSRecordFn func(apiToken, zoneID string, params cf.UpdateDNSRecordParams) (cf.DNSRecord, error)
	deleteDNSRecordFn func(apiToken, zoneID, recordID string) error
}

func (m *mockDNSClient) ListDNSRecords(apiToken, zoneID string, params cf.ListDNSRecordsParams) ([]cf.DNSRecord, error) {
	return m.listDNSRecordsFn(apiToken, zoneID, params)
}

func (m *mockDNSClient) CreateDNSRecord(apiToken, zoneID string, params cf.CreateDNSRecordParams) (cf.DNSRecord, error) {
	return m.createDNSRecordFn(apiToken, zoneID, params)
}

func (m *mockDNSClient) UpdateDNSRecord(apiToken, zoneID string, params cf.UpdateDNSRecordParams) (cf.DNSRecord, error) {
	return m.updateDNSRecordFn(apiToken, zoneID, params)
}

func (m *mockDNSClient) DeleteDNSRecord(apiToken, zoneID, recordID string) error {
	return m.deleteDNSRecordFn(apiToken, zoneID, recordID)
}

// 确保 mockDNSClient 实现 DNSClient 接口
var _ cloudflare.DNSClient = (*mockDNSClient)(nil)
