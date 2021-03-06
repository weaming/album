package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	fp "path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
)

const DEFAULT_PW = "admin"

var size int
var NEED_AUTH bool
var ADMIN, PASSWORD string

func main() {
	var LISTEN = flag.String("l", ":8000", "Listen [host]:port, default bind to 0.0.0.0")
	var OUTDIR = flag.String("o", "", "The directory of thumnail. Default [$ROOT/../thumbnail]")
	var MAX_WIDTH = flag.Uint("wd", 200, "The maximum width of output photo.")
	var MAX_HEIGHT = flag.Uint("ht", 200, "The maximum height of output photo.")

	flag.BoolVar(&NEED_AUTH, "a", true, "Whether need authorization.")
	flag.StringVar(&ADMIN, "u", "admin", "Basic authentication username")
	flag.StringVar(&PASSWORD, "p", DEFAULT_PW, "Basic authentication password")
	flag.IntVar(&size, "n", 20, "The maximum number of photos in each page.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] ROOT\nThe ROOT is the directory contains photos.\n\n", os.Args[0])
		flag.PrintDefaults()
	}
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

	// thumbnail cache directory
	var outdir = ""
	if *OUTDIR != "" {
		outdir = *OUTDIR
	} else {
		outdir = fp.Join(fp.Dir(ROOT), "thumbnail")
	}
	green(fmt.Sprintf("To be listed direcotry: [%v]", ROOT))
	green(fmt.Sprintf("Cache direcotry: [%v]", outdir))

	// basic authentication
	if NEED_AUTH {
		if PASSWORD == DEFAULT_PW {
			red("Warning: set yourself password by option -p")
		}
		green(fmt.Sprintf("Your basic auth name and password: [%v:%v]", ADMIN, PASSWORD))
	} else {
		red("Warning: please set your HTTP basic authentication")
	}

	go func() {
		for {
			thumb_directory(ROOT, outdir, *MAX_WIDTH, *MAX_HEIGHT)
			time.Sleep(time.Minute * 20)
		}
	}()

	Redirect("/", "/index")
	ServeFile("/favicon.ico", fp.Join(ROOT, "./favicon.ico"))

	http.Handle("/index/", gziphandler.GzipHandler(MyAlbum{root: ROOT}))
	ServeDir("/img/", ROOT)
	ServeDir("/thumb/", outdir)

	fmt.Printf("Open http://127.0.0.1:%v to enjoy!\n", strings.Split(*LISTEN, ":")[1])
	err = http.ListenAndServe(*LISTEN, nil)
	if err != nil {
		fmt.Println(err)
	}
}

type MyAlbum struct {
	root string
	dir  *Dir
}

func (album MyAlbum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logit(r)
	if NEED_AUTH {
		mybasicAuth(album.handlerFunc, ADMIN, PASSWORD)(w, r)
	} else {
		album.handlerFunc(w, r)
	}
}

