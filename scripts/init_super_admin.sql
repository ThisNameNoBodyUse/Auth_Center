-- 认证授权中心超级管理员初始化脚本
-- 注意：请在生产环境中修改默认密码

-- 1. 创建系统管理应用
INSERT IGNORE INTO applications (name, app_id, app_secret, description, status, created_at, updated_at) 
VALUES (
    '系统管理应用', 
    'system-admin', 
    'system-admin-secret-key-change-in-production', 
    '认证授权中心系统管理应用，用于管理所有应用', 
    1, 
    NOW(), 
    NOW()
);

-- 2. 为系统管理应用创建登录方式（密码登录）
INSERT IGNORE INTO providers (app_id, login_method, created_at, updated_at)
VALUES ('system-admin', 0, NOW(), NOW());

-- 3. 创建超级管理员用户
-- 密码: admin123 (请在生产环境中修改)
-- 密码哈希值: $argon2id$v=19$m=65536,t=3,p=4$r29vZ2xl$8J+QsQ==
INSERT IGNORE INTO users (app_id, username, email, password, is_super_admin, status, created_at, updated_at) 
VALUES (
    'system-admin',
    'superadmin',
    'admin@auth-center.com',
    '$argon2id$v=19$m=65536,t=3,p=4$r29vZ2xl$8J+QsQ==', -- 密码: admin123
    true,
    1,
    NOW(),
    NOW()
);  

-- 4. 创建超级管理员角色
INSERT IGNORE INTO roles (app_id, name, code, description, status, created_at, updated_at)
VALUES (
    'system-admin',
    '超级管理员',
    'super_admin',
    '系统超级管理员，拥有所有权限',
    1,
    NOW(),
    NOW()
);

-- 5. 创建系统管理权限
INSERT IGNORE INTO permissions (app_id, name, code, resource, action, description, status, created_at, updated_at)
VALUES 
    ('system-admin', '应用管理', 'app_manage', 'api', 'GET', '应用的所有操作权限', 1, NOW(), NOW());

-- 6. 将超级管理员角色分配给超级管理员用户
INSERT IGNORE INTO user_roles (user_id, role_id, app_id)
SELECT u.id, r.id, 'system-admin'
FROM users u, roles r
WHERE u.username = 'superadmin' 
  AND u.app_id = 'system-admin'
  AND r.code = 'super_admin'
  AND r.app_id = 'system-admin';

-- 7. 将系统管理权限分配给超级管理员角色
INSERT IGNORE INTO role_permissions (role_id, permission_id, app_id)
SELECT r.id, p.id, 'system-admin'
FROM roles r, permissions p
WHERE r.code = 'super_admin' 
  AND r.app_id = 'system-admin'
  AND p.app_id = 'system-admin';

-- 8. 创建系统管理相关的API权限
INSERT IGNORE INTO apis (app_id, path, method, description, permission_id, created_at, updated_at)
SELECT 
    'system-admin',
    '/api/v1/apps',
    'POST',
    '创建应用',
    p.id,
    NOW(),
    NOW()
FROM permissions p
WHERE p.code = 'app_manage' AND p.app_id = 'system-admin';

INSERT IGNORE INTO apis (app_id, path, method, description, permission_id, created_at, updated_at)
SELECT 
    'system-admin',
    '/api/v1/apps/*',
    'PUT',
    '更新应用',
    p.id,
    NOW(),
    NOW()
FROM permissions p
WHERE p.code = 'app_manage' AND p.app_id = 'system-admin';

INSERT IGNORE INTO apis (app_id, path, method, description, permission_id, created_at, updated_at)
SELECT 
    'system-admin',
    '/api/v1/apps/*',
    'DELETE',
    '删除应用',
    p.id,
    NOW(),
    NOW()
FROM permissions p
WHERE p.code = 'app_manage' AND p.app_id = 'system-admin';

-- 显示创建结果
SELECT '超级管理员初始化完成' as message;
SELECT '用户名: superadmin' as username;
SELECT '密码: admin123' as password;
SELECT '应用ID: system-admin' as app_id;
SELECT '应用密钥: system-admin-secret-key-change-in-production' as app_secret;
SELECT '请在生产环境中修改默认密码和应用密钥！' as warning;
