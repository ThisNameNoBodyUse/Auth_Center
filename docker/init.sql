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

-- 用户表：应用内用户
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
  app_id VARCHAR(64) NOT NULL COMMENT '所属应用ID',
  username VARCHAR(128) NOT NULL COMMENT '用户名（应用内唯一）',
  email VARCHAR(255) NULL COMMENT '邮箱（应用内唯一，可空）',
  phone VARCHAR(32) NULL COMMENT '手机号（可用于短信登录）',
  password VARCHAR(255) NOT NULL COMMENT '密码哈希',
  info JSON NULL COMMENT '用户扩展资料(JSON)',
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
('默认应用', 'default-app', 'default-secret', '系统默认应用', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入默认 provider（默认账号密码登录）
INSERT INTO providers (app_id, login_method) VALUES
('default-app', 0)
ON DUPLICATE KEY UPDATE login_method=VALUES(login_method);

-- 插入默认角色
INSERT INTO roles (app_id, name, code, description, status) VALUES
('default-app', '管理员', 'admin', '系统管理员角色', 1),
('default-app', '普通用户', 'user', '普通用户角色', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入默认权限
INSERT INTO permissions (app_id, name, code, resource, action, description, status) VALUES
('default-app', '用户管理', 'user:manage', 'user', 'manage', '用户管理权限', 1),
('default-app', '用户查看', 'user:read', 'user', 'read', '用户查看权限', 1),
('default-app', '用户创建', 'user:create', 'user', 'create', '用户创建权限', 1),
('default-app', '用户更新', 'user:update', 'user', 'update', '用户更新权限', 1),
('default-app', '用户删除', 'user:delete', 'user', 'delete', '用户删除权限', 1)
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 插入默认API
INSERT INTO apis (app_id, path, method, description, permission_id, status) VALUES
('default-app', '/api/v1/users', 'GET', '获取用户列表', 2, 1),
('default-app', '/api/v1/users', 'POST', '创建用户', 3, 1),
('default-app', '/api/v1/users/:id', 'PUT', '更新用户', 4, 1),
('default-app', '/api/v1/users/:id', 'DELETE', '删除用户', 5, 1)
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
