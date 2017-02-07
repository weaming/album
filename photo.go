package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"os"
	fp "path/filepath"
)
import (
	"github.com/nfnt/resize"
)

func thumbnail(path, outpath string) error {
	// open "test.jpg"
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	ext := fp.Ext(path)
	var img image.Image
	switch ext {
	case ".jpg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".gif":
		img, err = gif.Decode(file)
	default:
		img, err = jpeg.Decode(file)
	}
	if err != nil {
		return err
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	var max_size uint = 500
	m := resize.Thumbnail(max_size, max_size, img, resize.Bilinear)

	out, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer out.Close()

	// write new image to file
	switch ext {
	case ".jpg":
		jpeg.Encode(out, m, nil)
	case ".png":
		png.Encode(out, m)
	case ".gif":
		gif.Encode(out, m, nil)
	default:
		jpeg.Encode(out, m, nil)
	}
	return nil
}

func thumb_directory(tododir, outdir string) error {
	TODODIR := NewDir(tododir)
	err := os.MkdirAll(outdir, 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	for _, file := range TODODIR.AbsImages {
		rel_path, err := fp.Rel(tododir, file)
		if err != nil {
			return err
		}

		out_path := fp.Join(outdir, rel_path)
		if _, err := os.Stat(out_path); err == nil {
			fmt.Printf("Ignore existed thumbnail: %v\n", out_path)
			continue
		}

		fmt.Printf("Thumbnailing: %v\n", rel_path)
		thumbnail(file, out_path)
	}
	for index, dir := range TODODIR.AbsDirs {
		thumb_directory(dir, fp.Join(outdir, TODODIR.Dirs[index]))
	}
	return nil
}

func UrlEncoded(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return ""
	}
	return u.String()
}
