export interface TunnelConnectionVO {
  id: string
  colo_name: string
  client_id: string
  client_version: string
  opened_at: string
  origin_ip: string
}

export interface TunnelVO {
  id: string
  name: string
  status: string
  created_at?: string
  connections?: TunnelConnectionVO[]
}

export interface IngressRule {
  hostname: string
  service: string
  path?: string
  originRequest?: Record<string, unknown>
}

export interface TunnelConfig {
  ingress: IngressRule[]
  'warp-routing'?: Record<string, unknown>
  originRequest?: Record<string, unknown>
}

export interface TunnelConfigVO {
  tunnel_id: string
  config: TunnelConfig
  version: number
}

export interface TunnelTokenVO {
  token: string
}

export interface CreateTunnelRequest {
  credential_id: number
  name: string
}

export interface UpdateTunnelConfigRequest {
  credential_id: number
  config: Record<string, unknown>
}

export interface DNSRecordVO {
  id: string
  type: string
  name: string
  content: string
  proxied?: boolean
  ttl: number
}

export interface CreateDNSRecordRequest {
  credential_id: number
  zone_id: string
  name: string
  content: string
  proxied?: boolean
  ttl: number
}

export interface UpdateDNSRecordRequest {
  credential_id: number
  zone_id: string
  name: string
  content: string
  proxied?: boolean
  ttl: number
}

export interface DeleteDNSRecordRequest {
  credential_id: number
  zone_id: string
}

export interface ListDNSRecordsRequest {
  credential_id: number
  zone_id: string
  name?: string
  type?: string
}

export interface ListTunnelsQuery {
  credential_id: number
}

export interface GetTunnelQuery {
  credential_id: number
}

export interface DeleteTunnelQuery {
  credential_id: number
}

export interface GetTunnelTokenQuery {
  credential_id: number
}

export interface ListTunnelConnectionsQuery {
  credential_id: number
}

export interface GetTunnelConfigQuery {
  credential_id: number
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T | null
}