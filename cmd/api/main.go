package main

import (
	"fmt"
	"os"

	"web-coffee/config"
	"web-coffee/internal/handler"
	"web-coffee/pkg/database"
	"web-coffee/pkg/redis"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 2. Inisialisasi koneksi MySQL
	db, err := database.NewMySQL(cfg.MysqlDSN())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to MySQL: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Println("MySQL connected.")

	// 3. Inisialisasi koneksi Redis
	rdb, err := redis.NewRedis(cfg.RedisAddr(), cfg.RedisPassword)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to Redis: %v\n", err)
		os.Exit(1)
	}
	defer rdb.Close()
	fmt.Println("Redis connected.")

	// 4. Inisialisasi Gin Router
	r := handler.NewRouter(cfg.AppEnv)

	// 5. Jalankan HTTP Server
	fmt.Println("Starting server on port :" + cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		fmt.Fprintf(os.Stderr, "Server failed to run: %v\n", err)
		os.Exit(1)
	}
}
