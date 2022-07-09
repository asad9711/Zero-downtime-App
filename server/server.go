package server

import (
	"fmt"
	"listener_app/constants"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func Health(res http.ResponseWriter, req *http.Request) {
	brokers := strings.Split(constants.Brokers, ",")
	kafkaDialer := &kafka.Dialer{
		Timeout:   time.Duration(2000 * time.Millisecond),
		DualStack: true,
	}
	for _, broker := range brokers {
		log.Println("test-dial: kafka - trying to connect to broker - ", broker)
		conn, err := kafkaDialer.Dial("tcp", broker)
		if err != nil {
			// return 500, if the broker is unreachable. tested with locally hosted single kafka broker
			log.Println(fmt.Sprintf("test-dial: connect to kafka broker %s failed - %v", broker, err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = conn.Close(); err != nil {
			log.Println("error while closing connection to kafka - ", err.Error())
		}
	}

	log.Println("test-dial kafka passed")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}

func StartServer() {
	http.HandleFunc("/health", Health)
	log.Println("starting http server")

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Println("error while launching server - ", err.Error())
	} else {
		log.Println("launched server successfully at port 8090")
	}

}
