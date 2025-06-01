package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"entgo.io/contrib/schemast"

	"entgo.io/ent/dialect"

	_ "github.com/lib/pq"
)

// SchemaConverter is the interface that wraps the SchemaMutations method.
type SchemaConverter interface {
	// SchemaMutations imports a given schema from a data source and returns a list of schemast mutators.
	SchemaMutations(context.Context) ([]schemast.Mutator, error)
}

func NewConvert(opts ...ConvertOption) (SchemaConverter, error) {
	var (
		si  SchemaConverter
		err error
	)
	i := &ConvertOptions{}
	for _, apply := range opts {
		apply(i)
	}

	switch i.driver.Dialect {
	case dialect.MySQL:

	case dialect.Postgres:

	case "text":

	default:
		return nil, fmt.Errorf("entimport: unsupported dialect %q", i.driver.Dialect)
	}

	return si, err
}

func Convert(ctx context.Context, drv, dsn, schemaPath, outputPath *string, tables, excludeTables []string) {
	_ = os.MkdirAll(*outputPath, os.ModePerm)

	// 连接数据库
	db, err := sql.Open("postgres", "postgres://user:pass@localhost:5432/mydb?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	// 创建 Atlas 驱动
	driver, err := postgres.Open(db)
	if err != nil {
		log.Fatalf("failed to open driver: %v", err)
	}

	// 获取数据库模式
	sch, err := driver.InspectSchema(ctx, "public", &schema.InspectOptions{
		Tables: []string{"users"}, // 指定要检查的表，留空则获取所有表
	})
	if err != nil {
		log.Fatalf("failed to inspect schema: %v", err)
	}

	// 打印表结构
	for _, table := range sch.Tables {
		fmt.Printf("表名: %s\n", table.Name)
		fmt.Println("列信息:")
		//for _, col := range table.Columns {
		//	fmt.Printf("  - %s (%s)\n", col.Name, col.Type)
		//}
		//fmt.Println("索引:")
		//for _, idx := range table.Indexes {
		//	fmt.Printf("  - %s (unique: %v, columns: %v)\n",
		//		idx.Name, idx.Unique, cols)
		//}
		//fmt.Println("------------------------")
	}

}
