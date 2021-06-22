package db

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/knadh/koanf"
	"log"
	"time"
)

// dbConf contains database config required for connecting to a DB.
type dbConf struct {
	Host        string `koanf:"host"`
	Port        int    `koanf:"port"`
	User        string `koanf:"user"`
	Password    string `koanf:"password"`
	DBName      string `koanf:"database"`
	SSLMode     string `koanf:"ssl_mode"`
	MaxOpen     int    `koanf:"max_open"`
	MaxIdle     string `koanf:"max_idle"`
	MaxLifetime string `koanf:"max_lifetime"`
}

func InitDB(cfg *koanf.Koanf) (*pgxpool.Pool, error) {
	var dbCfg dbConf
	if err := cfg.Unmarshal("db", &dbCfg); err != nil {
		return nil, fmt.Errorf("error loading db config: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&pool_max_conns=%d&pool_max_conn_idle_time=%d&pool_max_conn_lifetime=%d",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode, dbCfg.MaxOpen, dbCfg.MaxIdle, dbCfg.MaxLifetime)

	log.Printf("connecting to db: %s:%d/%s", dbCfg.Host, dbCfg.Port, dbCfg.DBName)
	db, err := connectDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DB: %v", err)
	}
	log.Printf("checking for db migrations: %s:%d/%s", dbCfg.Host, dbCfg.Port, dbCfg.DBName)
	err = migrateDB(dbCfg)
	if err != nil {
		if err != migrate.ErrNoChange {
			log.Printf("db migration error: %s", err.Error())
			return nil, err
		} else {
			log.Printf("db migration: no changes detected")
		}

	}
	return db, nil
}

// The connectDB() function returns a pgx.Conn connection pool.
func connectDB(dsn string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error configuring the database: %v", err)
	}

	connectionPool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = connectionPool.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return connectionPool, nil
}

func migrateDB(dbCfg dbConf) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName, dbCfg.SSLMode)

	m, err := migrate.New(
		"file://./migrations",
		dsn)
	if err != nil {
		return fmt.Errorf("error configuring the database: %v", err)
	}
	err = m.Up()
	if err != nil {
		return err
	}
	return nil
}
