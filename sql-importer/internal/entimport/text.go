package entimport

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/google/uuid"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"

	"entgo.io/contrib/schemast"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

type Text struct {
	*ImportOptions
}

func NewText(i *ImportOptions) (*Text, error) {
	return &Text{
		ImportOptions: i,
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

func (t *Text) InspectSchema(ctx context.Context, sqlText string, opts *schema.InspectOptions, s *schema.Schema) (*schema.Schema, error) {
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

// SchemaMutations 实现 SchemaImporter 接口，用于解析 SQL 文本
func (t *Text) SchemaMutations(ctx context.Context) ([]schemast.Mutator, error) {
	// 加载 SQL 文本
	sqlText := t.loadSQLFromFile()
	if sqlText == "" {
		return nil, fmt.Errorf("无法加载 SQL 文件: %v", t.schemaPath)
	}

	inspectOptions := &schema.InspectOptions{
		Tables: t.tables,
	}

	var s schema.Schema

	// 分割 SQL 语句
	queries := strings.Split(sqlText, ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue // 跳过空查询
		}
		_, err := t.InspectSchema(ctx, query, inspectOptions, &s)
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

	return schemaMutations(t.field, tables)
}

func (t *Text) field(column *schema.Column) (f ent.Field, err error) {
	name := column.Name
	switch typ := column.Type.Type.(type) {
	case *schema.BinaryType:
		f = field.Bytes(name)
	case *schema.BoolType:
		f = field.Bool(name)
	case *schema.DecimalType:
		f = field.Float(name)
	case *schema.EnumType:
		f = field.Enum(name).Values(typ.Values...)
	case *schema.FloatType:
		f = t.convertFloat(typ, name)
	case *schema.IntegerType:
		f = t.convertInteger(typ, name)
	case *schema.JSONType:
		f = field.JSON(name, json.RawMessage{})
	case *schema.StringType:
		f = field.String(name)
	case *schema.TimeType:
		f = field.Time(name)

	case *postgres.SerialType:
		f = t.convertSerial(typ, name)
	case *postgres.UUIDType:
		f = field.UUID(name, uuid.New())

	default:
		return nil, fmt.Errorf("entimport: unsupported type %q for column %v", typ, column.Name)
	}
	applyColumnAttributes(f, column)
	return f, err
}

func (t *Text) convertFloat(typ *schema.FloatType, name string) (f ent.Field) {
	// Precision from 0 to 23 results in a 4-byte single-precision FLOAT column.
	// Precision from 24 to 53 results in an 8-byte double-precision DOUBLE column:
	// https://dev.mysql.com/doc/refman/8.0/en/floating-point-types.html
	switch typ.T {
	case mysql.TypeDouble:
		return field.Float(name)

	case mysql.TypeFloat:
		if typ.Precision == 0 {
			// If precision and scale are not specified, use Float32.
			return field.Float32(name)
		}
		// If precision is specified, use Float64.
		if typ.Precision > 23 {
			return field.Float(name)
		}

	case mysql.TypeReal:
		// MySQL's REAL is an alias for FLOAT, so we treat it as such.
		if typ.Precision == 0 {
			// If precision and scale are not specified, use Float32.
			return field.Float32(name)
		}
		// If precision is specified, use Float64.
		if typ.Precision > 23 {
			return field.Float(name)
		}
	}

	switch typ.T {
	case postgres.TypeReal:
		return field.Float32(name)
	case postgres.TypeDouble:
		return field.Float(name)
	case postgres.TypeFloat8:
		return field.Float(name)
	case postgres.TypeFloat4:
		return field.Float32(name)
	case postgres.TypeFloat:
		return field.Float(name)
	}

	return field.Float32(name)
}

func (t *Text) convertInteger(typ *schema.IntegerType, name string) (f ent.Field) {
	if typ.Unsigned {
		switch typ.T {
		case mTinyInt:
			f = field.Uint8(name)
		case mSmallInt:
			f = field.Uint16(name)
		case mMediumInt:
			f = field.Uint32(name)
		case mInt:
			f = field.Uint32(name)
		case mBigInt:
			f = field.Uint64(name)
		}
		return f
	}

	switch typ.T {
	case mTinyInt:
		f = field.Int8(name)
	case mSmallInt:
		f = field.Int16(name)
	case mMediumInt:
		f = field.Int32(name)
	case mInt:
		f = field.Int32(name)
	case mBigInt:
		// Int64 is not used on purpose.
		f = field.Int(name)

	case pInteger:
		f = field.Int32(name)
	}

	return f
}

// smallserial- 2 bytes - small autoincrementing integer 1 to 32767
// serial - 4 bytes autoincrementing integer 1 to 2147483647
// bigserial - 8 bytes large autoincrementing integer	1 to 9223372036854775807
func (t *Text) convertSerial(typ *postgres.SerialType, name string) ent.Field {
	return field.Uint(name).
		SchemaType(map[string]string{
			dialect.Postgres: typ.T, // Override Postgres.
		})
}
