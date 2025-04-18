package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	client := flag.Bool("client", false, "if the program should run as a client rather than server")
	start := flag.Int("start", 1, "initial request size")
	batchSize := flag.Int("batch", 10000, "how many requests to make")
	workers := flag.Int("workers", 5, "number of workers to spawn")
	flag.Parse()

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags)
	if *client {
		work := make(chan int)
		for range *workers {
			go func(w chan int) {
				c := &http.Client{}
				for {
					i := <-w
					u := fmt.Sprintf("http://size.docker.localhost:80/?size=%d", i)
					resp, err := c.Get(u)
					if err != nil {
						logger.Println(err)
					}
					b, err := io.ReadAll(resp.Body)
					if err != nil {
						logger.Println(err)
					}
					logger.Printf("Requested size: %v, Status: %v, Length: %v", i, resp.Status, len(b))
				}
			}(work)
		}
		stop := *start + *batchSize
		for i := *start; i < stop; i++ {
			work <- i
		}
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			size := r.URL.Query().Get("size")
			s, err := strconv.Atoi(size)
			if err != nil {
				s = 100
			}
			b := strings.Repeat("A", s)
			fmt.Fprint(w, string(b))
		})
		logger.Fatal(http.ListenAndServe(":8080", nil))
	}
}
