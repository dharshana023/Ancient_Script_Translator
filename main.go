package main

import (
        "context"
        "flag"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "ancient-script-decoder/api"
        "ancient-script-decoder/services"
        "ancient-script-decoder/utils"
)

func main() {
        // Parse command line flags
        configPath := flag.String("config", "config.yaml", "Path to configuration file")
        flag.Parse()

        // Initialize logger
        logger := utils.NewLogger()
        logger.Info("Starting Ancient Script Decoder service")

        // Load configuration
        config, err := utils.LoadConfig(*configPath)
        if err != nil {
                logger.Fatal("Failed to load configuration", "error", err)
        }

        // Create cancellation context for cleanup
        _, cancel := context.WithCancel(context.Background())
        defer cancel()

        // Initialize services
        imageProcessor := services.NewImageProcessor(config.ImageProcessing)
        translator := services.NewTranslator(config.Translation)
        
        // Initialize the new improved summarizer
        summarizer := services.NewSummarizer(config.Summarization)
        logger.Info("Initialized new context-aware summarizer with improved NLP techniques")

        // Create service handler
        serviceHandler := services.NewServiceHandler(imageProcessor, translator, summarizer, logger)

        // Start REST API server
        restServer := api.NewRESTServer(config.REST.Port, serviceHandler, logger)
        go func() {
                if err := restServer.Start(); err != nil && err != http.ErrServerClosed {
                        logger.Fatal("Failed to start REST server", "error", err)
                }
        }()
        logger.Info("REST API server started", "port", config.REST.Port)

        // Start gRPC server
        grpcServer := api.NewGRPCServer(config.GRPC.Port, serviceHandler, logger)
        go func() {
                if err := grpcServer.Start(); err != nil {
                        logger.Fatal("Failed to start gRPC server", "error", err)
                }
        }()
        logger.Info("gRPC server started", "port", config.GRPC.Port)

        // Wait for termination signal
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        logger.Info("Shutdown signal received")

        // Create shutdown context with timeout
        shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer shutdownCancel()

        // Shutdown servers gracefully
        if err := restServer.Stop(shutdownCtx); err != nil {
                logger.Error("Failed to gracefully shutdown REST server", "error", err)
        }
        grpcServer.Stop()

        logger.Info("Ancient Script Decoder service stopped")
}
