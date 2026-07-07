package cloudflare

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TunnelMetrics 隧道流量指标
type TunnelMetrics struct {
	BytesIn       int64
	BytesOut      int64
	TotalRequests int64
}

// MetricsClient Prometheus 指标客户端
type MetricsClient struct {
	httpClient *http.Client
}

// NewMetricsClient 创建指标客户端
func NewMetricsClient() *MetricsClient {
	return &MetricsClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// FetchMetrics 从 cloudflared 的 metrics 端点获取指标
func (c *MetricsClient) FetchMetrics(port int) (*TunnelMetrics, error) {
	url := fmt.Sprintf("http://localhost:%d/metrics", port)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求 metrics 端点失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metrics 端点返回状态码: %d", resp.StatusCode)
	}

	return parsePrometheusMetrics(resp.Body)
}

// parsePrometheusMetrics 解析 Prometheus 文本格式
func parsePrometheusMetrics(r io.Reader) (*TunnelMetrics, error) {
	metrics := &TunnelMetrics{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// 解析 cloudflared_tunnel_bytes_transferred
		if strings.HasPrefix(line, "cloudflared_tunnel_bytes_transferred") {
			value, labels := parseMetricLine(line)
			if value < 0 {
				continue
			}

			if strings.Contains(labels, `direction="ingress"`) {
				metrics.BytesIn = int64(value)
			} else if strings.Contains(labels, `direction="egress"`) {
				metrics.BytesOut = int64(value)
			}
		}

		// 解析 cloudflared_tunnel_total_requests (如果存在)
		if strings.HasPrefix(line, "cloudflared_tunnel_total_requests") {
			value, _ := parseMetricLine(line)
			if value >= 0 {
				metrics.TotalRequests = int64(value)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取 metrics 失败: %w", err)
	}

	return metrics, nil
}

// parseMetricLine 解析单行 Prometheus 指标
// 格式: metric_name{label="value"} 123.456
func parseMetricLine(line string) (float64, string) {
	// 找到第一个 { 和最后一个 }
	start := strings.Index(line, "{")
	end := strings.LastIndex(line, "}")

	if start == -1 || end == -1 {
		// 没有标签的指标
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			val, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return -1, ""
			}
			return val, ""
		}
		return -1, ""
	}

	labels := line[start+1 : end]
	valuePart := line[end+1:]
	parts := strings.Fields(strings.TrimSpace(valuePart))
	if len(parts) >= 1 {
		val, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return -1, labels
		}
		return val, labels
	}

	return -1, labels
}
