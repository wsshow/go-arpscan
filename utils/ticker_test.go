package utils

import (
	"log"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	tc := NewTicker(WithResetTime(5 * time.Second))
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			log.Println("reset")
			tc.Reset()
		}
	}()
	tc.Wait()
	log.Println("end")
}
