package main

import (
	"k8sman/internal/server"
)

func main() {
	ctx := server.NewRuntimeContext()
	ctx.Run()
}
