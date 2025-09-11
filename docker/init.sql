-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS auth_center CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 使用数据库
USE auth_center;

-- 应用表：多租户根实体
CREATE TABLE IF NOT EXISTS applications (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  name VARCHAR(128) NOT NULL COMMENT '应用名称',
  app_id VARCHAR(64) NOT NULL COMMENT '应用唯一ID，供集成端使用',
  app_secret VARCHAR(128) NOT NULL COMMENT '应用密钥',
  description VARCHAR(255) NULL COMMENT '描述',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY uk_app_app_id (app_id),
  UNIQUE KEY uk_app_name (name),
  KEY idx_app_deleted_at (deleted_at)
) COMMENT='应用';

-- 用户表：应用内用户（去除info）
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL COMMENT '所属应用ID',
  username VARCHAR(128) NOT NULL COMMENT '用户名（应用内唯一）',
  email VARCHAR(255) NULL COMMENT '邮箱（应用内唯一，可空）',
  phone VARCHAR(32) NULL COMMENT '手机号（可用于短信登录）',
  password VARCHAR(255) NOT NULL COMMENT '密码哈希',
  is_super_admin TINYINT NOT NULL DEFAULT 0 COMMENT '是否为超级管理员 1是 0否',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY uk_user_app_username (app_id, username),
  UNIQUE KEY uk_user_app_email (app_id, email),
  KEY idx_user_phone (phone),
  KEY idx_user_app_id (app_id),
  KEY idx_user_deleted_at (deleted_at)
) COMMENT='用户';

-- 系统管理员表
CREATE TABLE IF NOT EXISTS system_admins (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  username VARCHAR(128) NOT NULL COMMENT '用户名',
  email VARCHAR(255) NULL COMMENT '邮箱',
  phone VARCHAR(32) NULL COMMENT '手机号',
  password VARCHAR(255) NOT NULL COMMENT '密码哈希',
  admin_type VARCHAR(20) NOT NULL COMMENT '管理员类型 system: 系统级管理员, app: 应用级管理员',
  app_id VARCHAR(64) NULL COMMENT '应用级管理员关联的应用ID',
  is_active TINYINT NOT NULL DEFAULT 1 COMMENT '是否激活 1是 0否',
  last_login_at DATETIME NULL COMMENT '最后登录时间',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY uk_sys_admin_username (username),
  UNIQUE KEY uk_sys_admin_email (email),
  KEY idx_sys_admin_phone (phone),
  KEY idx_sys_admin_type (admin_type),
  KEY idx_sys_admin_app_id (app_id),
  KEY idx_sys_admin_deleted_at (deleted_at)
) COMMENT='系统管理员';

