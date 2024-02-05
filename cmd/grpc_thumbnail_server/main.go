package main

import (
	"github.com/OksidGen/grpc_thumbnail/server/internal/app"
	"github.com/OksidGen/grpc_thumbnail/server/pkg"
)

func main() {
	pkg.SetupLogger()
	app.Run()
}
