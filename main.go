package main

import (
	"listener_app/constants"
	"listener_app/consumers"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"listener_app/server"

	"github.com/segmentio/kafka-go"
)

func customLogger(msg string, obj ...interface{}) {
	log.Println("kafka reader ===>> ", msg, obj)
}

func createTopic() {
	// to create topics when auto.create.topics.enable='false'
	topic := constants.TopicName

	conn, err := kafka.Dial("tcp", "broker:9092")
	if err != nil {
		log.Println("unable to connect to kafka broker at broker:9092")
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
		log.Println("error in creating topic")
		panic(err.Error())
	} else {
		log.Println("successfully created topic -> ", topic)
	}

}

// create connection to kafka, and start listening -> processing incoming messages
func startListener() {
	// brokers := strings.Split("host.docker.internal:9092", ",")
	// brokers := strings.Split("docker.for.mac.host.internal:9092", ",")
	// brokers := strings.Split("192.168.65.2:9092", ",")
	brokers := strings.Split(constants.Brokers, ",")
	log.Println("BROKER ***** ", brokers)

	hostname, _ := os.Hostname()
	config := kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        hostname,
		Topic:          constants.TopicName,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 200 * time.Second,
		RetentionTime:  time.Minute,
		ErrorLogger:    kafka.LoggerFunc(customLogger),
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	var wg sync.WaitGroup
	log.Println("Starting to consume messages")

	go consumers.ListenForMessage(reader, &wg)

	processKillSignal := make(chan os.Signal, 1)
	// to get notified if pod gets restarted/stopped
	signal.Notify(processKillSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSTOP)

	log.Println("waiting to receive on sigterm signal")

	<-processKillSignal

	log.Println("received process kill signal. Waiting for processes in the waitgroup to complete")
	wg.Wait()
	log.Println("All processes in waitgroup complete.  Exiting ...........")
}

func main() {
	// this service should expose a health endpoint, which k8 can use to ascertain the liveness of the pod
	go server.StartServer()
	startListener()
}
