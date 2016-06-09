package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Web struct {
	mux *Mux
}

func NewWeb() *Web {
	return &Web{mux: NewMux()}
}

func (w *Web) ListenAndServe(s *Store) {

	r := gin.Default()

	r.Static("/assets", "./assets")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/assets/index.html")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		w.mux.Handle(c.Writer, c.Request)
	})

	r.GET("/highscores", func(c *gin.Context) {
		c.JSON(200, s.GetHighscores())
	})

	r.Run()
}
