package main

import (
	"Fuse/core"
	"Fuse/fuse"
	"Fuse/middleware"
)

func main() {
	app := fuse.New()

	app.Use(middleware.Defaults()...)

	httpSrv := app.HTTP()
	httpSrv.Get("/ping", func(c core.Ctx) core.Result {
		return core.OK(core.H{"message": "pong"})
	})

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
