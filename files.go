package main

import (
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"strings"
)

type DirStr struct {
	Root   string
	Dirs   []string
	Files  []string
	Images []string
}

func NewDirstr(path string) *DirStr {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if !fi.IsDir() {
		return nil
	} else {
		dir := DirStr{Root: path}
		files, _ := ioutil.ReadDir(path)
		for _, fi := range files {
			//tmpDirPath := filepath.Join(path, fi.Name())
			tmpDirPath := fi.Name()
			if fi.IsDir() {
				dir.Dirs = append(dir.Dirs, tmpDirPath)
			} else {
				dir.Files = append(dir.Files, tmpDirPath)
				switch strings.ToLower(fp.Ext(tmpDirPath)) {
				case ".jpg", ".png", ".gif":
					dir.Images = append(dir.Images, tmpDirPath)
				default:
				}
			}
		}
		return &dir
	}
}

func fileSize(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	size := fileInfo.Size()
	const ratio = 1024
	switch {
	case size < ratio:
		return fmt.Sprintf("%v %v", size, "B")
	case size/ratio < ratio:
		return fmt.Sprintf("%v %v", size/ratio, "KB")
	case size/ratio/ratio < ratio:
		return fmt.Sprintf("%v %v", size/ratio/ratio, "MB")
	case size/ratio/ratio/ratio < ratio:
		return fmt.Sprintf("%v %v", size/ratio/ratio/ratio, "GB")
	default:
		return ""
	}
}
