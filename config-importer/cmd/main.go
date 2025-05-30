package main

import "github.com/tx7do/kratos-cli/config-importer/internal"

func main() {
	i := internal.NewImporter()
	_ = i.Import()
}
