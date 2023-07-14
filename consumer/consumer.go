package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"k2c/config"
	"k2c/pkg/abspath"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	absPath := abspath.GetAbsolutePath()
	yamlFile, err := os.ReadFile(absPath + "kafka-to-clickhouse/config.yaml")
	if err != nil {
		panic(err)
	}

	var cfg *config.Config

	err = yaml.Unmarshal(yamlFile, &cfg)

	wg := sync.WaitGroup{}

	for _, topics := range cfg.Core.Files {
		for _, topic := range topics.Topics {
			wg.Add(1)
			go read(topic, &wg, cfg)
		}
	}

	wg.Wait()
}

func read(topic string, wg *sync.WaitGroup, config *config.Config) {
	defer wg.Done()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer r.Close()

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
	defer db.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("failed to read message: ", err)
			continue
		}

		err = db.Exec(
			context.Background(),
			"INSERT INTO logs (name, body, timestamp) VALUES (?, ?, ?)",
			"log_name",
			string(m.Value),
			time.Now().Format(time.DateTime))
		if err != nil {
			log.Println("failed to insert into DB: ", err)
			continue
		}
	}
}
