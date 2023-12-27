package main

import (
	_ "port-forwording/internal/network"
	"port-forwording/web"
)

func main() {
	web.Run()
}
