package main

import (
	"Fuse/fuse"
	"Fuse/middleware"
)

func main() {
	app := fuse.New()

	app.Use(middleware.Defaults()...)

	httpSrv := app.HTTP()
	httpSrv.Get("/ping", func(c fuse.Context) fuse.Result {
		return c.Success(fuse.H{"message": "pong"})
	})

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
