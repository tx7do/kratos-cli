package pkg

import (
	"bytes"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ReplaceInFile 读取文件 `path`，将所有 `old` 替换为 `new`，并原子写回。
// 如果文件内容没有变化，不会修改文件时间。
func ReplaceInFile(path, old, new string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s := string(b)
	newS := strings.ReplaceAll(s, old, new)
	if newS == s {
		return nil
	}
	return writeAtomic(path, []byte(newS))
}

// ReplaceRegexInFile 使用正则 `pattern` 将匹配的部分替换为 `replacement` 并写回。
// `pattern` 是普通的 Go 正则表达式。
func ReplaceRegexInFile(path, pattern, replacement string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	newB := re.ReplaceAll(b, []byte(replacement))
	if string(newB) == string(b) {
		return nil
	}
	return writeAtomic(path, newB)
}

// writeAtomic 在与目标相同目录下创建临时文件写入数据，然后重命名覆盖目标，保留原权限。
func writeAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)

	// 尝试获取原始文件权限
	var perm os.FileMode = 0o644
	if fi, err := os.Stat(path); err == nil {
		perm = fi.Mode().Perm()
	}

	tmp, err := os.CreateTemp(dir, "tmpfile-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	_, err = tmp.Write(data)
	if err1 := tmp.Close(); err == nil {
		err = err1
	}
	if err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	if err = os.Chmod(tmpName, perm); err != nil {
		_ = os.Remove(tmpName)
		return err
	}

	// 原子替换
	if err = os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}

func isFileExtIncluded(path string, includeExts []string) bool {
	ext := filepath.Ext(path)
	for _, e := range includeExts {
		if ext == "."+e {
			return true
		}
	}
	return false
}

// ReplaceTemplateInCurrentDir 遍历当前目录及子目录，替换指定的模板字符串。
// 返回被修改的文件数量和错误（如果有）。
func ReplaceTemplateInCurrentDir(rootDir, source, target string) (int, error) {
	var includeFileExt = []string{"go", "mod", "yaml"}

	var updated int
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//log.Printf("scanning: %s\n", path)

		// 跳过 .git 和 vendor 目录
		if d.IsDir() {
			switch d.Name() {
			case ".git", "vendor", "script":
				return filepath.SkipDir
			default:
				return nil
			}
		}

		if !isFileExtIncluded(path, includeFileExt) {
			//log.Printf("skipped (ext): %s\n", path)
			return nil
		}

		info, err := d.Info()
		if err != nil {
			log.Printf("failed to get file info: %s, error: %v\n", path, err)
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("failed to read file: %s, error: %v\n", path, err)
			return err
		}

		if !bytes.Contains(data, []byte(source)) {
			log.Printf("skipped: %s\n", path)
			return nil
		}

		newData := bytes.ReplaceAll(data, []byte(source), []byte(target))
		if err = os.WriteFile(path, newData, info.Mode().Perm()); err != nil {
			log.Printf("failed to write file: %s, error: %v\n", path, err)
			return err
		}

		updated++
		log.Printf("updated: %s\n", path)
		return nil
	})

	return updated, err
}

func IsDirExists(dir string) bool {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return true
	}
	return false
}

func IsFileExists(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err == nil {
		return !fi.IsDir()
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
