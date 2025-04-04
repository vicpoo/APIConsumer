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

// Middleware de ejemplo para pruebas
func TestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Permitir solicitudes desde cualquier origen para pruebas (CORS)
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Si es una solicitud OPTIONS, responder inmediatamente
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// Continuar con la siguiente función en la pila de middlewares
		c.Next()
	}
}

func main() {
	// Configurar Gin
	r := gin.Default()

	// Agregar el middleware de prueba
	r.Use(TestMiddleware())

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
