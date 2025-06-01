package internal

import "github.com/tx7do/kratos-cli/sql-proto/internal/mux"

type Options struct {
	Driver string `json:"driver"` // Driver name, e.g., "mysql", "postgres"
	Source string `json:"source"` // Data source name (DSN), e.g., "mysql://user:pass@tcp(localhost:3306)/dbname"

	Tables         []string `json:"tables"`
	ExcludedTables []string `json:"excluded_tables"` // Tables to exclude from inspection

	OutputPath string `json:"output_path"` // Path to save the generated code
}

type (
	// ConvertOptions are the options passed on to every SchemaConverter.
	ConvertOptions struct {
		tables         []string
		excludedTables []string
		protoPath      string
		driver         *mux.ConvertDriver
	}

	// ConvertOption allows for managing import configuration using functional options.
	ConvertOption func(*ConvertOptions)
)

// WithProtoPath provides a path for writing the generated protocol buffer definitions.
func WithProtoPath(path string) ConvertOption {
	return func(i *ConvertOptions) {
		i.protoPath = path
	}
}

// WithTables limits the schema import to a set of given tables (by all tables are imported)
func WithTables(tables []string) ConvertOption {
	return func(i *ConvertOptions) {
		i.tables = tables
	}
}

// WithExcludedTables supplies the set of tables to exclude.
func WithExcludedTables(tables []string) ConvertOption {
	return func(i *ConvertOptions) {
		i.excludedTables = tables
	}
}

// WithDriver provides an import driver to be used by SchemaConverter.
func WithDriver(drv *mux.ConvertDriver) ConvertOption {
	return func(i *ConvertOptions) {
		i.driver = drv
	}
}
