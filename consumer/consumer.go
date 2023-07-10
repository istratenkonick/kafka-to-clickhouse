package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"k2c/config"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	yamlFile, err := os.ReadFile("../config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg *config.Config

	err = yaml.Unmarshal(yamlFile, &cfg)

	wg := sync.WaitGroup{}
	wg.Add(2)

	for _, topics := range cfg.Core.Files {
		for _, topic := range topics.Topics {
			go read(topic, &wg, cfg)
		}
	}

	wg.Wait()
}

func read(topic string, wg *sync.WaitGroup, config *config.Config) {
	defer wg.Done()

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, 0)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	for {
		n, err := conn.ReadMessage(1e3)
		if err != nil {
			break
		}

		db, err := clickhouse.Open(&clickhouse.Options{
			Addr: []string{"localhost:9000"},
			Auth: clickhouse.Auth{
				Database: config.Clickhouse.Database,
				Username: config.Clickhouse.Username,
				Password: config.Clickhouse.Password,
			},
		})

		if err != nil {
			fmt.Println("failed to connect to DB: ", err)
		}

		_ = db.Exec(
			context.Background(),
			"INSERT INTO logs (name, body, timestamp) VALUES (?, ?, ?)",
			"log_name",
			string(n.Value),
			time.Now().String())
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

	fmt.Println("SUCCESS")
}
