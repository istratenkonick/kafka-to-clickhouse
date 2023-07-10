package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"gopkg.in/yaml.v3"
	"k2c/config"
	"k2c/pkg/abspath"
	"os"
	"strings"
)

func main() {
	var cfg *config.Config

	absPath := abspath.GetAbsolutePath()
	yamlFile, err := os.ReadFile(absPath + "cfg.yaml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)

	if err != nil {
		panic(err)
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: cfg.Clickhouse.Database,
			Username: cfg.Clickhouse.Username,
			Password: cfg.Clickhouse.Password,
		},
	})

	if err != nil {
		fmt.Println("Failed to connect to ClickHouse:", err)
		return
	}
	defer conn.Close()

	filePath := absPath + "migrations/create_table_logs.sql"

	sqlBytes, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error reading SQL file:", err)
		return
	}

	sqlStmts := strings.Split(string(sqlBytes), ";")

	for _, stmt := range sqlStmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		err := conn.Exec(context.TODO(), stmt)

		if err != nil {
			fmt.Println("Error executing SQL statement:", err)
			return
		}
	}

	fmt.Println("Table created successfully")
}
