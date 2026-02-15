// @title           My Service API
// @version         1.0.0
// @description     Go API Scaffold - RESTful API documentation
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Enter your Bearer token
package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-api-scaffold/internal/handler"
	"go-api-scaffold/internal/service"
	"go-api-scaffold/internal/store"
	"go-api-scaffold/pkg/config"
	"go-api-scaffold/pkg/logger"

	"google.golang.org/grpc"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Command line flags
	configPath := flag.String("c", "configs/config.yaml", "config file path")
	showVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s %s\n  Build: %s\n  Commit: %s\n", "my-service", Version, BuildTime, GitCommit)
		os.Exit(0)
	}

	// ====== 1. Load config ======
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// ====== 2. Init logger ======
	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Infof("starting %s %s (build: %s, commit: %s)", cfg.App.Name, Version, BuildTime, GitCommit)

	// ====== 3. Init database ======
	db, err := store.New(&cfg.Database)
	if err != nil {
		logger.Fatalf("failed to init database: %v", err)
	}
	defer db.Close()

	// ====== 4. Init service layer ======
	authSvc := service.NewAuthService(db, cfg.JWT.Secret, cfg.JWT.Expire, cfg.JWT.RefreshHours)
	exampleSvc := service.NewExampleService(db)

	// ====== 5. Start HTTP server ======
	r := handler.NewRouter(cfg, authSvc, exampleSvc)
	httpAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := &http.Server{
		Addr:         httpAddr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Infof("HTTP server started: http://%s", httpAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// ====== 6. Start gRPC server (optional) ======
	var grpcServer *grpc.Server
	if cfg.GRPC.Enabled {
		grpcServer = handler.NewGRPCServer(exampleSvc)
		grpcAddr := fmt.Sprintf(":%d", cfg.GRPC.Port)
		lis, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			logger.Fatalf("gRPC listen failed: %v", err)
		}
		go func() {
			logger.Infof("gRPC server started: %s", grpcAddr)
			if err := grpcServer.Serve(lis); err != nil {
				logger.Fatalf("gRPC server error: %v", err)
			}
		}()
	}

	// ====== 7. Graceful shutdown ======
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Infof("received signal: %v, shutting down...", sig)

	// Stop gRPC
	if grpcServer != nil {
		grpcServer.GracefulStop()
		logger.Info("gRPC server stopped")
	}

	// Stop HTTP (5s timeout)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Errorf("HTTP server shutdown error: %v", err)
	} else {
		logger.Info("HTTP server stopped")
	}

	logger.Info("service exited")
}
