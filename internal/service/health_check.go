package service

import (
	"context"
	"database/sql"
	"fmt"
	"middleware-platform/internal/model"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-zookeeper/zk"
	_ "github.com/lib/pq"
)

func (s *MiddlewareService) checkRedisHealth(m *model.Middleware) error {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", m.Host, m.Port),
		Password: m.Credentials,
		DB:       0,
	})
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	return client.Ping(ctx).Err()
}

func (s *MiddlewareService) checkMySQLHealth(m *model.Middleware) error {
	dsn := fmt.Sprintf("user:password@tcp(%s:%s)/", m.Host, m.Port)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	
	return db.Ping()
}

func (s *MiddlewareService) checkPgHealth(m *model.Middleware) error {
	dsn := fmt.Sprintf("host=%s port=%s user=postgres sslmode=disable", m.Host, m.Port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	
	return db.Ping()
}

func (s *MiddlewareService) checkZKHealth(m *model.Middleware) error {
	conn, _, err := zk.Connect([]string{fmt.Sprintf("%s:%s", m.Host, m.Port)}, time.Second*5)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	return nil
} 