package internal

import (
	"context"
	"fmt"
	"os"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

type Text struct {
	*ConvertOptions
}

func NewText(i *ConvertOptions) (*Text, error) {
	return &Text{
		ConvertOptions: i,
	}, nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func (t *Text) loadSQLFromFile() string {
	// 检查 schemaPath 是否为文件
	if !isFile(t.schemaPath) {
		// 如果不是文件，直接返回 schemaPath，认为它是一个 SQL 文本内容。
		return t.schemaPath
	}

	content, err := os.ReadFile(t.schemaPath)
	if err != nil {
		return ""
	}

	return string(content)
}

func (t *Text) ParseType(raw string) (schema.Type, error) {
	mysqlType, err := mysql.ParseType(raw)
	if err == nil {
		return mysqlType, nil
	}

	postgresType, err := postgres.ParseType(raw)
	if err == nil {
		return postgresType, nil
	}

	sqliteType, err := sqlite.ParseType(raw)
	if err == nil {
		return sqliteType, nil
	}

	return &schema.UnsupportedType{T: raw}, nil
}

func (t *Text) toColumnType(raw string, options []*sqlparser.ColumnOption) (*schema.ColumnType, error) {
	parsedType, err := t.ParseType(raw)
	if err != nil {
		return nil, err
	}

	isNull := true
	if options != nil {
		for _, option := range options {
			if option.Type == sqlparser.ColumnOptionNotNull {
				isNull = false
				break // 一旦找到 NOT NULL 选项，就可以停止检查
			}
		}
	}

	return &schema.ColumnType{
		Type: parsedType,
		Raw:  raw,
		Null: isNull,
	}, nil
}

func (t *Text) InspectSchema(sqlText string, s *schema.Schema) (*schema.Schema, error) {
	// 解析 SQL 文本
	orgStmt, err := sqlparser.Parse(sqlText)
	if err != nil {
		return nil, fmt.Errorf("解析失败: %v", err)
	}
	//fmt.Printf("stmt = %+v\n", orgStmt)

	switch stmt := orgStmt.(type) {
	case *sqlparser.CreateTable:
		//case *sqlparser.DDL:
		switch stmt.Action {
		case sqlparser.CreateStr:
			fmt.Println("解析到 CREATE TABLE 语句")
			fmt.Printf("表名: %v\n", stmt.NewName.Name.String())

			table := &schema.Table{
				Name:   stmt.NewName.Name.String(),
				Schema: s,
			}

			for _, opt := range stmt.Options {
				switch opt.Type {
				case sqlparser.TableOptionEngine:
					fmt.Printf("表引擎: %v\n", opt.String())

				case sqlparser.TableOptionCollate:
					fmt.Printf("表排序规则: %v\n", opt.String())
					table.Attrs = append(table.Attrs, &schema.Collation{
						V: opt.StrValue,
					})

				case sqlparser.TableOptionCharset:
					fmt.Printf("字符集: %v\n", opt.String())
					table.Attrs = append(table.Attrs, &schema.Charset{
						V: opt.StrValue,
					})

				case sqlparser.TableOptionComment:
					fmt.Printf("表注释: %v\n", opt.String())
					table.Attrs = append(table.Attrs, &schema.Comment{
						Text: opt.StrValue,
					})

				default:
					panic(fmt.Sprintf("unhandled default case %v", opt.Type))
				}
			}

			for _, col := range stmt.Columns {
				fmt.Printf("列名: %v, 类型: %v\n", col.Name, col.Type)

				colType, err := t.toColumnType(col.Type, col.Options)
				if err != nil {
					fmt.Printf("解析失败: %v\n", err)
					continue
				}
				//fmt.Printf("原始类型: %s, 解析后类型: %T %v\n", col.Type, colType.Type, colType.Null)

				column := &schema.Column{
					Name: col.Name,
					Type: colType,
				}

				table.Columns = append(table.Columns, column)
			}

			for _, idx := range stmt.Constraints {
				fmt.Printf("索引名: %v, 类型: %v\n", idx.Name, idx.Type)
				for _, key := range idx.Keys {
					fmt.Printf("索引列: %v 类型: %v\n", key.String(), idx.Type)

					switch idx.Type {
					case sqlparser.ConstraintPrimaryKey:

						var colFind *schema.Column
						for _, col := range table.Columns {
							if col.Name == key.String() {
								colFind = col
								continue
							}
						}

						if table.PrimaryKey == nil {
							table.PrimaryKey = &schema.Index{
								Table: table,
								Name:  idx.Name,
								Parts: []*schema.IndexPart{
									{
										C: colFind,
									},
								},
							}
						} else {
							table.PrimaryKey.Parts = append(table.PrimaryKey.Parts, &schema.IndexPart{
								C: colFind,
							})
						}

					case sqlparser.ConstraintForeignKey:
						table.ForeignKeys = append(table.ForeignKeys, &schema.ForeignKey{
							Table:  table,
							Symbol: idx.Name,
						})

					case sqlparser.ConstraintKey:
					case sqlparser.ConstraintUniq:

					case sqlparser.ConstraintUniqKey:

					case sqlparser.ConstraintIndex:
						table.Indexes = append(table.Indexes, &schema.Index{
							Table: table,
							Name:  idx.Name,
						})

					case sqlparser.ConstraintUniqIndex:
						table.Indexes = append(table.Indexes, &schema.Index{
							Table:  table,
							Name:   idx.Name,
							Unique: true,
						})

					default:
						panic("unhandled default case")
					}

				}
			}

			s.Tables = append(s.Tables, table)
		}
	}

	return s, nil
}

func (t *Text) SchemaTables(ctx context.Context) ([]*TableData, error) {
	// 加载 SQL 文本
	sqlText := t.loadSQLFromFile()
	if sqlText == "" {
		return nil, fmt.Errorf("无法加载 SQL 文件: %v", t.schemaPath)
	}

	//inspectOptions := &schema.InspectOptions{
	//	Tables: t.includedTables,
	//}

	var s schema.Schema

	// 分割 SQL 语句
	queries := strings.Split(sqlText, ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue // 跳过空查询
		}
		_, err := t.InspectSchema(query, &s)
		if err != nil {
			return nil, err
		}
	}

	tables := s.Tables
	if t.excludedTables != nil {
		tables = nil
		excludedTableNames := make(map[string]bool)
		for _, t := range t.excludedTables {
			excludedTableNames[t] = true
		}
		// filter out tables that are in excludedTables:
		for _, t := range s.Tables {
			if !excludedTableNames[t.Name] {
				tables = append(tables, t)
			}
		}
	}

	return schemaTables(t.fieldType, tables)
}

func (t *Text) fieldType(sqlType string) (f string) {
	if f = MySQLFieldType(sqlType); f != "" {
		return
	}

	f = PostgresFieldType(sqlType)

	return
}