func (album MyAlbum) handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Happy coding!")
	pathName := r.URL.Path

	page, err := getQueryInt(r, "page")
	if err != nil {
		target, _ := AddQuery(pathName, "page", "1")
		http.Redirect(w, r, target, http.StatusFound)
		return
	}

	obj := NewDir(fp.Join(album.root, pathName[6:]))
	if obj == nil {
		w.Write([]byte("Invalid URL"))
		return
	} else {
		album.dir = obj
	}

	pagination, htmlImages, returnPage := Img2Html(pathName, album.dir, page)
	if returnPage != page {
		target, _ := AddQuery(pathName, "page", strconv.Itoa(returnPage))
		http.Redirect(w, r, target, http.StatusFound)
	}

	w.Write([]byte(fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<title>My Photos</title>
		<style>
		.right{float: right;}
		.card{
		background-color: #fff;
		box-shadow: 0 2px 5px 0 rgba(0, 0, 0, .16), 0 2px 10px 0 rgba(0, 0, 0, .12);
		margin: 0 auto 1rem auto;
		padding: 1rem;
		max-width: 900px;
		border-radius: 3px;
		}
		.directory:hover
		{
		background-color: #eee;
		}

		div.pagination {
		min-height: 20px;
		}

		#pagination {
		display: grid;
		grid-template-columns: 1fr 150px;
		}

		div.pagination a{
		display: inline-block;
		border: 1px solid #aaa;
		padding: 5px 10px;
		margin: 5px 10px;
		border-radius: 4px;
		color: black;
		text-decoration: none;
		}
		div.pagination a:hover{
		box-shadow: 0 2px 5px 0 rgba(0, 0, 0, .16), 0 2px 10px 0 rgba(0, 0, 0, .12);
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
		}
		a.photo img.thumbnail{
		width: 100%%;
		height: 100%%;
		border: 1px solid #ccc;
		}
		a.photo img.thumbnail:hover{
		opacity: 0.7;
		border: 1px solid chocolate;
		box-shadow: 0 2px 5px 0 rgba(0, 0, 0, .16), 0 2px 10px 0 rgba(0, 0, 0, .12);
		}
		</style>
		</head>
		<body>
		<div class="card directories">
		<h3> Directories: %v <a href="/index" class="right">Home</a> </h3>
		<div>%v</div>
		</div>
		<div class="card photos">
		<h3 id="pagination">Photos: %v Size: %v
		<div class="pagiContainer">%v</div>
		</h3>
		<div class="container"> %v </div>
		</div>
		</body>
		</html>`,
		len(album.dir.Dirs),
		strings.Join(Dir2Html(pathName, album.dir), "\n"),
		len(album.dir.Images),
		some_files_size_str(album.dir.AbsImages),
		pagination,
		strings.Join(htmlImages, "\n"),
	)))
}

func Img2Html(pathName string, dir *Dir, page int) (string, []string, int) {
	var (
		pagination string
		htmlImages []string
	)

	_images, previous, next, page := Page(dir.Images, page, size)
	_abs_images, previous, next, page := Page(dir.AbsImages, page, size)

	// add pagination
	var htmlPrevious, htmlNext string
	if previous {
		newUrl, _ := AddQuery(pathName, "page", strconv.Itoa(page-1))
		htmlPrevious = fmt.Sprintf(`<a class="previous" href="%v">←</a>`, newUrl)
	}
	if next {
		newUrl, _ := AddQuery(pathName, "page", strconv.Itoa(page+1))
		htmlNext = fmt.Sprintf(`<a class="next" href="%v">→</a>`, newUrl)
	}
	if previous || next {
		pagination = fmt.Sprintf(`<div class="pagination">%v%v</div>`, htmlPrevious, htmlNext)
		//pagination = htmlPrevious + htmlNext
	}

	for index, file := range _images {
		u, _ := url.Parse(pathName[6:])
		u.Path = path.Join("/thumb/", u.Path, file)

		htmlImages = append(htmlImages, fmt.Sprintf(`<a class="photo" href="%v"><img src="%v" class="thumbnail" title="%v"></a>`,
			"/img/"+path.Join(pathName[6:], file),
			UrlEncoded(u.String()),
			fmt.Sprintf("%v [%v]", file, file_size_str(_abs_images[index]))))
	}
	return pagination, htmlImages, page
}

func Dir2Html(pathName string, dir *Dir) []string {
	rv := []string{}
	for index, file := range dir.Dirs {
		if hasPhoto(dir.AbsDirs[index]) {
			sub_dir := NewDir(dir.AbsDirs[index])
			rv = append(rv, fmt.Sprintf(
				`<div class="directory"><a class="link" href="%v">%v</a><span class="count right">[%v]</span><span class="right">%v</span></div>`,
				"/index/"+fp.Join(pathName[6:], file)+"/",
				file+"/",
				len(sub_dir.Images),
				dir_images_size_str(dir.AbsDirs[index])))
		}
	}
	return rv
}

func Page(items []string, page, size int) ([]string, bool, bool, int) {
	if len(items) == 0 {
		return []string{}, false, false, 1
	}

	end := size * page
	start := end - size
	next := end < len(items)

	if len(items) <= start {
		_page := len(items) / size
		if _page*size < len(items) {
			_page++
		}
		return Page(items, _page, size)
	}
	if !next {
		end = len(items)
	}
	return items[start:end], page > 1, next, page
}
