package database

import (
	"context"
	"fmt"
	"time"

	"backend/internal/config"

	"github.com/charmbracelet/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Connection Pool Settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("Koneksi database PostgreSQL berhasil diinisialisasi")
	return db, nil
}

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "", // tidak ada password default di docker compose
		DB:       0,  // default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke Redis: %w", err)
	}

	log.Info("Koneksi Redis berhasil diinisialisasi")
	return rdb, nil
}

func InitMinio(cfg *config.Config) (*minio.Client, error) {
	// Minio Host biasanya localhost:9000, useSSL = false
	minioClient, err := minio.New(cfg.MinioHost, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioUser, cfg.MinioPass, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal menginisialisasi MinIO client: %w", err)
	}

	log.Info("Koneksi MinIO Object Storage berhasil diinisialisasi")
	return minioClient, nil
}
