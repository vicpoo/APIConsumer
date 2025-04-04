// messaging_service.go
package infrastructure

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type MessagingService struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	hub  *Hub
}

const (
	exchangeName = "coffee_bed_orders"
	queueName    = "coffee_bed_orders_queue"
	routingKey   = "order.created"
)

func NewMessagingService(hub *Hub) *MessagingService {
	conn, err := amqp.Dial("amqp://reyhades:reyhades@44.223.218.9:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return nil
	}

	// Asegurar que la cola está declarada
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return nil
	}

	return &MessagingService{
		conn: conn,
		ch:   ch,
		hub:  hub,
	}
}

func (ms *MessagingService) ConsumeOrderMessages() error {
	msgs, err := ms.ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (false para reconocimiento manual)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received order message: %s", string(msg.Body))

			// Parsear el JSON recibido
			var orderData map[string]interface{}
			if err := json.Unmarshal(msg.Body, &orderData); err == nil {
				log.Printf("Parsed order data: %+v", orderData)
				// Aquí podrías procesar la orden recibida, guardarla en la base de datos, etc.
			} else {
				log.Printf("Error parsing order JSON: %v", err)
			}

			// Enviar a los clientes WebSocket
			ms.hub.broadcast <- msg.Body

			// Confirmar el mensaje
			msg.Ack(false)
		}
	}()

	return nil
}

func (ms *MessagingService) Close() {
	if ms.ch != nil {
		ms.ch.Close()
	}
	if ms.conn != nil {
		ms.conn.Close()
	}
}
