package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"auth-center/config"
	"auth-center/models"
)

// AppService 应用服务
type AppService struct{}

// generateAppID 生成应用ID
func (s *AppService) generateAppID() string {
	// 使用时间戳 + 随机数生成唯一ID
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("app_%d_%s", timestamp, randomStr)
}

// CreateAppRequest 创建应用请求
type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateAppResponse 创建应用响应
type CreateAppResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	AppID       string `json:"app_id"`
	AppSecret   string `json:"app_secret"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// UpdateAppRequest 更新应用请求
type UpdateAppRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

// AppListResponse 应用列表响应
type AppListResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	AppID       string `json:"app_id"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateApp 创建应用
func (s *AppService) CreateApp(req *CreateAppRequest) (*CreateAppResponse, error) {
	// 生成应用ID
	appID := s.generateAppID()
	
	// 检查应用ID是否已存在（理论上不应该重复，但为了安全起见）
	var existingApp models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&existingApp).Error; err == nil {
		// 如果重复，重新生成
		appID = s.generateAppID()
	}

	// 生成应用密钥
	appSecret, err := GenerateAppSecret()
	if err != nil {
		return nil, err
	}

	// 创建应用
	app := &models.Application{
		Name:        req.Name,
		AppID:       appID,
		AppSecret:   appSecret,
		Description: req.Description,
		Status:      1,
	}

	if err := config.DB.Create(app).Error; err != nil {
		// 如果因为唯一索引冲突（软删除未排除），返回更友好的错误
		if strings.Contains(err.Error(), "Duplicate entry") {
			return nil, fmt.Errorf("应用名称已存在: %s", req.Name)
		}
		return nil, fmt.Errorf("创建应用失败: %v", err)
	}

	return &CreateAppResponse{
		ID:          app.ID,
		Name:        app.Name,
		AppID:       app.AppID,
		AppSecret:   app.AppSecret,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetApp 获取应用信息
func (s *AppService) GetApp(appID string) (*AppListResponse, error) {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, fmt.Errorf("应用不存在")
	}

	return &AppListResponse{
		ID:          app.ID,
		Name:        app.Name,
		AppID:       app.AppID,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   app.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateApp 更新应用
func (s *AppService) UpdateApp(appID string, req *UpdateAppRequest) error {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return fmt.Errorf("应用不存在")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := config.DB.Model(&app).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新应用失败: %v", err)
	}

	return nil
}

// DeleteApp 删除应用
func (s *AppService) DeleteApp(appID string) error {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return fmt.Errorf("应用不存在")
	}

	if err := config.DB.Delete(&app).Error; err != nil {
		return fmt.Errorf("删除应用失败: %v", err)
	}

	return nil
}

// ListApps 获取应用列表
func (s *AppService) ListApps(page, pageSize int) ([]AppListResponse, int64, error) {
	var apps []models.Application
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总数
	if err := config.DB.Model(&models.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取应用列表
	if err := config.DB.Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	var responses []AppListResponse
	for _, app := range apps {
		responses = append(responses, AppListResponse{
			ID:          app.ID,
			Name:        app.Name,
			AppID:       app.AppID,
			Description: app.Description,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   app.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return responses, total, nil
}

// RegenerateAppSecret 重新生成应用密钥
func (s *AppService) RegenerateAppSecret(appID string) (string, error) {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return "", fmt.Errorf("应用不存在")
	}

	// 生成新的密钥
	newSecret, err := GenerateAppSecret()
	if err != nil {
		return "", err
	}

	// 更新密钥
	if err := config.DB.Model(&app).Update("app_secret", newSecret).Error; err != nil {
		return "", fmt.Errorf("更新密钥失败: %v", err)
	}

	return newSecret, nil
}

// GenerateAppSecret 生成应用密钥
func GenerateAppSecret() (string, error) {
	// 生成32字节的随机数据
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("生成随机数据失败: %v", err)
	}
	
	// 转换为十六进制字符串
	secret := hex.EncodeToString(bytes)
	
	// 添加前缀以便识别
	return fmt.Sprintf("app_%s", secret), nil
}