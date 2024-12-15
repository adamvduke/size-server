package main

import (
	"flag"
	"fmt"
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
	flag.Parse()

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, "", log.LstdFlags)
	if *client {
		work := make(chan int)
		for i := 0; i < 1; i++ {
			go func(w chan int) {
				for {
					i := <-w
					u := fmt.Sprintf("http://size.docker.localhost:80/?size=%d", i)
					resp, err := http.Get(u)
					if err != nil {
						logger.Println(err)
					}
					logger.Printf("Requested size: %v, Status: %v, Length: %v", i, resp.Status, resp.ContentLength)
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
