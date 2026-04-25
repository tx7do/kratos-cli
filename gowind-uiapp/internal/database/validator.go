package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// TestConnection 测试数据库连接并返回详细信息
func TestConnection(cfg DBConfig) (*ConnectionResult, error) {
	start := time.Now()

	// 建立连接
	conn, err := Connect(cfg)
	if err != nil {
		return &ConnectionResult{
			Success:  false,
			Message:  fmt.Sprintf("连接到 %s 数据库失败: %v", cfg.Type, err),
			Duration: time.Since(start).Milliseconds(),
			Error:    err.(*DBError),
		}, nil
	}
	defer conn.Close()

	// 获取数据库版本
	var version string
	switch cfg.Type {
	case DbTypeMySQL:
		err = conn.db.QueryRow("SELECT VERSION()").Scan(&version)
	case DbTypePostgreSQL:
		err = conn.db.QueryRow("SELECT version()").Scan(&version)
	case DbTypeSQLite:
		err = conn.db.QueryRow("SELECT sqlite_version()").Scan(&version)
	case DbTypeOracle:
		err = conn.db.QueryRow("SELECT * FROM v$version WHERE banner LIKE 'Oracle%'").Scan(&version)
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("获取版本失败: %v", err)
	}

	// 获取表数量
	tableCount, _ := getTableCount(conn, cfg.Type)

	duration := time.Since(start).Milliseconds()

	return &ConnectionResult{
		Success:   true,
		Message:   fmt.Sprintf("成功连接到 %s 数据库", cfg.Type),
		Database:  cfg.Database,
		ServerVer: cleanVersion(version),
		Duration:  duration,
		Tables:    tableCount,
		Connected: true,
	}, nil
}

// 获取表数量
func getTableCount(conn *DBConnection, dbType DbType) (int, error) {
	var query string
	switch dbType {
	case DbTypeMySQL:
		query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()"
	case DbTypePostgreSQL:
		query = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'"
	case DbTypeSQLite:
		query = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'"
	case DbTypeOracle:
		query = "SELECT COUNT(*) FROM user_tables"
	}

	var count int
	err := conn.db.QueryRow(query).Scan(&count)
	return count, err
}

// 清理版本字符串（移除冗余信息）
func cleanVersion(raw string) string {
	if raw == "" {
		return "Unknown"
	}
	// 移除 MySQL 的冗长版本信息
	if strings.Contains(strings.ToLower(raw), "mysql") {
		parts := strings.Split(raw, " ")
		if len(parts) > 0 {
			return parts[0]
		}
	}
	return raw
}
