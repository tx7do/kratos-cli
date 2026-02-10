module app

go 1.25.6

replace (
	github.com/tx7do/kratos-cli/generators => ../generators
	github.com/tx7do/kratos-cli/gowind => ../gowind
	github.com/tx7do/kratos-cli/sql-kratos => ../sql-kratos
	github.com/tx7do/kratos-cli/sql-orm => ../sql-orm
	github.com/tx7do/kratos-cli/sql-proto => ../sql-proto
)

require (
	github.com/go-sql-driver/mysql v1.9.3
	github.com/jackc/pgx/v5 v5.8.0
	github.com/sijms/go-ora/v2 v2.9.0
	github.com/tx7do/go-utils/ddl_parser v0.0.3
	github.com/tx7do/kratos-cli/sql-kratos v0.0.11
	github.com/wailsapp/wails/v2 v2.11.0
	modernc.org/sqlite v1.45.0
)
