package main

import (
	"Fuse/core"
	"Fuse/fuse"
	"Fuse/middleware"
	"Fuse/ssex"
	"log"
	"net/http"
	"time"
)

func main() {
	app := fuse.New()

	app.Use(middleware.Defaults()...)

	httpSrv := app.HTTP()
	httpSrv.GET("/ping/:id", func(c fuse.Context) fuse.Result {
		id := c.Param("id")
		log.Printf("id: %s", id)
		return c.Success(fuse.H{"message": "pong"}).WithHttpStatus(http.StatusInternalServerError)
	})

	httpSrv.GET("/ping/id/abc", func(c fuse.Context) fuse.Result {
		return c.Success(fuse.H{"message": "pong"}).WithHttpStatus(http.StatusOK)
	})

	httpSrv.GET("/sse", ssex.Upgrader(func(c core.Ctx, stream *ssex.Stream) error {
		words := []string{"王", "昊", "声", "是", "我", "儿"}
		for _, word := range words {
			err := stream.Send("message", word)
			if err != nil {
				return err
			}
			time.Sleep(1 * time.Second)
		}
		return stream.Send("done", nil)
	}))

	if err := app.Run(); err != nil {
		panic(err)
	}
}
