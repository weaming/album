package main

import (
	"fmt"
	"net/url"
)

func a(url, text string) string {
	return fmt.Sprintf(`<a href="%v">%v</a>`, url, text)
}

func img(url, title string) string {
	return fmt.Sprintf(`<img src="%v" title="%v">`, url, title)
}

func ele(T, class, text string) string {
	return fmt.Sprintf(`<%v class="%v">%v</%v>`, T, class, text, T)
}

func span(text, class string) string {
	return ele("span", class, text)
}

func div(text, class string) string {
	return ele("div", class, text)
}

func p(text, class string) string {
	return ele("p", class, text)
}

func UrlEncoded(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return ""
	}
	return u.String()
}

func AddQuery(pathName, key, value string) (string, error) {
	u, err := url.Parse(pathName)
	if err != nil {
		return pathName, err
	}
	q := u.Query()
	q.Set("page", value)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
