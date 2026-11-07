package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce/backend-service/configs"
	"ecommerce/backend-service/internal/auth"
	"ecommerce/backend-service/internal/cart"
	"ecommerce/backend-service/internal/database"
	"ecommerce/backend-service/internal/middleware"
	"ecommerce/backend-service/internal/order"
	"ecommerce/backend-service/internal/payment"
	"ecommerce/backend-service/internal/product"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func main() {
	cfg := configs.LoadConfig()

	// 1. Initialize PostgreSQL with GORM
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Fatal: Database initialization failed: %v", err)
	}

	// 2. Seed initial demo data if empty
	database.SeedInitialData(db)

	// 3. Instantiate Repositories
	authRepo := auth.NewRepository(db)
	prodRepo := product.NewRepository(db)
	cartRepo := cart.NewRepository(db)
	orderRepo := order.NewRepository(db)
	payRepo := payment.NewRepository(db)

	// 4. Instantiate Clients & Services
	dokuClient := payment.NewDokuClient(cfg.DokuClientID, cfg.DokuSecretKey, cfg.DokuIsProduction, cfg.DokuSimulateMode)

	authSvc := auth.NewService(authRepo, cfg.JWTSecret, cfg.RefreshTokenSecret)
	prodSvc := product.NewService(prodRepo)
	cartSvc := cart.NewService(cartRepo, prodRepo)
	orderSvc := order.NewService(orderRepo)
	paySvc := payment.NewService(payRepo, dokuClient, cfg)

	// 5. Instantiate Handlers
	authHandler := auth.NewHandler(authSvc)
	prodHandler := product.NewHandler(prodSvc)
	cartHandler := cart.NewHandler(cartSvc)
	orderHandler := order.NewHandler(orderSvc)
	payHandler := payment.NewHandler(paySvc)

	// 6. Router & Middlewares
	mux := http.NewServeMux()

	authMw := middleware.Auth(cfg.JWTSecret)
	buyerMw := middleware.RequireRole("buyer")
	sellerMw := middleware.RequireRole("seller")

	// Health check endpoints
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(HealthResponse{
			Status:  "ok",
			Service: "backend-service",
		})
	}
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/v1/health", healthHandler)

	// Register Domain Routes
	authHandler.RegisterRoutes(mux, authMw)
	prodHandler.RegisterRoutes(mux, authMw, sellerMw)
	cartHandler.RegisterRoutes(mux, authMw, buyerMw)
	orderHandler.RegisterRoutes(mux, authMw, buyerMw, sellerMw)
	payHandler.RegisterRoutes(mux, authMw, buyerMw)

	// Global Middleware Chain: Recovery -> Logger -> CORS -> ServeMux
	var rootHandler http.Handler = mux
	rootHandler = middleware.CORS(cfg.FrontendURL)(rootHandler)
	rootHandler = middleware.Logger(rootHandler)
	rootHandler = middleware.Recovery(rootHandler)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      rootHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Run Server asynchronously
	go func() {
		log.Printf("Server listening on http://localhost%s (Frontend origin: %s)\n", addr, cfg.FrontendURL)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server ListenAndServe error: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped cleanly")
}
