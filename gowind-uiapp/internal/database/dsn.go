package database

import (
	"fmt"
	"net/url"
	"strings"
)

// BuildDSN 根据配置智能构建 DSN
func BuildDSN(cfg DBConfig) (string, error) {
	// 优先使用自定义 DSN（高级模式）
	if cfg.UseDSN && cfg.DSN != "" {
		return cfg.DSN, nil
	}

	// 验证必要字段
	if cfg.Type == "" {
		return "", fmt.Errorf("数据库类型不能为空")
	}
	if cfg.Type != DbTypeSQLite && cfg.Host == "" {
		return "", fmt.Errorf("主机地址不能为空")
	}

	// 根据数据库类型构建 DSN
	switch cfg.Type {
	case DbTypeMySQL:
		return buildMySQLDSN(cfg)
	case DbTypePostgreSQL:
		return buildPostgresDSN(cfg)
	case DbTypeSQLite:
		return buildSQLiteDSN(cfg)
	case DbTypeOracle:
		return buildOracleDSN(cfg)
	default:
		return "", fmt.Errorf("不支持的数据库类型: %s", cfg.Type)
	}
}

// MySQL DSN: user:pass@tcp(host:port)/db?params
func buildMySQLDSN(cfg DBConfig) (string, error) {
	if cfg.Username == "" {
		return "", fmt.Errorf("MySQL 用户名不能为空")
	}
	// 密码可为空（如本地免密登录）
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}
	port := cfg.Port
	if port == 0 {
		port = 3306
	}

	// 构建参数
	params := []string{"charset=utf8mb4", "parseTime=true", "loc=Local"}
	if cfg.SSL {
		params = append(params, "tls=skip-verify") // 生产环境应使用真实证书
	} else {
		params = append(params, "tls=false")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.Username,
		cfg.Password, // 密码仅在此处拼接
		host,
		port,
		cfg.Database,
		strings.Join(params, "&"),
	), nil
}

// PostgreSQL DSN: postgres://user:pass@host:port/db?sslmode=disable
func buildPostgresDSN(cfg DBConfig) (string, error) {
	if cfg.Username == "" {
		return "", fmt.Errorf("PostgreSQL 用户名不能为空")
	}
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}
	port := cfg.Port
	if port == 0 {
		port = 5432
	}

	sslMode := "disable"
	if cfg.SSL {
		sslMode = "require"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&timezone=Asia/Shanghai",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password), // URL 编码避免特殊字符问题
		host,
		port,
		cfg.Database,
		sslMode,
	), nil
}

// SQLite DSN
func buildSQLiteDSN(cfg DBConfig) (string, error) {
	if cfg.DBPath == "" {
		return "", fmt.Errorf("SQLite 数据库路径不能为空")
	}
	// 支持内存数据库
	if cfg.DBPath == ":memory:" {
		return "file::memory:?cache=shared", nil
	}
	return fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", cfg.DBPath), nil
}

// Oracle DSN
func buildOracleDSN(cfg DBConfig) (string, error) {
	if cfg.Username == "" {
		return "", fmt.Errorf("oracle 用户名不能为空")
	}
	host := cfg.Host
	if host == "" {
		host = "localhost"
	}
	port := cfg.Port
	if port == 0 {
		port = 1521
	}
	serviceName := cfg.Database
	if serviceName == "" {
		serviceName = "ORCL"
	}

	return fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password),
		host,
		port,
		serviceName,
	), nil
}