-- 登录方式提供表（每应用一条）
CREATE TABLE IF NOT EXISTS providers (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL UNIQUE COMMENT '所属应用ID',
  login_method TINYINT NOT NULL DEFAULT 0 COMMENT '登录方式 0:账号密码 1:短信验证码',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT='应用登录提供方式';

-- 角色表：应用内角色
CREATE TABLE IF NOT EXISTS roles (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL COMMENT '所属应用ID',
  name VARCHAR(128) NOT NULL COMMENT '角色名称',
  code VARCHAR(128) NOT NULL COMMENT '角色编码（应用内唯一）',
  description VARCHAR(255) NULL COMMENT '描述',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY uk_role_app_code (app_id, code),
  KEY idx_role_app_id (app_id),
  KEY idx_role_deleted_at (deleted_at)
) COMMENT='角色';

-- 权限表：应用内权限点
CREATE TABLE IF NOT EXISTS permissions (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL COMMENT '所属应用ID',
  name VARCHAR(128) NOT NULL COMMENT '权限名称',
  code VARCHAR(128) NOT NULL COMMENT '权限编码（应用内唯一）',
  resource VARCHAR(64) NULL COMMENT '资源类型 menu/button/api等',
  action VARCHAR(64) NULL COMMENT '操作类型 read/write/delete等',
  description VARCHAR(255) NULL COMMENT '描述',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY uk_perm_app_code (app_id, code),
  KEY idx_perm_app_id (app_id),
  KEY idx_perm_deleted_at (deleted_at)
) COMMENT='权限';

-- API 表：受控接口资源
CREATE TABLE IF NOT EXISTS apis (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL COMMENT '所属应用ID',
  path VARCHAR(255) NOT NULL COMMENT 'API 路径（建议规范化）',
  method VARCHAR(16) NOT NULL COMMENT 'HTTP 方法',
  description VARCHAR(255) NULL COMMENT '描述',
  permission_id BIGINT UNSIGNED NULL COMMENT '绑定的权限ID',
  status TINYINT NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  deleted_at DATETIME NULL DEFAULT NULL COMMENT '软删除时间',
  UNIQUE KEY uk_api_app_path_method (app_id, path, method),
  KEY idx_api_app_id (app_id),
  KEY idx_api_permission_id (permission_id),
  KEY idx_api_deleted_at (deleted_at),
  CONSTRAINT fk_api_permission FOREIGN KEY (permission_id) REFERENCES permissions(id)
    ON DELETE SET NULL ON UPDATE CASCADE
) COMMENT='API 资源';

-- 用户-角色关联（多对多）
CREATE TABLE IF NOT EXISTS user_roles (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  app_id VARCHAR(64) NOT NULL COMMENT '应用ID（冗余用于隔离与约束）',
  UNIQUE KEY uk_user_role (app_id, user_id, role_id),
  KEY idx_ur_user (user_id),
  KEY idx_ur_role (role_id),
  CONSTRAINT fk_ur_user FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_ur_role FOREIGN KEY (role_id) REFERENCES roles(id)
    ON DELETE CASCADE ON UPDATE CASCADE
) COMMENT='用户-角色关联';

-- 角色-权限关联（多对多）
CREATE TABLE IF NOT EXISTS role_permissions (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
  permission_id BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
  app_id VARCHAR(64) NOT NULL COMMENT '应用ID（冗余用于隔离与约束）',
  UNIQUE KEY uk_role_perm (app_id, role_id, permission_id),
  KEY idx_rp_role (role_id),
  KEY idx_rp_perm (permission_id),
  CONSTRAINT fk_rp_role FOREIGN KEY (role_id) REFERENCES roles(id)
    ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_rp_perm FOREIGN KEY (permission_id) REFERENCES permissions(id)
    ON DELETE CASCADE ON UPDATE CASCADE
) COMMENT='角色-权限关联';

-- Token 表：记录签发令牌（审计/可选校验/黑名单辅助）
CREATE TABLE IF NOT EXISTS tokens (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL COMMENT '应用ID',
  user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  token TEXT NOT NULL COMMENT 'JWT 字符串（如需，可只存 jti）',
  jti VARCHAR(64) NULL COMMENT 'JWT ID，便于与黑名单联动',
  type VARCHAR(16) NOT NULL COMMENT '令牌类型 access/refresh',
  expires_at DATETIME NOT NULL COMMENT '过期时间',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  KEY idx_token_user (user_id),
  KEY idx_token_app (app_id),
  KEY idx_token_expires (expires_at),
  KEY idx_token_jti (jti),
  CONSTRAINT fk_token_user FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE ON UPDATE CASCADE
) COMMENT='令牌记录';



-- 插入默认应用
INSERT INTO applications (name, app_id, app_secret, description, status) VALUES
('默认应用', 'default-app', 'default-secret', '系统默认应用', 1),
('系统管理应用', 'system-admin', 'system-admin-secret-key-change-in-production', '认证授权中心系统管理应用，用于管理所有应用', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入默认 provider（默认账号密码登录）
INSERT INTO providers (app_id, login_method) VALUES
('default-app', 0),
('system-admin', 0)
ON DUPLICATE KEY UPDATE login_method=VALUES(login_method);

-- 插入默认角色
INSERT INTO roles (app_id, name, code, description, status) VALUES
('default-app', '管理员', 'admin', '系统管理员角色', 1),
('default-app', '普通用户', 'user', '普通用户角色', 1),
('system-admin', '超级管理员', 'super_admin', '系统超级管理员，拥有所有权限', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入默认权限
INSERT INTO permissions (app_id, name, code, resource, action, description, status) VALUES
('default-app', '用户管理', 'user:manage', 'user', 'manage', '用户管理权限', 1),
('default-app', '用户查看', 'user:read', 'user', 'read', '用户查看权限', 1),
('default-app', '用户创建', 'user:create', 'user', 'create', '用户创建权限', 1),
('default-app', '用户更新', 'user:update', 'user', 'update', '用户更新权限', 1),
('default-app', '用户删除', 'user:delete', 'user', 'delete', '用户删除权限', 1),
('system-admin', '应用管理', 'app_manage', 'application', 'all', '应用的所有操作权限', 1),
('system-admin', '用户管理', 'user_manage', 'user', 'all', '用户的所有操作权限', 1),
('system-admin', '角色管理', 'role_manage', 'role', 'all', '角色的所有操作权限', 1),
('system-admin', '权限管理', 'permission_manage', 'permission', 'all', '权限的所有操作权限', 1),
('system-admin', 'API管理', 'api_manage', 'api', 'all', 'API的所有操作权限', 1),
('system-admin', '系统监控', 'system_monitor', 'system', 'read', '系统监控权限', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入默认API
INSERT INTO apis (app_id, path, method, description, permission_id, status) VALUES
('default-app', '/api/v1/users', 'GET', '获取用户列表', 2, 1),
('default-app', '/api/v1/users', 'POST', '创建用户', 3, 1),
('default-app', '/api/v1/users/:id', 'PUT', '更新用户', 4, 1),
('default-app', '/api/v1/users/:id', 'DELETE', '删除用户', 5, 1),
('system-admin', '/api/v1/apps', 'POST', '创建应用', 6, 1),
('system-admin', '/api/v1/apps/*', 'PUT', '更新应用', 6, 1),
('system-admin', '/api/v1/apps/*', 'DELETE', '删除应用', 6, 1)
ON DUPLICATE KEY UPDATE description=VALUES(description);

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id, app_id) VALUES
(1, 1, 'default-app'),
(1, 2, 'default-app'),
(1, 3, 'default-app'),
(1, 4, 'default-app'),
(1, 5, 'default-app')
ON DUPLICATE KEY UPDATE permission_id=VALUES(permission_id);

-- 为普通用户角色分配查看权限
INSERT INTO role_permissions (role_id, permission_id, app_id) VALUES
(2, 2, 'default-app')
ON DUPLICATE KEY UPDATE permission_id=VALUES(permission_id);

-- 为超级管理员角色分配所有系统管理权限
INSERT INTO role_permissions (role_id, permission_id, app_id) VALUES
(3, 6, 'system-admin'),
(3, 7, 'system-admin'),
(3, 8, 'system-admin'),
(3, 9, 'system-admin'),
(3, 10, 'system-admin'),
(3, 11, 'system-admin')
ON DUPLICATE KEY UPDATE permission_id=VALUES(permission_id);

-- 插入超级管理员用户
INSERT INTO users (app_id, username, email, phone, password, is_super_admin, status) VALUES
('system-admin', 'superadmin', 'admin@auth-center.com', '', '$2b$10$UtbwZjygOigggJA.7So9v.cu0S1B.ibbBUNxdtA8GmwFVi86cZSye', 1, 1)
ON DUPLICATE KEY UPDATE username=VALUES(username);

-- 插入系统管理员
INSERT INTO system_admins (username, email, phone, password, admin_type, app_id, is_active) VALUES
('superadmin', 'admin@auth-center.com', '', '$2b$10$UtbwZjygOigggJA.7So9v.cu0S1B.ibbBUNxdtA8GmwFVi86cZSye', 'system', NULL, 1),
('appadmin', 'appadmin@auth-center.com', '13800138000', '$2b$10$UtbwZjygOigggJA.7So9v.cu0S1B.ibbBUNxdtA8GmwFVi86cZSye', 'app', 'default-app', 1)
ON DUPLICATE KEY UPDATE username=VALUES(username);

-- 为超级管理员用户分配超级管理员角色
INSERT INTO user_roles (user_id, role_id, app_id) VALUES
(1, 3, 'system-admin')
ON DUPLICATE KEY UPDATE role_id=VALUES(role_id);
