package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AuthClient 认证客户端
type AuthClient struct {
	BaseURL    string
	AppID      string
	AppSecret  string
	HTTPClient *http.Client
}

// NewAuthClient 创建认证客户端
func NewAuthClient(baseURL, appID, appSecret string) *AuthClient {
	return &AuthClient{
		BaseURL:   baseURL,
		AppID:     appID,
		AppSecret: appSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
	User         UserInfo `json:"user"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// PermissionCheckResponse 权限检查响应
type PermissionCheckResponse struct {
	HasPermission bool `json:"has_permission"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Login 用户登录
func (c *AuthClient) Login(req *LoginRequest) (*LoginResponse, error) {
	reqBody := map[string]interface{}{
		"app_id":   c.AppID,
		"username": req.Username,
		"password": req.Password,
	}

	return c.postJSON("/api/v1/auth/login", reqBody, &LoginResponse{})
}

// Register 用户注册
func (c *AuthClient) Register(req *RegisterRequest) error {
	reqBody := map[string]interface{}{
		"app_id":   c.AppID,
		"username": req.Username,
		"email":    req.Email,
		"password": req.Password,
	}

	_, err := c.postJSON("/api/v1/auth/register", reqBody, &APIResponse{})
	return err
}

// RefreshToken 刷新令牌
func (c *AuthClient) RefreshToken(req *RefreshTokenRequest) (*LoginResponse, error) {
	return c.postJSON("/api/v1/auth/refresh", req, &LoginResponse{})
}

// Logout 用户登出
func (c *AuthClient) Logout(token string) error {
	reqBody := map[string]interface{}{
		"token": token,
	}

	_, err := c.postJSON("/api/v1/auth/logout", reqBody, &APIResponse{})
	return err
}

// GetUserInfo 获取用户信息
func (c *AuthClient) GetUserInfo(token string) (*UserInfo, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	return c.getJSON("/api/v1/auth/user", headers, &UserInfo{})
}

// CheckPermission 检查权限
func (c *AuthClient) CheckPermission(token, permission string) (bool, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := fmt.Sprintf("/api/v1/permissions/check?permission=%s", permission)
	response, err := c.getJSON(url, headers, &PermissionCheckResponse{})
	if err != nil {
		return false, err
	}

	return response.HasPermission, nil
}

// CheckAPIPermission 检查API权限
func (c *AuthClient) CheckAPIPermission(token, path, method string) (bool, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := fmt.Sprintf("/api/v1/permissions/check-api?path=%s&method=%s", path, method)
	response, err := c.getJSON(url, headers, &PermissionCheckResponse{})
	if err != nil {
		return false, err
	}

	return response.HasPermission, nil
}

// GetUserPermissions 获取用户权限列表
func (c *AuthClient) GetUserPermissions(token string) ([]string, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	var response struct {
		Permissions []string `json:"permissions"`
	}

	_, err := c.getJSON("/api/v1/permissions/user", headers, &response)
	if err != nil {
		return nil, err
	}

	return response.Permissions, nil
}

// GetUserRoles 获取用户角色列表
func (c *AuthClient) GetUserRoles(token string) ([]RoleInfo, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	var response struct {
		Roles []RoleInfo `json:"roles"`
	}

	_, err := c.getJSON("/api/v1/permissions/roles", headers, &response)
	if err != nil {
		return nil, err
	}

	return response.Roles, nil
}

// postJSON 发送POST JSON请求
func (c *AuthClient) postJSON(path string, reqBody interface{}, response interface{}) (interface{}, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-App-Id", c.AppID)
	req.Header.Set("X-App-Secret", c.AppSecret)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errorResp APIResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("API错误: %s", errorResp.Error)
		}
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return nil, err
		}
		return response, nil
	}

	return nil, nil
}

// getJSON 发送GET JSON请求
func (c *AuthClient) getJSON(path string, headers map[string]string, response interface{}) (interface{}, error) {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errorResp APIResponse
		if err := json.Unmarshal(body, &errorResp); err == nil {
			return nil, fmt.Errorf("API错误: %s", errorResp.Error)
		}
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	if response != nil {
		if err := json.Unmarshal(body, response); err != nil {
			return nil, err
		}
		return response, nil
	}

	return nil, nil
}
