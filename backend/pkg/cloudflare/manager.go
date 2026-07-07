package cloudflare

import (
	"embed"
	"fmt"
	"os"
	"runtime"
	"sync"
)

//go:embed bin/*
var embeddedBinaries embed.FS

// Manager cloudflared 二进制管理器
type Manager struct {
	tmpPath string
	mu      sync.Mutex
}

// NewManager 创建二进制管理器
func NewManager() *Manager {
	return &Manager{}
}

// GetPath 获取 cloudflared 可执行文件路径
func (m *Manager) GetPath() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.tmpPath != "" {
		if _, err := os.Stat(m.tmpPath); err == nil {
			return m.tmpPath, nil
		}
	}

	key := fmt.Sprintf("bin/cloudflared-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		key += ".exe"
	}

	data, err := embeddedBinaries.ReadFile(key)
	if err != nil {
		return "", fmt.Errorf("未找到平台 %s/%s 的 cloudflared 二进制: %w", runtime.GOOS, runtime.GOARCH, err)
	}

	tmpFile, err := os.CreateTemp("", "cloudflared-*")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	if err := tmpFile.Chmod(0755); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("设置执行权限失败: %w", err)
	}

	tmpFile.Close()
	m.tmpPath = tmpFile.Name()
	return m.tmpPath, nil
}

// Cleanup 清理临时文件
func (m *Manager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.tmpPath != "" {
		os.Remove(m.tmpPath)
		m.tmpPath = ""
	}
}
