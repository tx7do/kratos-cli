package internal

// Importer 远程配置导入器
type Importer interface {
	// Import 导入所有的配置
	Import() error

	// ImportOneService 导入单个配置
	ImportOneService(app string) error
}
