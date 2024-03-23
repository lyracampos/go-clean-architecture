package app

import (
	"fmt"
	"time"

	"github.com/lyracampos/go-clean-architecture/config"
)

func RunWorker(config *config.Config) {
	keepRunning := true
	count := 1
	for keepRunning {
		time.Sleep(2 * time.Second)
		fmt.Println("...")
		fmt.Println("worker is running...%i:", count)
		count++
	}
}
