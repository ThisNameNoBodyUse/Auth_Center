package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"auth-center/controllers"
	"auth-center/service"
	"auth-center/utils"
)

func TestAuthController(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	// 创建测试路由
	r := gin.New()
	authController := &controllers.AuthController{}

	// 注册路由
	r.POST("/login", authController.Login)
	r.POST("/register", authController.Register)

	t.Run("用户注册", func(t *testing.T) {
		reqBody := service.RegisterRequest{
			AppID:    "test-app",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("期望状态码 %d，实际得到 %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("用户登录", func(t *testing.T) {
		reqBody := service.LoginRequest{
			AppID:    "test-app",
			Username: "testuser",
			Password: "password123",
		}

		jsonData, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("期望状态码 %d，实际得到 %d", http.StatusOK, w.Code)
		}

		var response service.LoginResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("解析响应失败: %v", err)
		}

		if response.AccessToken == "" {
			t.Error("访问令牌为空")
		}

		if response.RefreshToken == "" {
			t.Error("刷新令牌为空")
		}
	})
}

func TestPasswordUtils(t *testing.T) {
	password := "testpassword123"
	
	// 测试密码哈希
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Errorf("密码哈希失败: %v", err)
	}

	if hashedPassword == "" {
		t.Error("哈希密码为空")
	}

	// 测试密码验证
	valid, err := utils.VerifyPassword(password, hashedPassword)
	if err != nil {
		t.Errorf("密码验证失败: %v", err)
	}

	if !valid {
		t.Error("密码验证失败")
	}

	// 测试错误密码
	valid, err = utils.VerifyPassword("wrongpassword", hashedPassword)
	if err != nil {
		t.Errorf("密码验证失败: %v", err)
	}

	if valid {
		t.Error("错误密码应该验证失败")
	}
}

func TestJWTUtils(t *testing.T) {
	userID := uint(1)
	appID := "test-app"
	roles := []uint{1, 2}

	// 测试生成访问令牌
	accessToken, err := utils.GenerateAccessToken(userID, appID, roles)
	if err != nil {
		t.Errorf("生成访问令牌失败: %v", err)
	}

	if accessToken == "" {
		t.Error("访问令牌为空")
	}

	// 测试解析访问令牌
	claims, err := utils.ParseAccessToken(accessToken)
	if err != nil {
		t.Errorf("解析访问令牌失败: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("期望用户ID %d，实际得到 %d", userID, claims.UserID)
	}

	if claims.AppID != appID {
		t.Errorf("期望应用ID %s，实际得到 %s", appID, claims.AppID)
	}

	if len(claims.Roles) != len(roles) {
		t.Errorf("期望角色数量 %d，实际得到 %d", len(roles), len(claims.Roles))
	}

	// 测试生成刷新令牌
	refreshToken, err := utils.GenerateRefreshToken(userID, appID)
	if err != nil {
		t.Errorf("生成刷新令牌失败: %v", err)
	}

	if refreshToken == "" {
		t.Error("刷新令牌为空")
	}

	// 测试解析刷新令牌
	refreshClaims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		t.Errorf("解析刷新令牌失败: %v", err)
	}

	if refreshClaims.UserID != userID {
		t.Errorf("期望用户ID %d，实际得到 %d", userID, refreshClaims.UserID)
	}

	if refreshClaims.AppID != appID {
		t.Errorf("期望应用ID %s，实际得到 %s", appID, refreshClaims.AppID)
	}
}
