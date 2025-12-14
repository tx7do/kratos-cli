package project

import (
	"os"
	"path/filepath"
	"strings"
)

// getGoModProjectRoot gets the root directory of the Go module project.
func getGoModProjectRoot(dir string) string {
	if dir == filepath.Dir(dir) {
		return dir
	}
	if goModIsNotExistIn(dir) {
		return getGoModProjectRoot(filepath.Dir(dir))
	}
	return dir
}

// goModIsNotExistIn checks if go.mod file does not exist in the specified directory.
func goModIsNotExistIn(dir string) bool {
	_, e := os.Stat(filepath.Join(dir, "go.mod"))
	return os.IsNotExist(e)
}

// isDirExists checks if the directory exists.
func isDirExists(baseDir, projectName string) bool {
	to := filepath.Join(baseDir, projectName)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		return true
	}
	return false
}

// processProjectParams process project name and working dir.
func processProjectParams(projectName string, workingDir string) (projectNameResult, workingDirResult string) {
	_projectDir := projectName
	_workingDir := workingDir
	// Process ProjectModule with system variable
	if strings.HasPrefix(projectName, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// cannot get user home return fallback place dir
			return _projectDir, _workingDir
		}
		_projectDir = filepath.Join(homeDir, projectName[2:])
	}

	// check path is relative
	if !filepath.IsAbs(projectName) {
		absPath, err := filepath.Abs(projectName)
		if err != nil {
			return _projectDir, _workingDir
		}
		_projectDir = absPath
	}

	return filepath.Base(_projectDir), filepath.Dir(_projectDir)
}
