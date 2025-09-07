package database

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	 "net/smtp"
)

type mysql struct {
	client *gorm.DB
}

type Myredis struct {
	client *redis.Client
}
type mongodb struct {
	client *mongo.Client
}
type database struct {
	mysql   *mysql
	redis   *Myredis
	mongodb *mongodb
}
type STMP struct {
	client *smtp.Client
}

var (
	DB = &database{}
)
