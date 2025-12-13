package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceTemplateInCurrentDir(t *testing.T) {
	rootDir := "./go-wind-admin-template"
	newModuleName := "github.com/tx7do/my-wind-admin"
	oldModuleName := "github.com/tx7do/go-wind-admin-template"

	updateCount, err := ReplaceTemplateInCurrentDir(rootDir, oldModuleName, newModuleName)
	assert.Nil(t, err)
	t.Log("updateCount:", updateCount)
}
