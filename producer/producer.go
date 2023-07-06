package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
)

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

	_, err = conn.WriteMessages(kafka.Message{Value: []byte("This is second topic")})
	if err != nil {
		fmt.Println(err)
		return
	}
}
