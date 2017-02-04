package main

import (
	"flag"
	"fmt"
	fp "path/filepath"
	"strings"
)

import "github.com/kataras/iris"

var ROOT string
var LISTEN = ":8000"

func main() {
	flag.Parse()
	ROOT = flag.Arg(0)
	if flag.Arg(1) != "" {
		LISTEN = flag.Arg(1)
	}
	fmt.Printf("To be listed direcotry: [%v]\n", ROOT)

	iris.Config.IsDevelopment = true // reloads the templates on each request, defaults to false
	iris.Config.Gzip = true          // compressed gzip contents to the client, the same for Serializers also, defaults to false

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Writef(h_a("/public", "View your photos!"))
	})
	iris.StaticWeb("/img", ROOT)

	iris.Handle("GET", "/public/*path", MyAlbum{root: ROOT})
	iris.Listen(LISTEN)
}

type MyAlbum struct {
	root string
	dir  *DirStr
}

func (album MyAlbum) Serve(ctx *iris.Context) {
	path := ctx.Path()
	ext := strings.ToLower(fp.Ext(path))

	switch ext {
	case ".jpg", ".png", ".gif":
		ctx.WriteString("ok")
		//ctx.ServeFile(fp.Join(album.root, ctx.Param("path")))
	default:
		obj := NewDirstr(fp.Join(album.root, ctx.Param("path")))
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
					.img:hover{background-color: #eee;}
				</style>
			</head>
			<body>
				<div class="region">
					<h3>Directories: %v</h3>
					%v
				</div>
				<div class="region">
					<h3>Photos: %v</h3>
					%v
				</div>
			</body>
			</html>`,
			len(album.dir.Dirs),
			strings.Join(Dir2Html(album.dir.Root, album.dir.Dirs), "<br>"),
			len(album.dir.Images),
			strings.Join(Img2Html(path, album.dir.Root, album.dir.Images), "")))
	}
}

func Img2Html(path, root string, names []string) []string {
	rv := []string{}
	for _, file := range names {
		rv = append(rv, h_div(
			h_span(h_a("/img/"+fp.Join(path[8:], file), file), "link")+h_span(fileSize(fp.Join(root, file)), "size"), "img"))
	}
	return rv
}

func Dir2Html(root string, names []string) []string {
	rv := []string{}
	for _, file := range names {
		if len(NewDirstr(fp.Join(root, file)).Images) > 0 {
			rv = append(rv, h_a("/public/"+file, file+"/"))
		}
	}
	return rv
}
