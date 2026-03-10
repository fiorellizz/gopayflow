package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"

	"github.com/fiorellizz/gopayflow/internal/domain"
	"github.com/fiorellizz/gopayflow/internal/infrastructure/database"
)

func main() {

	dbURL := os.Getenv("DB_URL")
	rabbitURL := os.Getenv("RABBIT_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	_, err = channel.QueueDeclare(
		"orders_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := channel.Consume(
		"orders_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	repo := database.NewPostgresOrderRepository(db)

	log.Println("Worker started")

	for msg := range msgs {

		var order domain.Order

		err := json.Unmarshal(msg.Body, &order)
		if err != nil {
			log.Println(err)
			continue
		}

		status := processPayment(order.Amount)

		err = repo.UpdateStatus(context.Background(), order.ID, status)
		if err != nil {
			log.Println(err)
		}

		log.Println("order processed:", order.ID, status)
	}
}

func processPayment(amount float64) domain.OrderStatus {

	time.Sleep(2 * time.Second)

	if rand.Intn(2) == 0 {
		return domain.StatusApproved
	}

	return domain.StatusFailed
}