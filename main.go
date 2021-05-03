package main

import (
	"flag"

	"github.com/kiselev-nikolay/inflr-be/pkg/server"
)

func main() {
	flagDev := flag.Bool("dev", false, "Enable logs and development helpers")
	flag.Parse()
	mode := server.ModeProduction
	if *flagDev {
		mode = server.ModeDev
	}
	app := server.GetRouter(mode)
	app.Run(":8080")
}
