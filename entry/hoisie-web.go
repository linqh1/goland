package main

import (
	"github.com/hoisie/web"
)

func hello(val string) string {
	return "hello " + val
}

func main() {
	web.Get("/index(.*)", func(val string) string {
		return "other " + val
	})
	web.Get("/(.*)", func(val string) string {
		return "hello " + val
	})
	web.Run("0.0.0.0:9999")
}
