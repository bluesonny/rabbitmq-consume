package models

import (
	"fmt"
	. "mq-consume/config"
	"time"
)

var (
	Uri = fmt.Sprintf("amqp://%s:%s@%s/", ViperConfig.Mq.User, ViperConfig.Mq.Password, ViperConfig.Mq.Host)
	//flag.String("uri", "amqp://admintest:admin123@47.104.82.28:5672/", "AMQP URI")
	Exchange     = "amq.direct"         // flag.String("exchange", "amq.direct", "Durable, non-auto-deleted AMQP exchange name")
	ExchangeType = "direct"             //flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	Queue        = ViperConfig.Mq.Queue //flag.String("queue", "test", "Ephemeral AMQP queue name")
	BindingKey   = ViperConfig.Mq.Queue // flag.String("key", "queue", "AMQP binding key")
	ConsumerTag  = "consumer"           //flag.String("consumer-tag", "consumer", "AMQP consumer tag (should not be blank)")
	Lifetime     = -1 * time.Second     //flag.Duration("lifetime", -1*time.Second, "lifetime of process before shutdown (0s=infinite)")
)
