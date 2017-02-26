package main

import (
	"log"
	"net/http"
	"strconv"
)

func redirect_handler(to string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logit(r)
		http.Redirect(w, r, to, 301)
	}
}

func Redirect(from, to string) {
	http.Handle(from, redirect_handler(to))
}

func ServeDir(prefix, path string) {
	http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(path))))
}

func logit(r *http.Request) {
	log.Printf(`%v "%v %v %v"`, r.RemoteAddr, r.Method, r.RequestURI, r.Proto)
}

func getQuery(r *http.Request, name string) (v string) {
	v = r.URL.Query().Get("page")
	return
}

func getQueryInt(r *http.Request, name string) (v int, err error) {
	v_str := getQuery(r, name)
	v, err = strconv.Atoi(v_str)
	return
}
