package main

import (
	"flag"
	"message_callback/g"
	"github.com/streadway/amqp"
	"log"
	"message_callback/utils"
	"fmt"
)


func setupQueueChannel(queue g.QueueConfig){
	log.Printf("开始设置队列: %s\n", queue.Name)

	conn, err := amqp.Dial("amqp://root:root@116.196.102.109:5672/")
	utils.PanicOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.PanicOnError(err, "Failed to open a channel")
	defer ch.Close()

	errorQueueName := fmt.Sprintf("%s_error", queue.Name)
	_, err = ch.QueueDeclare(
		errorQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	utils.PanicOnError(err, "Failed to declare a queue: " + errorQueueName)


	retryQueueName := fmt.Sprintf("%s_retry", queue.Name)
	exchangeForRetryQueue := fmt.Sprintf("%s_ex", retryQueueName)
	err = ch.ExchangeDeclare(
		exchangeForRetryQueue,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.PanicOnError(err, "Failed to declare an exchange")

	exchangeForRequeue := fmt.Sprintf("%s_ex", queue.Name)
	err = ch.ExchangeDeclare(
		exchangeForRequeue,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.PanicOnError(err, "Failed to declare an exchange")

	// 工作队列
	log.Println("开始初始化工作队列" +  exchangeForRetryQueue)
	workerQueueOptions := map[string]interface{}{
		"x-dead-letter-exchange": exchangeForRetryQueue,
	}

	_, err = ch.QueueDeclare(
		queue.Name,
		true,
		false,
		false,
		false,
		workerQueueOptions,
	)
	utils.PanicOnError(err, "Failed to declare a queue: " + queue.Name)

	// 重试队列
	log.Println("开始初始化重试队列")
	retryQueueOptions := map[string]interface{}{
		"x-dead-letter-exchange": exchangeForRequeue,
		"x-message-ttl":          3000,
	}
	_, err = ch.QueueDeclare(
		retryQueueName,
		true,
		false,
		false,
		false,
		retryQueueOptions,
	)
	utils.PanicOnError(err, "Failed to declare a queue: " + retryQueueName)

	// 绑定exchange和对应的queue
	err = ch.QueueBind(queue.Name, "#", exchangeForRequeue, false, nil)
	utils.PanicOnError(err, "Failed to QueueBind")
}


func setupChannel() {
	log.Println("setupChannel")

	config := g.GetConfig()
	log.Println(config.Global.Retry, config.Global.RetryTTL)

	for _, queue := range config.Queues {
		setupQueueChannel(queue)
	}

}

func main() {
	cfg := flag.String("c", "config.json", "configuration file")
	flag.Parse()

	g.ParseConfig(*cfg)

	setupChannel()
}
