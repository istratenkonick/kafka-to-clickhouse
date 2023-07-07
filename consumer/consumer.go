package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
	"time"
)

type KafkaTopics struct {
	Main   string `yaml:"main_topic"`
	Second string `yaml:"second_topic"`
}

type Topic struct {
	Names KafkaTopics `yaml:"topics"`
}

func main() {
	yamlFile, err := os.ReadFile("../config.yaml")
	if err != nil {
		panic(err)
	}

	var topic Topic

	err = yaml.Unmarshal(yamlFile, &topic)

	topics := []string{topic.Names.Main, topic.Names.Second}
	wg := sync.WaitGroup{}
	wg.Add(len(topics))

	for _, topic := range topics {
		go read(topic, &wg)
	}

	wg.Wait()
}

func read(topic string, wg *sync.WaitGroup) {
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
				Database: "k2c",
				Username: "",
				Password: "",
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
