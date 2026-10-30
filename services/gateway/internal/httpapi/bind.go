package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func readJSON(r *http.Request, dst any) {
	_ = r.ParseForm()
	if strings.Contains(r.Header.Get("Content-Type"), "json") && r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(dst)
	}
}

func formOr(r *http.Request, jsonValue, key string) string {
	if strings.TrimSpace(jsonValue) != "" {
		return strings.TrimSpace(jsonValue)
	}
	return strings.TrimSpace(r.FormValue(key))
}

func queryInt(r *http.Request, key string) int64 {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		raw = r.FormValue(key)
	}
	n, _ := strconv.ParseInt(raw, 10, 64)
	return n
}

func pageLinks(path string, page, total, size int) map[string]string {
	if size <= 0 {
		size = 50
	}
	last := (total + size - 1) / size
	if last < 1 {
		last = 1
	}
	link := func(p int) string {
		if p < 1 || p > last {
			return ""
		}
		return path + "?page=" + strconv.Itoa(p)
	}
	prev := ""
	if page > 1 {
		prev = link(page - 1)
	}
	next := ""
	if page < last {
		next = link(page + 1)
	}
	return map[string]string{
		"self":  link(page),
		"first": link(1),
		"prev":  prev,
		"next":  next,
		"last":  link(last),
	}
}
