package main

import (
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"strings"
)

type DirStr struct {
	Root      string
	Dirs      []string
	Files     []string
	Images    []string
	AbsDirs   []string
	AbsFiles  []string
	AbsImages []string
}

func NewDirstr(path string) *DirStr {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !fi.IsDir() {
		return nil
	} else {
		dir := DirStr{Root: path}
		files, _ := ioutil.ReadDir(path)
		for _, fi := range files {
			absPath := fp.Join(path, fi.Name())
			relPath := fi.Name()
			if fi.IsDir() {
				dir.Dirs = append(dir.Dirs, relPath)
				dir.AbsDirs = append(dir.AbsDirs, absPath)
			} else {
				dir.Files = append(dir.Files, relPath)
				dir.AbsFiles = append(dir.AbsFiles, absPath)
				switch strings.ToLower(fp.Ext(relPath)) {
				case ".jpg", ".png", ".gif":
					dir.Images = append(dir.Images, relPath)
					dir.AbsImages = append(dir.AbsImages, absPath)
				default:
				}
			}
		}
		return &dir
	}
}

func size2text(size int64) string {
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

func getSize(path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}

func fileSize(path string) string {
	return size2text(getSize(path))
}

func allFilesSize(files []string) string {
	var total int64
	for _, path := range files {
		total += getSize(path)
	}
	return size2text(total)
}

func dirSize(dir string) string {
	return allFilesSize(NewDirstr(dir).AbsImages)
}
