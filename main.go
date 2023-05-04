package main

import (
	"flag"
	"fmt"
	"net/http"
	"zartekAssignment/variables"
	"zartekAssignment/visitors"
)

func init() {
	flag.IntVar(&variables.MaxRequests, "requests", variables.MaxRequests, "number of requests allowed per IP address in a given duration")
	flag.Parse()
}

func main() {
	v := visitors.NewVisitors()

	server := &http.Server{
		Addr:    ":8080",
		Handler: v,
	}

	fmt.Print("Server started at port http://localhost:8080\n")
	fmt.Printf("Max requests allowed per IP address in %v: %v\n", variables.Duration, variables.MaxRequests)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

	v.Wg.Wait()
}
