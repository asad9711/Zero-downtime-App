package business_logic

import (
	"log"
	"sync"
	"time"
)

func AsyncProcessMessage(msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		log.Println("logging message -  time - ", i, msg)
		time.Sleep(1 * time.Second)
	}
}
