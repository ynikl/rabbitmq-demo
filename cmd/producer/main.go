package main

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	// step 1: connect to the rabbitmq broker
	connectionUrl := "amqp://ian:ian1234@localhost:5672/playground"
	conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// step 2: open channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// step 3: 声明 excahnge 这步可以省略, 直接使用默认的 excahnge
	err = ch.ExchangeDeclare("hello-exchange", "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// step 4: declare a queue
	q, err := ch.QueueDeclare("hello", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// step 5: 将 声明的 Exchange 绑定到 hello 队列, 并使用 hellokey 作为 bindingkey
	err = ch.QueueBind(q.Name, "hellokey", "hello-exchange", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("queue name: ", q.Name)

	body := "Hello 世界!"

	for {

		// step 6: 推送到 hello-exchange 的 交换器, 并使用 hellokey 做为 routingKey 发送消息
		err = ch.Publish("hello-exchange", "hellokey", false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("publish success...")
		time.Sleep(1 * time.Second)
	}

	fmt.Println("I'm done")
}
