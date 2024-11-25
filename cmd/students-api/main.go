package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sayansaha934/students-api/internal/config"
	"github.com/sayansaha934/students-api/internal/http/handlers/student"
	"github.com/sayansaha934/students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg:=config.MustLoad()
	// Database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Database setup successfully", slog.String("env", cfg.Env))
	// setup router
	// router:=http.NewServeMux()
	router:=mux.NewRouter()
	router.HandleFunc("/api/students", student.New(storage)).Methods("POST")
	router.HandleFunc("/api/students/{id}", student.GetById(storage)).Methods("GET")
	router.HandleFunc("/api/students", student.GetList(storage)).Methods("GET")
	router.HandleFunc("/api/students/{id}", student.Delete(storage)).Methods("DELETE")
	router.HandleFunc("/api/students/{id}", student.Update(storage)).Methods("PUT")

	// setup server
	server:=http.Server{
		Addr:cfg.HTTPServer.Addr,
		Handler:router,
	}
	slog.Info("Server is started", slog.String("address", cfg.HTTPServer.Addr))
	done:=make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	go func(){
		err:=server.ListenAndServe()
		if err!=nil {
			log.Fatal("Failed to start server")
	}
	}()
	<-done
	slog.Info("Shutting down server")

	ctx,cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err:=server.Shutdown(ctx); err!=nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("Server shutdown successfully")
	
}
