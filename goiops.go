package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	tmpFile, err := ioutil.TempFile(".", "goiops-")
	if err != nil {
		log.Fatalf("Error opening tmpFile: %s", err)
	}

	randBytes := make([]byte, 16)

	for {
		select {
		case <-signals:
			log.Print("Received termination signal.")
			return
		default:
			preRun := time.Now()

			for i := 0; i < 500; i++ {
				if _, err := rand.Read(randBytes); err != nil {
					log.Printf("Error reading random bytes: %s", err)
				}

				if _, err := tmpFile.Write(randBytes); err != nil {
					log.Fatalf("Error writing to tmpFile: %s", err)
				}

				time.Sleep(time.Second / 500)
			}

			preSync := time.Now()

			if err := tmpFile.Sync(); err != nil {
				log.Fatalf("Error syncing tmpFile: %s", err)
			}

			syncTime := time.Since(preSync)

			if _, err := tmpFile.Seek(0, 0); err != nil {
				log.Fatalf("Error seeking tmpFile: %s", err)
			}

			if err := tmpFile.Truncate(0); err != nil {
				log.Fatalf("Error truncating tmpFile: %s", err)
			}

			runTime := time.Since(preRun)
			log.Printf("Whole: %.2fs, Sync: %7.2fms", float64(runTime)/float64(time.Second), float64(syncTime)/float64(time.Millisecond))
		}
	}
}
