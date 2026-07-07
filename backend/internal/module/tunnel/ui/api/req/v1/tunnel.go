package v1

import (
	"time"

	cf "github.com/cloudflare/cloudflare-go"
)

// ==================== 隧道请求 ====================

// CreateTunnelRequest 创建隧道请求
type CreateTunnelRequest struct {
	CredentialID int64  `json:"credential_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
}

// UpdateTunnelConfigRequest 更新隧道配置请求
type UpdateTunnelConfigRequest struct {
	CredentialID int64                  `json:"credential_id" binding:"required"`
	Config       cf.TunnelConfiguration `json:"config" binding:"required"`
}

// ==================== DNS 请求 ====================

// ListDNSRecordsRequest 查询 DNS 记录列表请求
type ListDNSRecordsRequest struct {
	CredentialID int64  `form:"credential_id" binding:"required"`
	ZoneID       string `form:"zone_id" binding:"required"`
	Name         string `form:"name,omitempty"`
	Type         string `form:"type,omitempty"`
}

// CreateDNSRecordRequest 创建 DNS 记录请求
type CreateDNSRecordRequest struct {
	CredentialID int64  `json:"credential_id" binding:"required"`
	ZoneID       string `json:"zone_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Proxied      *bool  `json:"proxied"`
	TTL          int    `json:"ttl"`
}

// UpdateDNSRecordRequest 更新 DNS 记录请求
type UpdateDNSRecordRequest struct {
	CredentialID int64  `json:"credential_id" binding:"required"`
	ZoneID       string `json:"zone_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Content      string `json:"content" binding:"required"`
	Proxied      *bool  `json:"proxied"`
	TTL          int    `json:"ttl"`
}

// DeleteDNSRecordRequest 删除 DNS 记录请求
type DeleteDNSRecordRequest struct {
	CredentialID int64  `json:"credential_id" binding:"required"`
	ZoneID       string `json:"zone_id" binding:"required"`
}

// ==================== 响应 VO ====================

// TunnelVO 隧道详情
type TunnelVO struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Status      string                `json:"status"`
	CreatedAt   *time.Time            `json:"created_at,omitempty"`
	Connections []TunnelConnectionVO  `json:"connections,omitempty"`
}

// TunnelConnectionVO 隧道连接信息
type TunnelConnectionVO struct {
	ID            string `json:"id"`
	ColoName      string `json:"colo_name"`
	ClientID      string `json:"client_id"`
	ClientVersion string `json:"client_version"`
	OpenedAt      string `json:"opened_at"`
	OriginIP      string `json:"origin_ip"`
}

// TunnelConfigVO 隧道配置
type TunnelConfigVO struct {
	TunnelID string                  `json:"tunnel_id"`
	Config   cf.TunnelConfiguration  `json:"config"`
	Version  int                     `json:"version"`
}

// TunnelTokenVO 隧道连接 Token
type TunnelTokenVO struct {
	Token string `json:"token"`
}

// DNSRecordVO DNS 记录
type DNSRecordVO struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied *bool  `json:"proxied"`
	TTL     int    `json:"ttl"`
}
