import api from '../api'
import type {
  TunnelVO,
  TunnelConfigVO,
  TunnelTokenVO,
  CreateTunnelRequest,
  UpdateTunnelConfigRequest,
  DNSRecordVO,
  CreateDNSRecordRequest,
  UpdateDNSRecordRequest,
  DeleteDNSRecordRequest,
  ListDNSRecordsRequest,
  ListTunnelsQuery,
  GetTunnelQuery,
  DeleteTunnelQuery,
  GetTunnelTokenQuery,
  ListTunnelConnectionsQuery,
  GetTunnelConfigQuery,
  ApiResponse
} from './types'

export async function listTunnels(query: ListTunnelsQuery): Promise<ApiResponse<TunnelVO[]>> {
  const response = await api.get<ApiResponse<TunnelVO[]>>('/v1/tunnels', { params: query })
  return response.data
}

export async function createTunnel(data: CreateTunnelRequest): Promise<ApiResponse<TunnelVO>> {
  const response = await api.post<ApiResponse<TunnelVO>>('/v1/tunnels', data)
  return response.data
}

export async function getTunnel(id: string, query: GetTunnelQuery): Promise<ApiResponse<TunnelVO>> {
  const response = await api.get<ApiResponse<TunnelVO>>(`/v1/tunnels/${id}`, { params: query })
  return response.data
}

export async function deleteTunnel(id: string, query: DeleteTunnelQuery): Promise<ApiResponse<null>> {
  const response = await api.delete<ApiResponse<null>>(`/v1/tunnels/${id}`, { params: query })
  return response.data
}

export async function getTunnelToken(id: string, query: GetTunnelTokenQuery): Promise<ApiResponse<TunnelTokenVO>> {
  const response = await api.get<ApiResponse<TunnelTokenVO>>(`/v1/tunnels/${id}/token`, { params: query })
  return response.data
}

export async function listTunnelConnections(id: string, query: ListTunnelConnectionsQuery): Promise<ApiResponse<TunnelVO>> {
  const response = await api.get<ApiResponse<TunnelVO>>(`/v1/tunnels/${id}/connections`, { params: query })
  return response.data
}

export async function getTunnelConfig(id: string, query: GetTunnelConfigQuery): Promise<ApiResponse<TunnelConfigVO>> {
  const response = await api.get<ApiResponse<TunnelConfigVO>>(`/v1/tunnels/${id}/config`, { params: query })
  return response.data
}

export async function updateTunnelConfig(id: string, data: UpdateTunnelConfigRequest): Promise<ApiResponse<TunnelConfigVO>> {
  const response = await api.put<ApiResponse<TunnelConfigVO>>(`/v1/tunnels/${id}/config`, data)
  return response.data
}

export async function listDNSRecords(query: ListDNSRecordsRequest): Promise<ApiResponse<DNSRecordVO[]>> {
  const response = await api.get<ApiResponse<DNSRecordVO[]>>('/v1/dns/records', { params: query })
  return response.data
}

export async function createDNSRecord(data: CreateDNSRecordRequest): Promise<ApiResponse<DNSRecordVO>> {
  const response = await api.post<ApiResponse<DNSRecordVO>>('/v1/dns/records', data)
  return response.data
}

export async function updateDNSRecord(id: string, data: UpdateDNSRecordRequest): Promise<ApiResponse<DNSRecordVO>> {
  const response = await api.put<ApiResponse<DNSRecordVO>>(`/v1/dns/records/${id}`, data)
  return response.data
}

export async function deleteDNSRecord(id: string, data: DeleteDNSRecordRequest): Promise<ApiResponse<null>> {
  const response = await api.delete<ApiResponse<null>>(`/v1/dns/records/${id}`, { data })
  return response.data
}