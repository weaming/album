package main

import (
	"flag"
	"fmt"
	"os"
	fp "path/filepath"
	"strings"
	"time"
)

import "github.com/kataras/iris"
import "github.com/iris-contrib/middleware/basicauth"

const DEFAULT_PW = "admin"

func main() {
	var LISTEN = flag.String("l", ":8000", "Listen [host]:port, default bind to 0.0.0.0")
	var ADMIN = flag.String("u", "admin", "Basic authentication username")
	var PASSWORD = flag.String("p", DEFAULT_PW, "Basic authentication password")
	var OUTDIR = flag.String("o", "", "The directory of thumnail. Default [$ROOT/../thumbnail]")
	var outdir = ""
	flag.Parse()

	// check the directory path
	ROOT, _ := fp.Abs(flag.Arg(0))
	fi, err := os.Stat(ROOT)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !fi.IsDir() {
		fmt.Fprintln(os.Stderr, "The path should be a directory!!")
		os.Exit(1)
	}

	if *OUTDIR != "" {
		outdir = *OUTDIR
	} else {
		outdir = fp.Join(fp.Dir(ROOT), "thumbnail")
	}
	fmt.Printf("To be listed direcotry: [%v]\n", ROOT)
	thumb_directory(ROOT, outdir)

	fmt.Printf("Your basic authentication username: [%v]\n", *ADMIN)
	fmt.Printf("Your basic authentication password: [%v]\n", *PASSWORD)
	if *PASSWORD == DEFAULT_PW {
		fmt.Println("Warning: set yourself password")
	}

	authConfig := basicauth.Config{
		Users:   map[string]string{*ADMIN: *PASSWORD},
		Expires: time.Duration(5) * time.Minute,
	}
	auth := basicauth.New(authConfig)

	iris.Config.Gzip = true // compressed gzip contents to the client, the same for Serializers also, defaults to false

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Redirect("/index")
	})
	iris.StaticWeb("/img", ROOT)
	iris.StaticWeb("/thumb", outdir)

	needAuth := iris.Party("/index", auth)
	{
		needAuth.Handle("GET", "/*path", MyAlbum{root: ROOT})
	}
	iris.Listen(*LISTEN)
}

type MyAlbum struct {
	root string
	dir  *Dir
}

func (album MyAlbum) Serve(ctx *iris.Context) {
	path := ctx.Path()
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
				.right{float: right;}
				.region{
					background-color: #fff;
					box-shadow: 0 2px 5px 0 rgba(0, 0, 0, .16), 0 2px 10px 0 rgba(0, 0, 0, .12);
					margin: 0 auto 1rem auto;
					padding: 1rem;
					max-width: 900px;
				}
				.photo:hover,
				.directory:hover
				{
					background-color: #eee;
				}

				div.photos div.container{
					display: flex;
					justify-content: space-around;
					flex-wrap: wrap;
					box-sizing: border-box;
				}

				a.photo{
					max-width: 200px;
					max-height: 200px;
					margin: 5px;
					border: 1px solid #ccc;
				}
				a.photo:hover{
					border: 1px solid chocolate;
				}
				a.photo img.thumbnail{
					display: inline-block;
					width: 100%%;
					height: 100%%;
				}
				a.photo img.thumbnail:hover{
					opacity: 0.7;
				}
			</style>
		</head>
		<body>
			<div class="region directories">
				<h3> Directories: %v <a href="/index" class="right">Home</a> </h3>
				%v
			</div>
			<div class="region photos">
				<h3>Photos: %v Size: %v</h3>
				<div class="container"> %v </div>
			</div>
		</body>
		</html>`,
		len(album.dir.Dirs),
		strings.Join(Dir2Html(path, album.dir), "\n"),
		len(album.dir.Images),
		some_files_size_str(album.dir.AbsImages),
		strings.Join(Img2Html(path, album.dir), "\n")))
}

func Img2Html(path string, dir *Dir) []string {
	rv := []string{}
	for index, file := range dir.Images {
		rv = append(rv, fmt.Sprintf(`<a class="photo" href="%v"><img src="%v" class="thumbnail" title="%v"></a>`,
			"/img/"+fp.Join(path[7:], file),
			UrlEncoded("/thumb/"+fp.Join(path[7:], file)),
			fmt.Sprintf("%v [%v]", file, file_size_str(dir.AbsImages[index]))))
	}
	return rv
}

func Dir2Html(path string, dir *Dir) []string {
	rv := []string{}
	for index, file := range dir.Dirs {
		if hasPhoto(dir.AbsDirs[index]) {
			rv = append(rv, h_div(
				h_span(h_a("/index/"+fp.Join(path[7:], file), file+"/"), "link")+h_span(dir_images_size_str(dir.AbsDirs[index]), "right"), "directory"))
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
