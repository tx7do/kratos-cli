package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// GetTables 获取数据库表列表
func GetTables(conn *DBConnection, dbType DbType) ([]TableInfo, error) {
	var query string
	switch dbType {
	case DbTypeMySQL:
		query = `
			SELECT 
				t.TABLE_NAME AS table_name,
				t.TABLE_TYPE AS table_type,
				t.ENGINE AS table_engine,
				t.TABLE_ROWS AS table_rows,
				t.CREATE_TIME AS create_time,
				t.TABLE_COMMENT AS table_comment,
				(SELECT COUNT(*) FROM information_schema.COLUMNS c 
				 WHERE c.TABLE_SCHEMA = t.TABLE_SCHEMA AND c.TABLE_NAME = t.TABLE_NAME) AS table_columns,
				(SELECT COUNT(*) FROM information_schema.STATISTICS s 
				 WHERE s.TABLE_SCHEMA = t.TABLE_SCHEMA AND s.TABLE_NAME = t.TABLE_NAME) AS table_indexes
			FROM information_schema.TABLES t
			WHERE t.TABLE_SCHEMA = DATABASE()
			ORDER BY t.TABLE_NAME
		`
	case DbTypePostgreSQL:
		query = `
			SELECT 
				t.table_name AS table_name,
				t.table_type AS table_type,
				NULL AS table_engine,
				NULL AS table_rows,
				NULL AS create_time,
				obj_description(c.oid) AS table_comment,
				(SELECT COUNT(*) FROM information_schema.columns 
				 WHERE table_name = t.table_name AND table_schema = 'public') AS table_columns,
				(SELECT COUNT(*) FROM pg_indexes 
				 WHERE tablename = t.table_name AND schemaname = 'public') AS table_indexes
			FROM information_schema.tables t
			LEFT JOIN pg_class c ON c.relname = t.table_name
			WHERE t.table_schema = 'public' AND t.table_type = 'BASE TABLE'
			ORDER BY t.table_name
		`
	case DbTypeSQLite:
		query = `
			SELECT 
				name AS table_name,
				'table' AS table_type,
				NULL AS table_engine,
				0 AS table_rows,
				NULL AS create_time,
				NULL AS table_comment,
				(SELECT COUNT(*) FROM pragma_table_info(name)) AS table_columns,
				(SELECT COUNT(*) FROM pragma_index_list(name)) AS table_indexes
			FROM sqlite_master 
			WHERE type='table' AND name NOT LIKE 'sqlite_%'
			ORDER BY name
		`
	case DbTypeOracle:
		query = `
			SELECT 
				t.table_name AS table_name,
				'BASE TABLE' AS table_type,
				NULL AS table_engine,
				num_rows AS table_rows,
				last_analyzed AS create_time,
				NULL AS table_comment,
				(SELECT COUNT(*) FROM all_tab_columns c 
				 WHERE c.table_name = t.table_name AND c.owner = USER) AS table_columns,
				(SELECT COUNT(*) FROM all_indexes i 
				 WHERE i.table_name = t.table_name AND i.owner = USER) AS table_indexes
			FROM user_tables t
			ORDER BY t.table_name
		`
	}

	rows, err := conn.db.Query(query)
	if err != nil {
		return nil, wrapDBError(err, dbType)
	}
	defer rows.Close()

	//log.Printf("Executing table metadata query for %s database\n", dbType)

	var tables []TableInfo
	for rows.Next() {
		var table TableInfo
		var tableEngine sql.NullString
		var tableRows sql.NullInt64
		var createTime sql.NullTime
		var tableComment sql.NullString

		err = rows.Scan(
			&table.Name,
			&table.Type,
			&tableEngine,
			&tableRows,
			&createTime,
			&tableComment,
			&table.Columns,
			&table.Indexes,
		)
		if err != nil {
			log.Printf("Error scanning table row: %v\n", err)
			continue
		}

		if tableEngine.Valid {
			table.Engine = tableEngine.String
		}
		if tableRows.Valid {
			table.Rows = uint64(tableRows.Int64)
		}
		if createTime.Valid {
			table.CreateTime = createTime.Time
		}
		if tableComment.Valid {
			table.Comment = tableComment.String
		}

		//log.Printf("Fetched table: %+v\n", table)

		tables = append(tables, table)
	}

	return tables, rows.Err()
}

