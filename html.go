package main

import (
	"fmt"
	"net/url"
)

func h_a(url, text string) string {
	return fmt.Sprintf(`<a href="%v">%v</a>`, url, text)
}

func h_img(url, title string) string {
	return fmt.Sprintf(`<img src="%v" title="%v">`, url, title)
}

func h_ele(T, class, text string) string {
	return fmt.Sprintf(`<%v class="%v">%v</%v>`, T, class, text, T)
}

func h_span(text, class string) string {
	return h_ele("span", class, text)
}

func h_div(text, class string) string {
	return h_ele("div", class, text)
}

func h_p(text, class string) string {
	return h_ele("p", class, text)
}

func UrlEncoded(str string) string {
	u, err := url.Parse(str)
	if err != nil {
		return ""
	}
	return u.String()
}
