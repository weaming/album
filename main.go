package main

import (
	"flag"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"
)

import "github.com/kataras/iris"

var ROOT string
var LISTEN = ":8000"

func main() {
	flag.Parse()
	ROOT, _ := fp.Abs(flag.Arg(0))
	fi, err := os.Stat(ROOT)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(11)
	}
	if !fi.IsDir() {
		fmt.Fprintln(os.Stderr, "The path should be a directory!!")
		os.Exit(12)
	}

	if flag.Arg(1) != "" {
		LISTEN = flag.Arg(1)
	}
	fmt.Printf("To be listed direcotry: [%v]\n", ROOT)

	iris.Config.IsDevelopment = true // reloads the templates on each request, defaults to false
	iris.Config.Gzip = true          // compressed gzip contents to the client, the same for Serializers also, defaults to false

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("/index")
	})
	iris.StaticWeb("/img", ROOT)

	iris.Handle("GET", "/index/*path", MyAlbum{root: ROOT})
	iris.Listen(LISTEN)
}

type MyAlbum struct {
	root string
	dir  *Dir
}

func (album MyAlbum) Serve(ctx *iris.Context) {
	path := ctx.Path()
	ext := strings.ToLower(fp.Ext(path))

	switch ext {
	case ".jpg", ".png", ".gif":
		ctx.WriteString("ok")
		//ctx.ServeFile(fp.Join(album.root, ctx.Param("path")))
	default:
		obj := NewDir(fp.Join(album.root, ctx.Param("path")))
		if obj == nil {
			ctx.WriteString("Invalid URL")
			return
		} else {
			album.dir = obj
		}
		ctx.WriteString(fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>My Photos</title>
				<style>
					.size{float: right;}
					.region{
					background-color: #fff;
					box-shadow: 0 2px 5px 0 rgba(0, 0, 0, .16), 0 2px 10px 0 rgba(0, 0, 0, .12);
					margin: 0 auto 1rem auto;
					padding: 1rem;
					max-width: 900px;
					}
					.img:hover,
					.directory:hover
					{background-color: #eee;}
				</style>
			</head>
			<body>
				<div class="region">
					<h3> Directories: %v <a href="/index" style="float: right;">Home</a> </h3>
					%v
				</div>
				<div class="region">
					<h3>Photos: %v Size: %v</h3>
					%v
				</div>
			</body>
			</html>`,
			len(album.dir.Dirs),
			strings.Join(Dir2Html(path, album.dir), ""),
			len(album.dir.Images),
			some_files_size_str(album.dir.AbsImages),
			strings.Join(Img2Html(path, album.dir), "")))
	}
}

func Img2Html(path string, dir *Dir) []string {
	rv := []string{}
	for index, file := range dir.Images {
		rv = append(rv, h_div(
			h_span(h_a("/img/"+fp.Join(path[7:], file), file), "link")+h_span(file_size_str(dir.AbsImages[index]), "size"), "img"))
	}
	return rv
}

func Dir2Html(path string, dir *Dir) []string {
	rv := []string{}
	for index, file := range dir.Dirs {
		if hasPhoto(dir.AbsDirs[index]) {
			rv = append(rv, h_div(
				h_span(h_a("/index/"+fp.Join(path[7:], file), file+"/"), "link")+h_span(dir_images_size_str(dir.AbsDirs[index]), "size"), "directory"))
		}
	}
	return rv
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
