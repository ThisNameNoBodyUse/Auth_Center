package main

import (
	"fmt"
	"log"

	"auth-sdk"
)

func main() {
	// 创建认证客户端
	client := auth.NewAuthClient("http://localhost:8080", "your-app-id", "your-app-secret")

	// 用户注册
	registerReq := &auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	if err := client.Register(registerReq); err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Println("注册成功")
	}

	// 用户登录
	loginReq := &auth.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResp, err := client.Login(loginReq)
	if err != nil {
		log.Printf("登录失败: %v", err)
		return
	}

	fmt.Printf("登录成功，访问令牌: %s\n", loginResp.AccessToken)

	// 获取用户信息
	userInfo, err := client.GetUserInfo(loginResp.AccessToken)
	if err != nil {
		log.Printf("获取用户信息失败: %v", err)
	} else {
		fmt.Printf("用户信息: %+v\n", userInfo)
	}

	// 检查权限
	hasPermission, err := client.CheckPermission(loginResp.AccessToken, "user:read")
	if err != nil {
		log.Printf("检查权限失败: %v", err)
	} else {
		fmt.Printf("是否有user:read权限: %t\n", hasPermission)
	}

	// 检查API权限
	hasAPIPermission, err := client.CheckAPIPermission(loginResp.AccessToken, "/api/v1/users", "GET")
	if err != nil {
		log.Printf("检查API权限失败: %v", err)
	} else {
		fmt.Printf("是否有API权限: %t\n", hasAPIPermission)
	}

	// 获取用户权限列表
	permissions, err := client.GetUserPermissions(loginResp.AccessToken)
	if err != nil {
		log.Printf("获取权限列表失败: %v", err)
	} else {
		fmt.Printf("用户权限: %v\n", permissions)
	}

	// 获取用户角色列表
	roles, err := client.GetUserRoles(loginResp.AccessToken)
	if err != nil {
		log.Printf("获取角色列表失败: %v", err)
	} else {
		fmt.Printf("用户角色: %v\n", roles)
	}

	// 用户登出
	if err := client.Logout(loginResp.AccessToken); err != nil {
		log.Printf("登出失败: %v", err)
	} else {
		fmt.Println("登出成功")
	}
}
