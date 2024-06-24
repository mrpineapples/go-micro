package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err := ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			payload := Payload{}
			json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message on [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log what we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// authenticate
		// add as many cases as long as logic exists
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	data, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"
	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		return err
	}

	return nil

}
