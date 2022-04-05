package consumers

import (
	"context"
	"fmt"
	"listener_app/business_logic"
	"sync"

	"github.com/segmentio/kafka-go"
)

func ListenForMessage(reader *kafka.Reader, wg *sync.WaitGroup) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		wg.Add(1)
		go business_logic.AsyncProcessMessage(string(m.Value), wg)

	}
}
