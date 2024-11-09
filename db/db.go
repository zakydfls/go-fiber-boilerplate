package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type DB struct {
	*sql.DB
}

var db *gorm.DB

func Init() {
	dbInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	var err error
	db, err = ConnectDB(dbInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func ConnectDB(dataSourceName string) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	gormDB.Use(
		dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{postgres.Open(dataSourceName)},
			// Replicas: []gorm.Dialector{postgres.Open(dataSourceName), postgres.Open(dataSourceName)},
			Policy: dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(100).
			SetMaxOpenConns(150).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour),
	)

	//dbmap.TraceOn("[gorp]", log.New(os.Stdout, "golang-gin:", log.Lmicroseconds)) //Trace database requests
	return gormDB, nil
}

func GetDB() *gorm.DB {
	return db
}

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(selectDB ...int) {
	var redisHost = os.Getenv("REDIS_HOST")
	var redisPassword = os.Getenv("REDIS_PASSWORD")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       selectDB[0],
		// Optional configuration
		// DialTimeout:        10 * time.Second,
		// ReadTimeout:        30 * time.Second,
		// WriteTimeout:       30 * time.Second,
		// PoolSize:           10,
		// PoolTimeout:        30 * time.Second,
		// IdleTimeout:        500 * time.Millisecond,
		// IdleCheckFrequency: 500 * time.Millisecond,
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: true,
		// },
	})
}

// GetRedis ...
func GetRedis() *redis.Client {
	return RedisClient
}

func SetRedisCache(key string, value interface{}, expiredTime time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, key, p, expiredTime).Err()
}

func GetRedisCache(key string, dest interface{}) error {
	p, err := RedisClient.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(p, dest)
}

func DelRedisCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}
