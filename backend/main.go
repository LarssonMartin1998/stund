package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"backend/blog"
	"backend/config"
	"backend/contact"
	"backend/database"
	"backend/router"
	"backend/utils"
)

func main() {
	cfg := utils.Must(config.Load())
	db := utils.Must(database.NewSQLiteDB(cfg))
	r := router.New(cfg)

	r.SetupRoutes(&router.RoutingContext{
		Providers: []router.RouteProvider{
		},
	}, cfg)

	port := ":" + cfg.Port
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	fmt.Fprintf(writer, "✅ Starting HTTP server on port %s\n", port)
	r.WriteRoutes(writer)

	r.ClearRouteTree()

	writer.Flush()

	server := &http.Server{
		Addr:           port,
		Handler:        r.GetHTTPHandler(),
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		log.Println(buf.String())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), max(cfg.Server.ReadTimeout, cfg.Server.WriteTimeout))
	defer cancel()

	server.Shutdown(ctx)
	db.Close()
}
