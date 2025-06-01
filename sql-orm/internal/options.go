package internal

type Options struct {
	ORM string `json:"orm"` // ORM type, e.g., "gorm", "sqlx", "ent"

	Driver string `json:"driver"` // Driver name, e.g., "mysql", "postgres"
	Source string `json:"source"` // Data source name (DSN), e.g., "mysql://user:pass@tcp(localhost:3306)/dbname"

	Tables         []string `json:"tables"`
	ExcludedTables []string `json:"excluded_tables"` // Tables to exclude from inspection

	SchemaPath string `json:"schema_path"` // Path to save the generated schema files
	DaoPath    string `json:"dao_path"`    // Path to save the generated DAO code (for gorm)
}
