package main

import (
	"Fuse/core"
	"Fuse/fuse"
	"Fuse/middleware"
	"Fuse/ssex"
	"Fuse/wsx"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

	httpSrv.GET("/sse", ssex.Upgrade(func(c core.Ctx, stream *ssex.Stream) error {
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

	httpSrv.GET("/ws", wsx.Upgrade(func(c core.Ctx, conn *websocket.Conn) error {
		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					return
				}
			}
		}()
		words := []string{"大家好", "我是你爹"}

		for _, word := range words {
			err := conn.WriteMessage(websocket.TextMessage, []byte(word))
			if err != nil {
				return err
			}
		}
		time.Sleep(1 * time.Second)

		return nil
	}))

	if err := app.Run(); err != nil {
		panic(err)
	}
}
