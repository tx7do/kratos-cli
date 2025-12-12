package pkg

import "testing"

func TestNewModuleInspectorFromGo(t *testing.T) {
	inspector, err := NewModuleInspectorFromGo("")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Root: %s, ModPath: %s", inspector.Root, inspector.ModPath)
}
