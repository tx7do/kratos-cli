package database

import (
	"time"
)

// DbType 数据库类型枚举
type DbType string

const (
	DbTypeMySQL      DbType = "mysql"
	DbTypePostgreSQL DbType = "postgresql"
	DbTypeSQLite     DbType = "sqlite"
	DbTypeOracle     DbType = "oracle"
)

// DBConfig 数据库连接配置（前端传入）
type DBConfig struct {
	Type     DbType `json:"type"`     // mysql/postgresql/sqlite/oracle
	Host     string `json:"host"`     // 主机地址 (localhost)
	Port     int    `json:"port"`     // 端口 (3306)
	Database string `json:"database"` // 数据库名
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码（敏感信息）
	SSL      bool   `json:"ssl"`      // SSL 连接

	// SQLite 特有
	DBPath string `json:"dbPath"` // SQLite 文件路径

	// 高级模式
	UseDSN bool   `json:"useDSN,omitempty"` // 是否使用自定义 DSN
	DSN    string `json:"dsn,omitempty"`    // 自定义 DSN（仅当 UseDSN=true 时生效）

	SQLContent string `json:"sqlContent,omitempty"` // 导入的 SQL 内容

	// 连接参数（可选）
	Timeout      time.Duration `json:"timeout,omitempty"`      // 连接超时（秒）
	MaxOpenConns int           `json:"maxOpenConns,omitempty"` // 最大连接数
}

// ConnectionResult 连接测试结果
type ConnectionResult struct {
	Success   bool     `json:"success"`   // 是否成功
	Message   string   `json:"message"`   // 详细消息
	Database  string   `json:"database"`  // 数据库名
	ServerVer string   `json:"serverVer"` // 服务器版本
	Duration  int64    `json:"duration"`  // 耗时(毫秒)
	Tables    int      `json:"tables"`    // 表数量
	Connected bool     `json:"connected"` // 是否保持连接
	Error     *DBError `json:"error,omitempty"`
}

// DBError 标准化错误
type DBError struct {
	Code    string `json:"code"`    // 错误码 (AUTH_FAILED, CONN_TIMEOUT...)
	Message string `json:"message"` // 用户友好消息
	Details string `json:"details"` // 技术细节（开发用）
}

// Error 实现 error 接口
func (e *DBError) Error() string {
	if e == nil {
		return ""
	}
	if e.Details != "" {
		return e.Code + ": " + e.Message + " (" + e.Details + ")"
	}
	return e.Code + ": " + e.Message
}

// TableInfo 表元数据
type TableInfo struct {
	Name       string    `json:"table_name"`
	Type       string    `json:"table_type"`    // BASE TABLE, VIEW...
	Engine     string    `json:"table_engine"`  // InnoDB, MyISAM...
	Rows       uint64    `json:"table_rows"`    // 行数（估算）
	Comment    string    `json:"table_comment"` // 表注释
	Columns    int       `json:"table_columns"` // 列数量
	Indexes    int       `json:"table_indexes"` // 索引数量
	CreateTime time.Time `json:"create_time"`   // 创建时间
}

// ColumnInfo 列元数据
type ColumnInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`       // VARCHAR(255), INT...
	Nullable   bool   `json:"nullable"`   // 是否可为空
	PrimaryKey bool   `json:"primaryKey"` // 是否主键
	Default    string `json:"default"`    // 默认值
	Comment    string `json:"comment"`    // 列注释
	Extra      string `json:"extra"`      // AUTO_INCREMENT...
}

// QueryResult SQL 查询结果
type QueryResult struct {
	Columns      []string        `json:"columns"` // 列名
	Rows         [][]interface{} `json:"rows"`    // 行数据
	RowsAffected int64           `json:"rowsAffected"`
	Duration     int64           `json:"duration"` // 毫秒
}
