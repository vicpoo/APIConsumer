// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github/vicpoo/APIConsumer/Consumer/infrastructure"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configurar Gin
	r := gin.Default()

	// Inicializar el hub de WebSocket y el consumidor
	hub := infrastructure.NewHub()
	go hub.Run()

	// Configurar servicio de mensajería
	messagingService := infrastructure.NewMessagingService(hub)
	defer messagingService.Close()

	// Configurar rutas
	infrastructure.SetupRoutes(r, hub)

	// Iniciar consumidor de RabbitMQ para órdenes
	if err := messagingService.ConsumeOrderMessages(); err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}

	// Manejar señales para apagado limpio
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en una goroutine
	go func() {
		if err := r.Run(":8005"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started on port 8005")
	log.Println("Order consumer started")

	// Esperar señal de apagado
	<-sigChan
	log.Println("Shutting down server...")
}
