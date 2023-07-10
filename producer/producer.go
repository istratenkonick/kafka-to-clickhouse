package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Data struct {
	Body string `yaml:"data"`
}

func main() {
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		"localhost:9092",
		"my_topic2",
		0)

	if err != nil {
		fmt.Println(err)
		return
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

	var info Data

	err = yaml.Unmarshal(yamlFile, &info)

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = conn.WriteMessages(kafka.Message{Value: []byte(info.Body)})
	if err != nil {
		fmt.Println(err)
		return
	}
}
