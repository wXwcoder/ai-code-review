package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// FileStorage 文件存储
type FileStorage struct {
	basePath string
}

// ReviewItem 审查记录项
type ReviewItem struct {
	RepoPath  string    `json:"repoPath"`
	Revision  int64     `json:"revision"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	FileName  string    `json:"fileName"`
}

// NewFileStorage 创建新的文件存储
func NewFileStorage(basePath string) *FileStorage {
	return &FileStorage{
		basePath: basePath,
	}
}

// sanitizePath 清理路径，把 / 替换成 _，避免文件名问题
func sanitizePath(path string) string {
	// 如果路径为空，使用默认名
	if path == "" {
		return "default"
	}
	// 替换掉文件系统不支持的字符
	safe := strings.ReplaceAll(path, "/", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	return safe
}

// GetReviewFile 获取审查文件路径
func (s *FileStorage) GetReviewFile(repoPath string, revision int64, extension string) string {
	safePath := sanitizePath(repoPath)
	fileName := fmt.Sprintf("%s_r%d%s", safePath, revision, extension)
	return filepath.Join(s.basePath, fileName)
}

// WriteFile 写入文件
func (s *FileStorage) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// ReadFile 读取文件
func (s *FileStorage) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// DeleteFile 删除文件
func (s *FileStorage) DeleteFile(path string) error {
	return os.Remove(path)
}

// FileExists 检查文件是否存在
func (s *FileStorage) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ListReviews 列出所有审查记录
func (s *FileStorage) ListReviews() ([]ReviewItem, error) {
	var items []ReviewItem

	// 确保存储目录存在
	if err := os.MkdirAll(s.basePath, 0755); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// 只处理 .md 文件，跳过 .processing 文件
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		// 解析文件名
		// 文件名格式: repoPath_r12345.md 或 default_r12345.md
		name := strings.TrimSuffix(file.Name(), ".md")
		parts := strings.Split(name, "_r")
		
		if len(parts) != 2 {
			continue
		}

		revision, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			continue
		}

		// 重建 repoPath（把 _ 还原回 /）
		repoPath := parts[0]
		if repoPath == "default" {
			repoPath = ""
		}

		// 读取文件内容
		contentPath := filepath.Join(s.basePath, file.Name())
		content, err := os.ReadFile(contentPath)
		if err != nil {
			continue
		}

		// 获取文件修改时间
		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		items = append(items, ReviewItem{
			RepoPath:  repoPath,
			Revision:  revision,
			Content:   string(content),
			Timestamp: fileInfo.ModTime(),
			FileName:  file.Name(),
		})
	}

	// 按时间倒序排列
	sort.Slice(items, func(i, j int) bool {
		return items[i].Timestamp.After(items[j].Timestamp)
	})

	return items, nil
}

// GetReview 获取单个审查记录
func (s *FileStorage) GetReview(repoPath string, revision int64) (*ReviewItem, error) {
	// 构建文件路径
	filePath := s.GetReviewFile(repoPath, revision, ".md")
	
	// 检查文件是否存在
	if !s.FileExists(filePath) {
		return nil, os.ErrNotExist
	}

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	return &ReviewItem{
		RepoPath:  repoPath,
		Revision:  revision,
		Content:   string(content),
		Timestamp: fileInfo.ModTime(),
		FileName:  filepath.Base(filePath),
	}, nil
}
