package main

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	fp "path/filepath"
)
import (
	"github.com/nfnt/resize"
)

func thumbnail(path, outpath, ext string) error {
	// open "test.jpg"
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
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
	case ".jpg", ".jpeg":
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
	_, err := os.Stat(outdir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(outdir, 0755)
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
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
		err = thumbnail(file, out_path, fp.Ext(file))
		if err != nil {
			err := thumbnail(file, out_path, fp.Ext(".jpg"))
			if err != nil {
				return err
			}
		}
	}
	for index, dir := range TODODIR.AbsDirs {
		thumb_directory(dir, fp.Join(outdir, TODODIR.Dirs[index]))
	}
	return nil
}
