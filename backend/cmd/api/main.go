package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"

	"github.com/fiorellizz/gopayflow/internal/application"
	"github.com/fiorellizz/gopayflow/internal/infrastructure/database"
	httpInterface "github.com/fiorellizz/gopayflow/internal/interfaces/http"
	"github.com/streadway/amqp"
	"github.com/fiorellizz/gopayflow/internal/infrastructure/messaging"
)

func main() {

	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rabbitURL := os.Getenv("RABBIT_URL")

	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	queue, err := channel.QueueDeclare(
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

	publisher := messaging.NewRabbitMQPublisher(channel, queue.Name)

	orderRepo := database.NewPostgresOrderRepository(db)

	createOrderUseCase := application.NewCreateOrderUseCase(orderRepo, publisher)
	getOrderByIDUseCase := application.NewGetOrderByIDUseCase(orderRepo)
	listOrdersUseCase := application.NewListOrdersUseCase(orderRepo)

	orderHandler := httpInterface.NewOrderHandler(
		createOrderUseCase,
		getOrderByIDUseCase,
		listOrdersUseCase,
	)

	router := gin.Default()

	router.SetTrustedProxies([]string{"nginx"})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
	})

	router.POST("/orders", orderHandler.CreateOrder)
	router.GET("/orders", orderHandler.ListOrders)
	router.GET("/orders/:id", orderHandler.GetOrderByID)

	log.Println("Starting server on port 8080")
	router.Run(":8080")
}