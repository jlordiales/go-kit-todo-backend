package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/handlers"
	"github.com/jlordiales/go-kit-todo-backend/pkg/endpoints"
	todohttp "github.com/jlordiales/go-kit-todo-backend/pkg/http"
	"github.com/jlordiales/go-kit-todo-backend/pkg/service"
	"github.com/jlordiales/go-kit-todo-backend/pkg/service/inmemory"
)

func main() {
	addr := flag.String("PORT", ":8080", "HTTP listen address")
	basePath := flag.String("base_path", "http://localhost", "Base url for requests")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stdout)
		logger = log.With(logger, "timestamp", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var s service.Service
	{
		s = inmemory.NewService(logger)
	}

	endpoints := endpoints.New(s, fmt.Sprintf("%s%s", *basePath, *addr), logger)

	var handler http.Handler
	{
		handler = todohttp.MakeHandler(log.With(logger, "handler", "HTTP"), endpoints)
		handler = handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PATCH", "OPTIONS", "HEAD"}), handlers.AllowedOrigins([]string{"*"}))(handler)
	}

	srv := &http.Server{
		Handler:      handler,
		Addr:         *addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", *addr)
		errs <- srv.ListenAndServe()
	}()

	logger.Log("exit", <-errs)
}