// GetColumns 获取表的列信息
func GetColumns(conn *DBConnection, dbType DbType, tableName string) ([]ColumnInfo, error) {
	var query string
	args := []any{tableName}

	switch dbType {
	case DbTypeMySQL:
		query = `
			SELECT 
				COLUMN_NAME,
				CONCAT(COLUMN_TYPE, IF(CHARACTER_MAXIMUM_LENGTH IS NOT NULL AND CHARACTER_MAXIMUM_LENGTH > 0, 
					CONCAT('(', CHARACTER_MAXIMUM_LENGTH, ')'), '')) AS full_type,
				IS_NULLABLE = 'YES' AS is_nullable,
				COLUMN_KEY = 'PRI' AS is_primary,
				COLUMN_DEFAULT,
				COLUMN_COMMENT,
				EXTRA
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
			ORDER BY ORDINAL_POSITION
		`
	case DbTypePostgreSQL:
		query = `
			SELECT 
				a.attname AS column_name,
				pg_catalog.format_type(a.atttypid, a.atttypmod) AS full_type,
				NOT a.attnotnull AS is_nullable,
				EXISTS (
					SELECT 1 FROM pg_constraint 
					WHERE conrelid = a.attrelid AND conkey[1] = a.attnum AND contype = 'p'
				) AS is_primary,
				pg_get_expr(d.adbin, d.adrelid) AS column_default,
				col_description(a.attrelid, a.attnum) AS column_comment,
				CASE WHEN a.attidentity != '' THEN 'GENERATED' ELSE '' END AS extra
			FROM pg_attribute a
			LEFT JOIN pg_attrdef d ON a.attrelid = d.adrelid AND a.attnum = d.adnum
			WHERE a.attrelid = $1::regclass AND a.attnum > 0 AND NOT a.attisdropped
			ORDER BY a.attnum
		`
		args = []any{tableName}
	case DbTypeSQLite:
		query = fmt.Sprintf("PRAGMA table_info(%s)", quoteSQLiteIdentifier(tableName))
		// SQLite 需要特殊处理
		return getSQLiteColumns(conn, tableName)
	case DbTypeOracle:
		query = `
			SELECT 
				column_name,
				data_type || 
					CASE WHEN data_precision IS NOT NULL THEN '(' || data_precision || 
						CASE WHEN data_scale > 0 THEN ',' || data_scale ELSE '' END || ')' 
					ELSE '' END AS full_type,
				nullable = 'Y' AS is_nullable,
				EXISTS (
					SELECT 1 FROM all_constraints c 
					JOIN all_cons_columns cc ON c.constraint_name = cc.constraint_name
					WHERE c.constraint_type = 'P' AND cc.table_name = :1 AND cc.column_name = column_name
				) AS is_primary,
				data_default,
				NULL AS column_comment,
				CASE WHEN column_id = 1 THEN 'FIRST' ELSE '' END AS extra
			FROM all_tab_columns
			WHERE table_name = :1 AND owner = USER
			ORDER BY column_id
		`
		args = []any{tableName}
	}

	rows, err := conn.db.Query(query, args...)
	if err != nil {
		return nil, wrapDBError(err, dbType)
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		var nullable, primaryKey sql.NullBool
		var colDefault, comment, extra sql.NullString

		err = rows.Scan(
			&col.Name,
			&col.Type,
			&nullable,
			&primaryKey,
			&colDefault,
			&comment,
			&extra,
		)
		if err != nil {
			continue
		}

		col.Nullable = nullable.Bool
		col.PrimaryKey = primaryKey.Bool
		if colDefault.Valid {
			col.Default = colDefault.String
		}
		if comment.Valid {
			col.Comment = comment.String
		}
		if extra.Valid {
			col.Extra = extra.String
		}

		columns = append(columns, col)
	}

	return columns, rows.Err()
}

// SQLite 特殊处理
func getSQLiteColumns(conn *DBConnection, tableName string) ([]ColumnInfo, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", quoteSQLiteIdentifier(tableName))
	rows, err := conn.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dfltValue sql.NullString

		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			continue
		}

		columns = append(columns, ColumnInfo{
			Name:       name,
			Type:       typ,
			Nullable:   notnull == 0,
			PrimaryKey: pk > 0,
			Default:    dfltValue.String,
			Comment:    "",
			Extra:      "",
		})
	}
	return columns, rows.Err()
}

// 安全转义 SQLite 标识符
func quoteSQLiteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
