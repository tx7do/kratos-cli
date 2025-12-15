::指定起始文件夹
set DIR=%cd%

cd %DIR%\generators
go get all
go mod tidy

cd %DIR%\config-exporter
go get all
go mod tidy

cd %DIR%\sql-orm
go get all
go mod tidy

cd %DIR%\sql-proto
go get all
go mod tidy

cd %DIR%\sql-kratos
go get all
go mod tidy

cd %DIR%\gowind
go get all
go mod tidy
