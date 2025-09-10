package service

import (
	"errors"
	"time"

	"auth-center/config"
	"auth-center/models"
	"auth-center/utils"

	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct{}

// LoginRequest 登录请求
type LoginRequest struct {
	AppID    string `json:"app_id" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Code     string `json:"code"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int64    `json:"expires_in"`
	TokenType    string   `json:"token_type"`
	User         UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint       `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Info     string     `json:"info"`
	Roles    []RoleInfo `json:"roles"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	AppID    string `json:"app_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required,min=6"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 验证应用是否存在
	var app models.Application
	if err := config.DB.Where("app_id = ? AND status = 1", req.AppID).First(&app).Error; err != nil {
		return nil, errors.New("应用不存在或已禁用")
	}

	// 读取应用登录方式
	loginMethod, err := s.getLoginMethod(req.AppID)
	if err != nil {
		return nil, err
	}

	// 根据 provider 决定校验方式
	var user models.User
	switch loginMethod {
	case 0: // 账号密码登录
		if req.Username == "" || req.Password == "" {
			return nil, errors.New("用户名与密码必填")
		}
		if err := config.DB.Where("username = ? AND app_id = ? AND status = 1", req.Username, req.AppID).First(&user).Error; err != nil {
			return nil, errors.New("用户不存在或已禁用")
		}
		valid, verr := utils.VerifyPassword(req.Password, user.Password)
		if verr != nil || !valid {
			return nil, errors.New("用户名或密码错误")
		}
	case 1: // 手机验证码登录
		if req.Phone == "" || req.Code == "" {
			return nil, errors.New("手机号与验证码必填")
		}
		if err := config.DB.Where("phone = ? AND app_id = ? AND status = 1", req.Phone, req.AppID).First(&user).Error; err != nil {
			return nil, errors.New("用户不存在或已禁用")
		}
		if ok := s.verifyOTP(req.AppID, req.Phone, req.Code); !ok {
			return nil, errors.New("验证码错误或已过期")
		}
	default:
		return nil, errors.New("不支持的登录方式")
	}

	// 获取用户角色
	roles, err := s.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}

	// 生成令牌
	accessToken, err := utils.GenerateAccessToken(user.ID, req.AppID, roles)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, req.AppID)
	if err != nil {
		return nil, err
	}

	// 获取角色信息
	roleInfos, err := s.GetRoleInfos(roles)
	if err != nil {
		return nil, err
	}

	// 保存令牌到数据库
	if err := s.saveToken(user.ID, req.AppID, accessToken, "access"); err != nil {
		return nil, err
	}
	if err := s.saveToken(user.ID, req.AppID, refreshToken, "refresh"); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.GetConfig().JWT.TTL,
		TokenType:    "Bearer",
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Info:     user.Info,
			Roles:    roleInfos,
		},
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) error {
	// 验证应用是否存在
	var app models.Application
	if err := config.DB.Where("app_id = ? AND status = 1", req.AppID).First(&app).Error; err != nil {
		return errors.New("应用不存在或已禁用")
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ? AND app_id = ?", req.Username, req.AppID).First(&existingUser).Error; err == nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在（如果提供了邮箱）
	if req.Email != "" {
		if err := config.DB.Where("email = ? AND app_id = ?", req.Email, req.AppID).First(&existingUser).Error; err == nil {
			return errors.New("邮箱已存在")
		}
	}

	// 哈希密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	// 创建用户
	user := models.User{
		AppID:    req.AppID,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Status:   1,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error) {
	// 解析刷新令牌
	claims, err := utils.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("id = ? AND app_id = ? AND status = 1", claims.UserID, claims.AppID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在或已禁用")
	}

	// 获取用户角色
	roles, err := s.GetUserRoles(user.ID)
	if err != nil {
		return nil, err
	}

	// 生成新的访问令牌
	accessToken, err := utils.GenerateAccessToken(user.ID, claims.AppID, roles)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(user.ID, claims.AppID)
	if err != nil {
		return nil, err
	}

	// 保存新令牌到数据库
	if err := s.saveToken(user.ID, claims.AppID, accessToken, "access"); err != nil {
		return nil, err
	}
	if err := s.saveToken(user.ID, claims.AppID, refreshToken, "refresh"); err != nil {
		return nil, err
	}

	// 获取角色信息
	roleInfos, err := s.GetRoleInfos(roles)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    config.GetConfig().JWT.TTL,
		TokenType:    "Bearer",
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Info:     user.Info,
			Roles:    roleInfos,
		},
	}, nil
}

// Logout 用户登出
func (s *AuthService) Logout(token string) error {
	// 解析令牌
	claims, err := utils.ParseAccessToken(token)
	if err != nil {
		return errors.New("无效的令牌")
	}

	// 将令牌加入黑名单
	blacklistKey := utils.TokenBlacklistPrefix + claims.JTI
	if err := utils.Set(blacklistKey, "1", time.Duration(config.GetConfig().JWT.TTL)*time.Second); err != nil {
		return err
	}

	return nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID uint, appID string) (*UserInfo, error) {
	var user models.User
	if err := config.DB.Where("id = ? AND app_id = ? AND status = 1", userID, appID).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在或已禁用")
	}

	// 获取用户角色
	roles, err := s.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// 获取角色信息
	roleInfos, err := s.GetRoleInfos(roles)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Info:     user.Info,
		Roles:    roleInfos,
	}, nil
}

// getUserRoles 获取用户角色ID列表
func (s *AuthService) GetUserRoles(userID uint) ([]uint, error) {
	var userRoles []models.UserRole
	if err := config.DB.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	return roleIDs, nil
}

// getRoleInfos 获取角色信息
func (s *AuthService) GetRoleInfos(roleIDs []uint) ([]RoleInfo, error) {
	if len(roleIDs) == 0 {
		return []RoleInfo{}, nil
	}

	var roles []models.Role
	if err := config.DB.Where("id IN ? AND status = 1", roleIDs).Find(&roles).Error; err != nil {
		return nil, err
	}

	var roleInfos []RoleInfo
	for _, role := range roles {
		roleInfos = append(roleInfos, RoleInfo{
			ID:   role.ID,
			Name: role.Name,
			Code: role.Code,
		})
	}

	return roleInfos, nil
}

// saveToken 保存令牌到数据库
func (s *AuthService) saveToken(userID uint, appID, token, tokenType string) error {
	var expiresAt time.Time
	if tokenType == "access" {
		expiresAt = time.Now().Add(time.Duration(config.GetConfig().JWT.TTL) * time.Second)
	} else {
		expiresAt = time.Now().Add(time.Duration(config.GetConfig().JWT.RefreshTTL) * time.Second)
	}

	tokenRecord := models.Token{
		AppID:     appID,
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}

	return config.DB.Create(&tokenRecord).Error
}

// getLoginMethod 获取应用登录方式（0:密码 1:短信验证码）
func (s *AuthService) getLoginMethod(appID string) (int, error) {
	var p models.Provider
	if err := config.DB.Where("app_id = ?", appID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return p.LoginMethod, nil
}

// verifyOTP 校验短信验证码（示例占位：从 Redis 读取）
func (s *AuthService) verifyOTP(appID, phone, code string) bool {
	key := "otp:" + appID + ":" + phone
	val, err := utils.Get(key)
	if err != nil {
		return false
	}
	return val == code
}
