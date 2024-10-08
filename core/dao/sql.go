package dao

import "strings"

var DropTableSql = map[string]string{
	"app":         "DROP TABLE IF EXISTS `app`;",
	"config":      "DROP TABLE IF EXISTS `config`;",
	"release":     "DROP TABLE IF EXISTS `release`;",
	"release_log": "DROP TABLE IF EXISTS `release_log`;",
	"user":        "DROP TABLE IF EXISTS `user`;",
}

var CreateTableSql = map[string]string{
	"app": `
		CREATE TABLE IF NOT EXISTS "app" (
			"app_id" bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
			"name" varchar(255) NOT NULL COMMENT '应用代码',
			"code" varchar(255) NOT NULL COMMENT '应用代号',
			"desc" varchar(255) DEFAULT NULL COMMENT '描述',
			"api_key" varchar(255) DEFAULT NULL COMMENT 'API密钥',
			"release_index" int(11) NOT NULL DEFAULT '0' COMMENT '发布INDEX',
			"create_time" datetime NOT NULL COMMENT '创建时间',
			"update_time" datetime NOT NULL COMMENT '更新时间',
			PRIMARY KEY ("app_id"),
			UNIQUE KEY "uniq_code" ("code"),
			UNIQUE KEY "uniq_api_key" ("api_key"),
			KEY "idx_name" ("name")
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='App信息表';
	`,
	"config": `
		CREATE TABLE IF NOT EXISTS "config" (
			"config_id" bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
			"app_id" bigint(20) NOT NULL COMMENT '应用ID',
			"key" varchar(255) NOT NULL COMMENT '配置Key',
			"value" longtext CHARACTER SET utf8mb4 NOT NULL COMMENT '配置Value',
			"desc" varchar(255) NOT NULL COMMENT '配置描述',
			"status" tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态（1-待发布、2-已发布）',
			"operate" tinyint(4) NOT NULL DEFAULT '1' COMMENT '操作标志（1-新增、2-更新、3-删除）',
			"create_time" datetime NOT NULL COMMENT '创建时间',
			"create_by" varchar(255) NOT NULL COMMENT '创建者',
			"update_time" datetime NOT NULL COMMENT '修改时间',
			"update_by" varchar(255) NOT NULL COMMENT '修改者',
			"release_by" varchar(255) CHARACTER SET utf8 DEFAULT NULL COMMENT '发布者',
			"release_time" datetime DEFAULT NULL COMMENT '发布时间',
			PRIMARY KEY ("config_id"),
			UNIQUE KEY "uniq_app_key" ("app_id","key") USING BTREE,
			KEY "idx_app_id" ("app_id"),
			KEY "idx_key" ("key")
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='配置明细表';
	`,
	"release": `
		CREATE TABLE IF NOT EXISTS "release" (
			"app_id" bigint(20) NOT NULL COMMENT '应用ID',
			"config_list" longtext CHARACTER SET utf8mb4 NOT NULL COMMENT '配置列表',
			"release_time" datetime NOT NULL COMMENT '修改时间',
			"release_index" int(11) NOT NULL COMMENT '发布序号',
			PRIMARY KEY ("app_id") USING BTREE,
			KEY "index_app_id" ("app_id")
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='发布版本表';
	`,
	"release_log": `
		CREATE TABLE IF NOT EXISTS "release_log" (
			"id" bigint(20) NOT NULL AUTO_INCREMENT COMMENT '配置ID',
			"app_id" bigint(20) NOT NULL COMMENT '应用ID',
			"config_list" longtext CHARACTER SET utf8mb4 NOT NULL COMMENT '配置列表',
			"release_time" datetime NOT NULL COMMENT '发布时间',
			"release_index" int(11) NOT NULL COMMENT '发布序号',
			"release_by" varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '发布者',
			PRIMARY KEY ("id"),
			KEY "index_app_id" ("app_id")
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='版本历史表';
	`,
	"user": `
		CREATE TABLE IF NOT EXISTS "user" (
			"user_id" bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
			"name" varchar(255) NOT NULL COMMENT '用户名',
			"password" varchar(255) NOT NULL COMMENT '用户密码',
			"permission" tinyint(4) NOT NULL DEFAULT '1' COMMENT '权限（1-普通用户、2-管理员）',
			"create_time" datetime NOT NULL COMMENT '创建时间',
			"update_time" datetime NOT NULL COMMENT '修改时间',
			PRIMARY KEY ("user_id"),
			UNIQUE KEY "uniq_name" ("name")
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';
	`,
}

var InsertDefaultUserSql string = "INSERT IGNORE INTO `user` VALUES (1, 'admin', '123456', 2, '2024-01-01 23:33:33', '2024-01-01 23:33:33');"

func init() {
	for k, v := range CreateTableSql {
		CreateTableSql[k] = strings.ReplaceAll(v, "\"", "`")
	}
}
