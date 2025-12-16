package pkg

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func kratosHome() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	home := filepath.Join(dir, ".kratos")
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func kratosHomeWithDir(dir string) string {
	home := filepath.Join(kratosHome(), dir)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		if err := os.MkdirAll(home, 0o700); err != nil {
			log.Fatal(err)
		}
	}
	return home
}

func copyFile(src, dst string, replaces []string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	var old string
	for i, next := range replaces {
		if i%2 == 0 {
			old = next
			continue
		}
		buf = bytes.ReplaceAll(buf, []byte(old), []byte(next))
	}
	return os.WriteFile(dst, buf, srcinfo.Mode())
}

func copyDir(src, dst string, replaces, ignores []string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, srcinfo.Mode())
	if err != nil {
		return err
	}

	fds, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		if hasSets(fd.Name(), ignores) {
			continue
		}
		srcfp := filepath.Join(src, fd.Name())
		dstfp := filepath.Join(dst, fd.Name())
		var e error
		if fd.IsDir() {
			e = copyDir(srcfp, dstfp, replaces, ignores)
		} else {
			e = copyFile(srcfp, dstfp, replaces)
		}
		if e != nil {
			return e
		}
	}
	return nil
}

func hasSets(name string, sets []string) bool {
	for _, ig := range sets {
		if ig == name {
			return true
		}
	}
	return false
}

func Tree(path string, dir string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err == nil && info != nil && !info.IsDir() {
			fmt.Printf("%s %s (%v bytes)\n", color.GreenString("CREATED"), strings.ReplaceAll(path, dir+"/", ""), info.Size())
		}
		return nil
	})
}

func SplitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}

// HasCmdAndConfigs 检查 dir（为空则使用工作目录）下是否存在 cmd 和 configs 目录。
// 返回 (hasCmd, hasConfigs, error)。
func HasCmdAndConfigs(dir string) (bool, bool, error) {
	var d string
	if len(dir) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			return false, false, err
		}
		d = wd
	} else {
		d = dir
	}

	hasMain := false
	mainPath := filepath.Join(d, "cmd", "server", "main.go")
	if fi, err := os.Stat(mainPath); err == nil {
		if !fi.IsDir() {
			hasMain = true
		}
	} else if !os.IsNotExist(err) {
		return false, false, err
	}

	hasConfigs := false
	configsPath := filepath.Join(d, "configs")
	if fi, err := os.Stat(configsPath); err == nil {
		if fi.IsDir() {
			hasConfigs = true
		}
	} else if !os.IsNotExist(err) {
		return false, false, err
	}

	return hasMain, hasConfigs, nil
}
