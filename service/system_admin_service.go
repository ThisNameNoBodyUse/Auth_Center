package service

import (
	"errors"
	"time"

	"auth-center/config"
	"auth-center/models"
	"auth-center/utils"
	"github.com/golang-jwt/jwt/v5"
)

// SystemAdminService 系统管理员服务
type SystemAdminService struct{}

// SystemLoginRequest 系统登录请求
type SystemLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SystemLoginResponse 系统登录响应
type SystemLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Admin        struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		Email      string `json:"email"`
		AdminType  string `json:"admin_type"`
		AppID      string `json:"app_id,omitempty"`
		IsActive   bool   `json:"is_active"`
	} `json:"admin"`
}

// SystemRegisterRequest 系统管理员注册请求
type SystemRegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
	AdminType string `json:"admin_type" binding:"required,oneof=system app"`
	AppID     string `json:"app_id,omitempty"` // 应用级管理员需要指定应用ID
}

// SystemRegisterResponse 系统管理员注册响应
type SystemRegisterResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AdminType string `json:"admin_type"`
	AppID     string `json:"app_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// SystemLogin 系统管理员登录
func (s *SystemAdminService) SystemLogin(req *SystemLoginRequest) (*SystemLoginResponse, error) {
	// 查找系统管理员
	var admin models.SystemAdmin
	if err := config.DB.Where("username = ? AND is_active = ?", req.Username, true).First(&admin).Error; err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	valid, err := utils.VerifyPassword(req.Password, admin.Password)
	if err != nil || !valid {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	config.DB.Model(&admin).Update("last_login_at", now)

	// 生成JWT令牌
	accessToken, refreshToken, err := s.generateTokens(admin.ID, admin.Username, admin.AdminType, admin.AppID)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 保存令牌到数据库
	if err := s.saveTokens(admin.ID, accessToken, refreshToken); err != nil {
		return nil, errors.New("保存令牌失败")
	}

	response := &SystemLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1小时
	}

	response.Admin.ID = admin.ID
	response.Admin.Username = admin.Username
	response.Admin.Email = admin.Email
	response.Admin.AdminType = admin.AdminType
	response.Admin.AppID = admin.AppID
	response.Admin.IsActive = admin.IsActive

	return response, nil
}

// SystemRegister 系统管理员注册
func (s *SystemAdminService) SystemRegister(req *SystemRegisterRequest) (*SystemRegisterResponse, error) {
	// 验证应用级管理员必须指定应用ID
	if req.AdminType == "app" && req.AppID == "" {
		return nil, errors.New("应用级管理员必须指定应用ID")
	}

	// 验证应用是否存在（仅应用级管理员）
	if req.AdminType == "app" {
		var app models.Application
		if err := config.DB.Where("app_id = ? AND status = ?", req.AppID, 1).First(&app).Error; err != nil {
			return nil, errors.New("指定的应用不存在或已禁用")
		}
	}

	// 检查用户名是否已存在
	var existingAdmin models.SystemAdmin
	if err := config.DB.Where("username = ?", req.Username).First(&existingAdmin).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingAdmin).Error; err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建系统管理员
	admin := &models.SystemAdmin{
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		Password:  hashedPassword,
		AdminType: req.AdminType,
		AppID:     req.AppID,
		IsActive:  true,
	}

	if err := config.DB.Create(admin).Error; err != nil {
		return nil, errors.New("创建管理员失败")
	}

	response := &SystemRegisterResponse{
		ID:        admin.ID,
		Username:  admin.Username,
		Email:     admin.Email,
		AdminType: admin.AdminType,
		AppID:     admin.AppID,
		CreatedAt: admin.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// SystemRefreshToken 刷新系统管理员令牌
func (s *SystemAdminService) SystemRefreshToken(refreshToken string) (*SystemLoginResponse, error) {
	// 验证刷新令牌
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 查找管理员
	var admin models.SystemAdmin
	if err := config.DB.Where("id = ? AND is_active = ?", claims["user_id"], true).First(&admin).Error; err != nil {
		return nil, errors.New("管理员不存在或已禁用")
	}

	// 生成新的令牌
	accessToken, newRefreshToken, err := s.generateTokens(admin.ID, admin.Username, admin.AdminType, admin.AppID)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 保存新令牌
	if err := s.saveTokens(admin.ID, accessToken, newRefreshToken); err != nil {
		return nil, errors.New("保存令牌失败")
	}

	response := &SystemLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600,
	}

	response.Admin.ID = admin.ID
	response.Admin.Username = admin.Username
	response.Admin.Email = admin.Email
	response.Admin.AdminType = admin.AdminType
	response.Admin.AppID = admin.AppID
	response.Admin.IsActive = admin.IsActive

	return response, nil
}

// generateTokens 生成JWT令牌
func (s *SystemAdminService) generateTokens(adminID uint, username, adminType, appID string) (string, string, error) {
	// 生成访问令牌
	accessClaims := jwt.MapClaims{
		"user_id":    adminID,
		"username":   username,
		"admin_type": adminType,
		"app_id":     appID,
		"type":       "system_admin",
		"exp":        time.Now().Add(time.Hour).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        "auth-center",
		"sub":        "access-token",
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshClaims := jwt.MapClaims{
		"user_id":    adminID,
		"username":   username,
		"admin_type": adminType,
		"app_id":     appID,
		"type":       "system_admin",
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
		"iss":        "auth-center",
		"sub":        "refresh-token",
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// validateToken 验证令牌
func (s *SystemAdminService) validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的令牌")
}

// saveTokens 保存令牌到数据库
func (s *SystemAdminService) saveTokens(adminID uint, accessToken, refreshToken string) error {
	// 删除旧的令牌
	config.DB.Where("user_id = ? AND type IN (?)", adminID, []string{"access", "refresh"}).Delete(&models.Token{})

	// 保存访问令牌
	accessTokenRecord := &models.Token{
		AppID:     "system-admin",
		UserID:    adminID,
		Token:     accessToken,
		Type:      "access",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// 保存刷新令牌
	refreshTokenRecord := &models.Token{
		AppID:     "system-admin",
		UserID:    adminID,
		Token:     refreshToken,
		Type:      "refresh",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := config.DB.Create(accessTokenRecord).Error; err != nil {
		return err
	}

	if err := config.DB.Create(refreshTokenRecord).Error; err != nil {
		return err
	}

	return nil
}

// GetSystemAdminInfo 获取系统管理员信息
func (s *SystemAdminService) GetSystemAdminInfo(adminID uint) (*models.SystemAdmin, error) {
	var admin models.SystemAdmin
	if err := config.DB.Where("id = ? AND is_active = ?", adminID, true).First(&admin).Error; err != nil {
		return nil, errors.New("管理员不存在")
	}
	return &admin, nil
}
