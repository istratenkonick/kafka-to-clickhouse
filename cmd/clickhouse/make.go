package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"io/ioutil"
	"strings"
)

func main() {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "miniOPcore",
			Username: "",
			Password: "",
		},
	})

	if err != nil {
		fmt.Println("Failed to connect to ClickHouse:", err)
		return
	}
	defer conn.Close()

	filePath := "../../migrations/create_table_logs.sql"
	sqlBytes, err := ioutil.ReadFile(filePath)

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
