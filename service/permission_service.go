package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"auth-center/config"
	"auth-center/models"
	"auth-center/utils"
)

// PermissionService 权限服务
type PermissionService struct{}

// CheckUserPermission 检查用户权限
func (s *PermissionService) CheckUserPermission(userID uint, appID, permission string) (bool, error) {
	// 先尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("%s%d:%s", utils.UserPermissionPrefix, userID, appID)
	permissions, err := utils.SMembers(cacheKey)
	if err != nil || len(permissions) == 0 {
		// 缓存未命中，从数据库查询
		permissions, err = s.GetUserPermissionsFromDB(userID, appID)
		if err != nil {
			return false, err
		}

		// 缓存到Redis
		if len(permissions) > 0 {
			utils.SAdd(cacheKey, permissions)
			utils.Expire(cacheKey, 12*time.Hour)
		}
	}

	// 检查权限
	for _, perm := range permissions {
		if perm == permission {
			return true, nil
		}
	}

	return false, nil
}

// CheckAPIPermission 检查API权限
func (s *PermissionService) CheckAPIPermission(userID uint, appID, apiPath, apiMethod string) (bool, error) {
	// 获取用户角色
	var userRoles []models.UserRole
	if err := config.DB.Where("user_id = ? AND app_id = ?", userID, appID).Find(&userRoles).Error; err != nil {
		return false, err
	}

	if len(userRoles) == 0 {
		return false, nil
	}

	// 检查每个角色的权限
	for _, userRole := range userRoles {
		// 获取角色权限
		rolePermissions, err := s.getRolePermissions(userRole.RoleID, appID)
		if err != nil {
			return false, err
		}

		// 检查API权限
		for _, permissionID := range rolePermissions {
			hasPermission, err := s.checkAPIPermissionByPermissionID(permissionID, appID, apiPath, apiMethod)
			if err != nil {
				return false, err
			}
			if hasPermission {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetUserPermissionsFromDB 从数据库获取用户权限
func (s *PermissionService) GetUserPermissionsFromDB(userID uint, appID string) ([]string, error) {
	var permissions []string

	// 查询用户角色
	var userRoles []models.UserRole
	if err := config.DB.Where("user_id = ? AND app_id = ?", userID, appID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	// 查询每个角色的权限
	for _, userRole := range userRoles {
		rolePermissions, err := s.getRolePermissions(userRole.RoleID, appID)
		if err != nil {
			return nil, err
		}

		// 获取权限代码
		for _, permissionID := range rolePermissions {
			var permission models.Permission
			if err := config.DB.Where("id = ? AND app_id = ? AND status = 1", permissionID, appID).First(&permission).Error; err != nil {
				continue
			}
			permissions = append(permissions, permission.Code)
		}
	}

	return permissions, nil
}

// getRolePermissions 获取角色权限
func (s *PermissionService) getRolePermissions(roleID uint, appID string) ([]uint, error) {
	// 先尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("%s%d:%s", utils.RolePermissionPrefix, roleID, appID)
	permissions, err := utils.SMembers(cacheKey)
	if err != nil || len(permissions) == 0 {
		// 缓存未命中，从数据库查询
		var rolePermissions []models.RolePermission
		if err := config.DB.Where("role_id = ? AND app_id = ?", roleID, appID).Find(&rolePermissions).Error; err != nil {
			return nil, err
		}

		var permissionIDs []uint
		for _, rp := range rolePermissions {
			permissionIDs = append(permissionIDs, rp.PermissionID)
		}

		// 缓存到Redis
		if len(permissionIDs) > 0 {
			var interfaceSlice []interface{}
			for _, id := range permissionIDs {
				interfaceSlice = append(interfaceSlice, id)
			}
			utils.SAdd(cacheKey, interfaceSlice...)
			utils.Expire(cacheKey, 12*time.Hour)
		}

		return permissionIDs, nil
	}

	// 从缓存结果转换
	var permissionIDs []uint
	for _, perm := range permissions {
		if id, err := strconv.ParseUint(perm, 10, 32); err == nil {
			permissionIDs = append(permissionIDs, uint(id))
		}
	}
	return permissionIDs, nil
}

// checkAPIPermissionByPermissionID 根据权限ID检查API权限
func (s *PermissionService) checkAPIPermissionByPermissionID(permissionID uint, appID, apiPath, apiMethod string) (bool, error) {
	// 先尝试从Redis缓存获取
	cacheKey := fmt.Sprintf("%s%d:%s", utils.APIPermissionPrefix, permissionID, appID)
	apis, err := utils.SMembers(cacheKey)
	if err != nil || len(apis) == 0 {
		// 缓存未命中，从数据库查询
		apis, err = s.getAPIsByPermissionID(permissionID, appID)
		if err != nil {
			return false, err
		}

		// 缓存到Redis
		if len(apis) > 0 {
			utils.SAdd(cacheKey, apis)
			utils.Expire(cacheKey, 24*time.Hour)
		}
	}

	// 检查API权限
	requestKey := apiPath + ":" + apiMethod
	for _, api := range apis {
		if api == requestKey {
			return true, nil
		}
	}

	return false, nil
}

// getAPIsByPermissionID 根据权限ID获取API列表
func (s *PermissionService) getAPIsByPermissionID(permissionID uint, appID string) ([]string, error) {
	var apis []models.API
	if err := config.DB.Where("permission_id = ? AND app_id = ? AND status = 1", permissionID, appID).Find(&apis).Error; err != nil {
		return nil, err
	}

	var apiList []string
	for _, api := range apis {
		apiList = append(apiList, api.Path+":"+api.Method)
	}

	return apiList, nil
}

// ValidateAppCredentials 验证应用凭据
func ValidateAppCredentials(appID, appSecret string) (*models.Application, error) {
	var app models.Application
	if err := config.DB.Where("app_id = ? AND app_secret = ?", appID, appSecret).First(&app).Error; err != nil {
		return nil, errors.New("无效的应用凭据")
	}
	return &app, nil
}

// CheckUserPermission 检查用户权限（全局函数）
func CheckUserPermission(userID uint, appID, permission string) (bool, error) {
	service := &PermissionService{}
	return service.CheckUserPermission(userID, appID, permission)
}

// CheckAPIPermission 检查API权限（全局函数）
func CheckAPIPermission(userID uint, appID, apiPath, apiMethod string) (bool, error) {
	service := &PermissionService{}
	return service.CheckAPIPermission(userID, appID, apiPath, apiMethod)
}
