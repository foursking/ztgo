package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v2/config"
)

// InitFileConfig 初始化文件配置
func InitFileConfig(path string) (err error) {
	paths, err := readAllPaths(path)
	if err != nil {
		return
	}
	if len(paths) == 0 {
		return
	}
	for _, p := range paths {
		if err = config.LoadFile(p); err != nil {
			return
		}
	}
	return
}

func readAllPaths(base string) ([]string, error) {
	fs, err := os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("check local config file fail! error: %s", err)
	}
	// dir or file to paths
	if !fs.IsDir() {
		return []string{base}, nil
	}
	var paths []string
	files, err := ioutil.ReadDir(base)
	if err != nil {
		return nil, fmt.Errorf("read dir %s error: %s", base, err)
	}
	for _, f := range files {
		if !f.IsDir() && !isHiddenFile(f.Name()) {
			paths = append(paths, path.Join(base, f.Name()))
		}
	}
	return paths, nil
}

func isHiddenFile(name string) bool {
	// TODO: support windows
	return strings.HasPrefix(filepath.Base(name), ".")
}
