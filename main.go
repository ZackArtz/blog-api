package main

import "github.com/zackartz/blog-api/fiber/api"

func main() {
	server := api.Server{}

	server.Initialize()
}
