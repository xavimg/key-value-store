package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct {
	Storage    store.Storer[string, string]
	listenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		Storage:    store.NewKVStore[string, string](),
		listenAddr: listenAddr,
	}
}

func (s *Server) handlePush(c echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	s.Storage.Push(key, value)

	return c.JSON(http.StatusOK, map[string]string{"msg ": "ok"})
}

func (s *Server) handleGet(c echo.Context) error {
	key := c.Param("key")

	value, err := s.Storage.Get(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"value": value})
}

func (s *Server) Start() {
	fmt.Println("HTTP server is running on", s.listenAddr)

	e := echo.New()

	e.GET("push/:key/:value", s.handlePush)
	e.GET("get/:key", s.handleGet)

	e.Start(s.listenAddr)
}

func main() {
	s := NewServer(":3000")

	s.Start()
}
