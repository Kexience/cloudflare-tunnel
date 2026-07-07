package v1

// DashboardStatsVO 仪表盘统计数据
type DashboardStatsVO struct {
	RunningCount  int   `json:"running_count"`
	TotalCount    int   `json:"total_count"`
	BytesIn       int64 `json:"bytes_in"`
	BytesOut      int64 `json:"bytes_out"`
	TotalRequests int64 `json:"total_requests"`
}
