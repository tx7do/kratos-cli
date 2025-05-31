package gorm

import (
	"context"
	"testing"
)

func TestImporter(t *testing.T) {
	ctx := context.Background()
	drv := "postgres"
	//dsn := "mysql://root:*Abcd123456@tcp(localhost:3306)/example"
	dsn := "postgres://postgres:*Abcd123456@localhost:5432/example?sslmode=disable"
	schemaPath := "./models/"
	daoPath := "./daos/"
	var tables []string
	var excludeTables []string
	_ = Importer(ctx, &drv, &dsn, &schemaPath, &daoPath, tables, excludeTables)
}
