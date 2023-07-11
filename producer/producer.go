package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type Data struct {
	Body string `yaml:"data"`
}

func main() {
	var info Data

	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		"localhost:9092",
		"my_topic2",
		0)

	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	err = conn.SetDeadline(time.Now().Add(time.Second * 2))
	if err != nil {
		fmt.Println(err)
		return
	}

	yamlFile, err := os.ReadFile("../data.yaml")

	if err != nil {
		fmt.Println(err)
		return
	}

	err = yaml.Unmarshal(yamlFile, &info)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.WriteMessages(kafka.Message{Value: []byte(info.Body)})
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
}
