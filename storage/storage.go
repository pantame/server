package storage

import (
	"fmt"
	"github.com/pantame/server/config"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/storage/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

const NotFound = "storage: not found"

type MemoryStorage interface {
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetBytes(key string) ([]byte, error)
	Set(key string, value interface{}, exp time.Duration) error
	IncrBy(key string, value int64) error
	Delete(key string) error
	Reset() error
	Close() error
}

var (
	DB    *gorm.DB
	Cache MemoryStorage
)

func ConnectDB() {
	conf := config.Database()
	var err error
	DB, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		conf.Host, conf.User, conf.Pass, conf.Db, conf.Port)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func AutoMigrateDB() {
	log.Println("Rodando migrações...")
	DB.AutoMigrate(&entities.User{}, &entities.AccessPass{}, &entities.Session{}, &entities.IPData{}, &entities.File{})
}

func ConnectCache() {
	Cache = redis.New(config.Redis())
}
