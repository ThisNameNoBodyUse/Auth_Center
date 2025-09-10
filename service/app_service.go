package service

import (
	"errors"
	"time"

	"auth-center/config"
	"auth-center/models"
	"github.com/google/uuid"
)

// AppService 应用服务
type AppService struct{}

// CreateAppRequest 创建应用请求
type CreateAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateAppResponse 创建应用响应
type CreateAppResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	AppID       string    `json:"app_id"`
	AppSecret   string    `json:"app_secret"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// UpdateAppRequest 更新应用请求
type UpdateAppRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      *int   `json:"status"`
}

// AppListResponse 应用列表响应
type AppListResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	AppID       string    `json:"app_id"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateApp 创建应用
func (s *AppService) CreateApp(req *CreateAppRequest) (*CreateAppResponse, error) {
	// 检查应用名称是否已存在
	var existingApp models.Application
	if err := config.DB.Where("name = ?", req.Name).First(&existingApp).Error; err == nil {
		return nil, errors.New("应用名称已存在")
	}

	// 生成应用ID和密钥
	appID := generateAppID()
	appSecret := generateAppSecret()

	// 创建应用
	app := models.Application{
		Name:        req.Name,
		AppID:       appID,
		AppSecret:   appSecret,
		Description: req.Description,
		Status:      1,
	}

	if err := config.DB.Create(&app).Error; err != nil {
		return nil, err
	}

	return &CreateAppResponse{
		ID:          app.ID,
		Name:        app.Name,
		AppID:       app.AppID,
		AppSecret:   app.AppSecret,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt,
	}, nil
}

// GetApp 获取应用信息
func (s *AppService) GetApp(appID string) (*AppListResponse, error) {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在")
	}

	return &AppListResponse{
		ID:          app.ID,
		Name:        app.Name,
		AppID:       app.AppID,
		Description: app.Description,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}, nil
}

// UpdateApp 更新应用
func (s *AppService) UpdateApp(appID string, req *UpdateAppRequest) error {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return errors.New("应用不存在")
	}

	// 更新字段
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

	if len(updates) > 0 {
		if err := config.DB.Model(&app).Updates(updates).Error; err != nil {
			return err
		}
	}

	return nil
}

// DeleteApp 删除应用
func (s *AppService) DeleteApp(appID string) error {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return errors.New("应用不存在")
	}

	// 软删除
	if err := config.DB.Delete(&app).Error; err != nil {
		return err
	}

	return nil
}

// ListApps 获取应用列表
func (s *AppService) ListApps(page, pageSize int) ([]AppListResponse, int64, error) {
	var apps []models.Application
	var total int64

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 查询总数
	if err := config.DB.Model(&models.Application{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 查询应用列表
	if err := config.DB.Offset(offset).Limit(pageSize).Find(&apps).Error; err != nil {
		return nil, 0, err
	}

	var appList []AppListResponse
	for _, app := range apps {
		appList = append(appList, AppListResponse{
			ID:          app.ID,
			Name:        app.Name,
			AppID:       app.AppID,
			Description: app.Description,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		})
	}

	return appList, total, nil
}

// RegenerateAppSecret 重新生成应用密钥
func (s *AppService) RegenerateAppSecret(appID string) (string, error) {
	var app models.Application
	if err := config.DB.Where("app_id = ?", appID).First(&app).Error; err != nil {
		return "", errors.New("应用不存在")
	}

	// 生成新的应用密钥
	newSecret := generateAppSecret()

	// 更新数据库
	if err := config.DB.Model(&app).Update("app_secret", newSecret).Error; err != nil {
		return "", err
	}

	return newSecret, nil
}

// generateAppID 生成应用ID
func generateAppID() string {
	return "app_" + uuid.New().String()[:8]
}

// generateAppSecret 生成应用密钥
func generateAppSecret() string {
	return "secret_" + uuid.New().String()
}
