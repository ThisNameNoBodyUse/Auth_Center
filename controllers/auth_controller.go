package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"auth-center/service"
)

// AuthController 认证控制器
type AuthController struct{}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录请求"
// @Success 200 {object} service.LoginResponse "登录成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "认证失败"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.AuthService{}
	response, err := authService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册新账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册请求"
// @Success 201 {object} map[string]string "注册成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 409 {object} map[string]string "用户已存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.AuthService{}
	if err := authService.Register(&req); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "注册成功"})
}

// RefreshToken 刷新令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} service.LoginResponse "刷新成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "令牌无效"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req service.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.AuthService{}
	response, err := authService.RefreshToken(&req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，使令牌失效
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LogoutRequest true "登出请求"
// @Success 200 {object} map[string]string "登出成功"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "令牌无效"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	var req service.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authService := &service.AuthService{}
	if err := authService.Logout(req.Token); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.UserInfo "用户信息"
// @Failure 401 {object} map[string]string "未认证"
// @Failure 404 {object} map[string]string "用户不存在"
// @Failure 500 {object} map[string]string "服务器错误"
// @Router /auth/user [get]
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	appID, exists := ctx.Get("app_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "应用未认证"})
		return
	}

	authService := &service.AuthService{}
	userInfo, err := authService.GetUserInfo(userID.(uint), appID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}
