package models

import (
	"time"

	"gorm.io/gorm"
)

// Application 应用模型
type Application struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	AppID       string         `json:"app_id" gorm:"uniqueIndex;not null"`
	AppSecret   string         `json:"app_secret" gorm:"not null"`
	Description string         `json:"description"`
	Status      int            `json:"status" gorm:"default:1"` // 1:启用 0:禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// User 用户模型
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	AppID        string         `json:"app_id" gorm:"index;not null"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex"`
	Phone        string         `json:"phone" gorm:"index"`
	Password     string         `json:"-" gorm:"not null"`                   // 不返回给前端
	IsSuperAdmin bool           `json:"is_super_admin" gorm:"default:false"` // 是否为超级管理员
	Status       int            `json:"status" gorm:"default:1"`             // 1:启用 0:禁用
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AppID       string         `json:"app_id" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	AppID       string         `json:"app_id" gorm:"index;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Code        string         `json:"code" gorm:"not null"`
	Resource    string         `json:"resource"` // 资源类型：menu, button, api等
	Action      string         `json:"action"`   // 操作类型：read, write, delete等
	Description string         `json:"description"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// API API接口模型
type API struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	AppID        string         `json:"app_id" gorm:"index;not null"`
	Path         string         `json:"path" gorm:"not null"`
	Method       string         `json:"method" gorm:"not null"`
	Description  string         `json:"description"`
	PermissionID uint           `json:"permission_id" gorm:"index"`
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// UserRole 用户角色关联表
type UserRole struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `json:"user_id" gorm:"index"`
	RoleID uint   `json:"role_id" gorm:"index"`
	AppID  string `json:"app_id" gorm:"index"`
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	RoleID       uint   `json:"role_id" gorm:"index"`
	PermissionID uint   `json:"permission_id" gorm:"index"`
	AppID        string `json:"app_id" gorm:"index"`
}

// Token 令牌模型（用于令牌管理）
type Token struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AppID     string    `json:"app_id" gorm:"index;not null"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Type      string    `json:"type"` // access, refresh
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Provider 登录方式提供方
type Provider struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AppID       string    `json:"app_id" gorm:"uniqueIndex;not null"`
	LoginMethod int       `json:"login_method" gorm:"not null;default:0"` // 0:账号密码 1:短信验证码
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SystemAdmin 系统管理员模型
type SystemAdmin struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;not null"`
	Email        string         `json:"email" gorm:"uniqueIndex"`
	Phone        string         `json:"phone" gorm:"index"`
	Password     string         `json:"-" gorm:"not null"`                   // 不返回给前端
	AdminType    string         `json:"admin_type" gorm:"not null"`          // system: 系统级管理员, app: 应用级管理员
	AppID        string         `json:"app_id" gorm:"index"`                 // 应用级管理员关联的应用ID
	IsActive     bool           `json:"is_active" gorm:"default:true"`       // 是否激活
	LastLoginAt  *time.Time     `json:"last_login_at"`                       // 最后登录时间
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 方法用于指定表名
func (Application) TableName() string {
	return "applications"
}

func (User) TableName() string {
	return "users"
}

func (Role) TableName() string {
	return "roles"
}

func (Permission) TableName() string {
	return "permissions"
}

func (API) TableName() string {
	return "apis"
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func (Token) TableName() string {
	return "tokens"
}

func (Provider) TableName() string {
	return "providers"
}

func (SystemAdmin) TableName() string {
	return "system_admins"
}
