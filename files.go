package main

import (
	"fmt"
	"io/ioutil"
	"os"
	fp "path/filepath"
	"strings"
)

type Dir struct {
	Root      string
	Dirs      []string
	Files     []string
	Images    []string
	AbsDirs   []string
	AbsFiles  []string
	AbsImages []string
}

func NewDir(path string) *Dir {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !fi.IsDir() {
		return nil
	} else {
		dir := Dir{Root: path}
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
				case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
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
	size_float := float64(size)
	const ratio = 1024

	switch {
	case size < ratio:
		return fmt.Sprintf("%.2f %v", size_float, "B")
	case size/ratio < ratio:
		return fmt.Sprintf("%.2f %v", size_float/ratio, "KB")
	case size/ratio/ratio < ratio:
		return fmt.Sprintf("%.2f %v", size_float/ratio/ratio, "MB")
	case size/ratio/ratio/ratio < ratio:
		return fmt.Sprintf("%.2f %v", size_float/ratio/ratio/ratio, "GB")
	default:
		return ""
	}
}

func get_size(path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}

func some_files_size_int64(files []string) (total int64) {
	for _, path := range files {
		total += get_size(path)
	}
	return
}

func some_sub_dir_images_size_int64(dirs []string) (total int64) {
	for _, path := range dirs {
		tmp := NewDir(path)
		total = total + some_files_size_int64(tmp.AbsImages) + some_sub_dir_images_size_int64(tmp.AbsDirs)
	}
	return
}

func file_size_str(path string) string {
	return size2text(get_size(path))
}

func some_files_size_str(files []string) string {
	var total int64
	for _, file := range files {
		total += get_size(file)
	}
	return size2text(total)
}

func dir_images_size_str(dir string) string {
	tmp := NewDir(dir)
	return size2text(some_files_size_int64(tmp.AbsImages) + some_sub_dir_images_size_int64(tmp.AbsDirs))
}

func hasPhoto(path string) bool {
	dir := NewDir(path)
	if len(dir.Images) > 0 {
		return true
	} else {
		for _, subpath := range dir.AbsDirs {
			if hasPhoto(subpath) {
				return true
			}
		}
	}
	return false
}
