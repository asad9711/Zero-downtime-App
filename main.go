package main

import (
	"fmt"
	"listener_app/constants"
	"listener_app/consumers"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func customLogger(msg string, obj ...interface{}) {
	fmt.Println("kafka reader ===>> ", msg, obj)
}

func createTopic() {
	// to create topics when auto.create.topics.enable='false'
	topic := constants.TopicName

	conn, err := kafka.Dial("tcp", "broker:9092")
	if err != nil {
		fmt.Println("unable to connect to kafka broker at broker:9092")
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		fmt.Println("error in creating topic")
		panic(err.Error())
	} else {
		fmt.Println("successfully created topic -> ", topic)
	}

}

func startListener() {
	// fmt.Println("sleeping for 60 secs")
	// time.Sleep(60 * time.Second)

	// brokers := strings.Split("host.docker.internal:9092", ",")
	// brokers := strings.Split("docker.for.mac.host.internal:9092", ",")
	// brokers := strings.Split("192.168.65.2:9092", ",")
	brokers := strings.Split(constants.BROKERS, ",")
	fmt.Println("BROKER ***** ", brokers)

	// fmt.Println("sleeping for 5 secs")
	// time.Sleep(5 * time.Second)

	hostname, _ := os.Hostname()
	config := kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        hostname,
		Topic:          constants.TopicName,
		MinBytes:       1,    // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 200 * time.Second,
		RetentionTime:  time.Minute,
		Logger:         kafka.LoggerFunc(customLogger),
		ErrorLogger:    kafka.LoggerFunc(customLogger),
		//MaxWait:         100 * time.Microsecond, // Maximum amount of time to wait for new data to come when fetching batches of messages dividerChannel kafka.
		//ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	var wg sync.WaitGroup
	fmt.Println("Starting to consume messages")

	go consumers.ListenForMessage(reader, &wg)

	processKillSignal := make(chan os.Signal, 1)
	signal.Notify(processKillSignal, syscall.SIGINT)

	fmt.Println("waiting to receive on sigterm signal")

	<-processKillSignal

	fmt.Println("received process kill signal. Waiting for processes in the waitgroup to complete")
	wg.Wait()
	fmt.Println("All processes in waitgroup complete.  Exiting ...........")
}
func main() {

	startListener()
}
