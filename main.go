package main

import (
	"log"
	. "mq-consume/handlers"
	. "mq-consume/models"
	"time"
)

func main() {

	c, err := NewConsumer(false, Uri, Exchange, ExchangeType, Queue, BindingKey, ConsumerTag)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if Lifetime > 0 {
		log.Printf("running for %s", Lifetime)
		time.Sleep(Lifetime)
	} else {
		log.Printf("running forever %s", Lifetime)
		select {}
	}

	log.Printf("shutting down")

	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
}
