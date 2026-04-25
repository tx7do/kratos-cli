package database

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestConnect_WithDSN(t *testing.T) {
	// 测试使用自定义 DSN 连接
	cfg := DBConfig{
		Type:    DbTypeSQLite,
		UseDSN:  true,
		DSN:     ":memory:",
		Timeout: 5 * time.Second,
	}

	conn, err := Connect(cfg)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 验证连接有效
	if err := conn.Ping(context.Background()); err != nil {
		t.Fatalf("Ping 失败: %v", err)
	}

	t.Log("✓ 使用自定义 DSN 连接成功")
}

func TestConnect_WithConfig(t *testing.T) {
	// 测试使用配置字段构建 DSN 连接
	cfg := DBConfig{
		Type:    DbTypeSQLite,
		UseDSN:  false,
		DBPath:  ":memory:",
		Timeout: 5 * time.Second,
	}

	conn, err := Connect(cfg)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 验证连接有效
	if err := conn.Ping(context.Background()); err != nil {
		t.Fatalf("Ping 失败: %v", err)
	}

	t.Log("✓ 使用配置字段连接成功")
}

func TestConnect_InvalidDSN(t *testing.T) {
	// 测试无效的 DSN
	cfg := DBConfig{
		Type:   DbTypeMySQL,
		UseDSN: true,
		DSN:    "invalid dsn",
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("期望返回错误，但没有错误")
	}

	var dbErr *DBError
	ok := errors.As(err, &dbErr)
	if !ok {
		t.Fatalf("期望返回 DBError，但得到 %T", err)
	}

	t.Logf("✓ 正确返回 DBError: %s", dbErr.Error())
}

func TestConnect_EmptyDSN(t *testing.T) {
	// 测试 DSN 为空的情况
	cfg := DBConfig{
		Type:   DbTypeMySQL,
		UseDSN: true,
		DSN:    "",
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("期望返回错误，但没有错误")
	}

	var dbErr *DBError
	ok := errors.As(err, &dbErr)
	if !ok {
		t.Fatalf("期望返回 DBError")
	}

	if dbErr.Code != "INVALID_CONFIG" {
		t.Errorf("期望错误码 INVALID_CONFIG，但得到 %s", dbErr.Code)
	}

	t.Logf("✓ 正确处理空 DSN: %s", dbErr.Error())
}

func TestConnect_UnsupportedType(t *testing.T) {
	// 测试不支持的数据库类型
	cfg := DBConfig{
		Type:   DbType("unsupported"),
		UseDSN: true,
		DSN:    "some dsn",
	}

	_, err := Connect(cfg)
	if err == nil {
		t.Fatal("期望返回错误，但没有错误")
	}

	var dbErr *DBError
	ok := errors.As(err, &dbErr)
	if !ok {
		t.Fatalf("期望返回 DBError")
	}

	if dbErr.Code != "UNSUPPORTED_DB" {
		t.Errorf("期望错误码 UNSUPPORTED_DB，但得到 %s", dbErr.Code)
	}

	t.Logf("✓ 正确处理不支持的数据库类型: %s", dbErr.Error())
}

func TestGenerateConnKey(t *testing.T) {
	// 测试使用自定义 DSN 时的连接键生成
	cfg1 := DBConfig{
		Type:     DbTypeMySQL,
		UseDSN:   true,
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		Username: "root",
	}
	key1 := generateConnKey(cfg1)
	if key1 != "mysql://custom-dsn" {
		t.Errorf("期望 'mysql://custom-dsn'，但得到 '%s'", key1)
	}
	t.Logf("✓ UseDSN=true 时的连接键: %s", key1)

	// 测试使用配置字段时的连接键生成
	cfg2 := DBConfig{
		Type:     DbTypePostgreSQL,
		UseDSN:   false,
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		Username: "admin",
	}
	key2 := generateConnKey(cfg2)
	if key2 != "postgresql://admin@localhost:5432/testdb" {
		t.Errorf("期望 'postgresql://admin@localhost:5432/testdb'，但得到 '%s'", key2)
	}
	t.Logf("✓ UseDSN=false 时的连接键: %s", key2)
}

func TestBuildDSN_MySQL(t *testing.T) {
	cfg := DBConfig{
		Type:     DbTypeMySQL,
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
		Username: "root",
		Password: "password",
		SSL:      false,
	}

	dsn, err := buildDSN(cfg)
	if err != nil {
		t.Fatalf("buildDSN 失败: %v", err)
	}

	if dsn == "" {
		t.Fatal("DSN 为空")
	}

	t.Logf("✓ MySQL DSN: %s", dsn)
}

func TestBuildDSN_PostgreSQL(t *testing.T) {
	cfg := DBConfig{
		Type:     DbTypePostgreSQL,
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		Username: "postgres",
		Password: "password",
		SSL:      true,
	}

	dsn, err := buildDSN(cfg)
	if err != nil {
		t.Fatalf("buildDSN 失败: %v", err)
	}

	if dsn == "" {
		t.Fatal("DSN 为空")
	}

	t.Logf("✓ PostgreSQL DSN: %s", dsn)
}

func TestBuildDSN_SQLite(t *testing.T) {
	cfg := DBConfig{
		Type:   DbTypeSQLite,
		DBPath: "/tmp/test.db",
	}

	dsn, err := buildDSN(cfg)
	if err != nil {
		t.Fatalf("buildDSN 失败: %v", err)
	}

	if dsn == "" {
		t.Fatal("DSN 为空")
	}

	t.Logf("✓ SQLite DSN: %s", dsn)
}

func TestBuildDSN_SQLite_NoPath(t *testing.T) {
	cfg := DBConfig{
		Type:   DbTypeSQLite,
		DBPath: "", // 缺失路径
	}

	_, err := buildDSN(cfg)
	if err == nil {
		t.Fatal("期望返回错误，但没有错误")
	}

	t.Logf("✓ 正确处理缺失的 SQLite 路径: %v", err)
}

func TestBuildDSN_Oracle(t *testing.T) {
	cfg := DBConfig{
		Type:     DbTypeOracle,
		Host:     "localhost",
		Port:     1521,
		Database: "ORCL",
		Username: "system",
		Password: "password",
	}

	dsn, err := buildDSN(cfg)
	if err != nil {
		t.Fatalf("buildDSN 失败: %v", err)
	}

	if dsn == "" {
		t.Fatal("DSN 为空")
	}

	t.Logf("✓ Oracle DSN: %s", dsn)
}

func TestConnection_Ping(t *testing.T) {
	// 创建内存数据库连接
	cfg := DBConfig{
		Type:   DbTypeSQLite,
		UseDSN: true,
		DSN:    ":memory:",
	}

	conn, err := Connect(cfg)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 测试 Ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.Ping(ctx); err != nil {
		t.Fatalf("Ping 失败: %v", err)
	}

	t.Log("✓ Ping 成功")
}

func TestConnection_Close(t *testing.T) {
	cfg := DBConfig{
		Type:   DbTypeSQLite,
		UseDSN: true,
		DSN:    ":memory:",
	}

	conn, err := Connect(cfg)
	if err != nil {
		t.Fatalf("连接失败: %v", err)
	}

	conn.Close()
	t.Log("✓ Close 成功")
}
