package business_logic

import (
	"fmt"
	"sync"
	"time"
)

func AsyncProcessMessage(msg string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		fmt.Println("Printing message -  time - ", i, msg)
		time.Sleep(1 * time.Second)
	}
}
