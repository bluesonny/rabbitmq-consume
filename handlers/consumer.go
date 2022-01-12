package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/streadway/amqp"
	"log"
	. "rabbitmq-consume/models"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
	//连接异常结束
	ConnNotifyClose chan *amqp.Error
	//通道异常接收
	ChNotifyClose chan *amqp.Error
	//用于关闭进程
	CloseProcess chan bool
}

func NewConsumer(isClose bool, amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {

	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}
	if isClose {
		c.CloseProcess <- true
	}

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		log.Printf("Dial: %v", err)

	}
	/*
		go func() {
			fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		}()
	*/
	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	deliveries, err := c.channel.Consume(
		queue.Name, // name
		c.tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	/*
		go func() {
			time.Sleep(10 * time.Second)
			c.channel.Close()
			//fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		}()
	*/
	c.CloseProcess = make(chan bool, 1)
	go c.ConsumerReConnect()
	go handle(deliveries, c.done)

	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf("----业务队列开始处理----")
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		//业务处理

		msg := Message{}
		json.Unmarshal(d.Body, &msg)
		log.Printf(
			"got id %d name %s", msg.Id, msg.Type)
		var sta bool = do(msg)
		if sta == true {
			d.Ack(false)
		}
	}
	log.Printf("handle: deliveries channel closed")
	done <- nil
}

//ConsumerReConnect 消费者重连
func (c *Consumer) ConsumerReConnect() {
	log.Printf("重连开始......")
closeTag:
	for {
		log.Printf("重连开始进入循环等待")
		//fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		c.ConnNotifyClose = c.conn.NotifyClose(make(chan *amqp.Error))
		c.ChNotifyClose = c.channel.NotifyClose(make(chan *amqp.Error))
		//var err *amqp.Error
		select {
		case err := <-c.ConnNotifyClose:
			if err != nil {
				log.Printf("rabbit消费者连接异常:%v", err)
			}
			for err := range c.ConnNotifyClose {
				log.Println(err)
			}
			ants.Submit(func() {
				log.Printf("***连接重连")
				NewConsumer(false, Uri, Exchange, ExchangeType, Queue, BindingKey, ConsumerTag)
				log.Printf("***连接重连完毕")
			})
			break closeTag
		case err := <-c.ChNotifyClose:
			log.Printf("rabbit连接ch关闭异常:%v", err)
			if err != nil {
				log.Printf("rabbit连接ch关闭异常:%v", err)
			}
			// 判断连接是否关闭
			if !c.conn.IsClosed() {
				if err := c.conn.Close(); err != nil {
					log.Printf("rabbit连接关闭异常:%v", err)
				}
			}
			// IMPORTANT: 必须清空 Notify，否则死连接不会释放
			for err := range c.ChNotifyClose {
				log.Println(err)
			}
			ants.Submit(func() {
				log.Printf("***ch连接重连")
				NewConsumer(false, Uri, Exchange, ExchangeType, Queue, BindingKey, ConsumerTag)
				log.Printf("***ch连接重连完毕")
			})
			break closeTag

		case <-c.CloseProcess:
			break closeTag
		}
	}
	log.Printf("结束消费者旧进程")
}
