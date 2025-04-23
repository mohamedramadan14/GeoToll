package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/mohamedramadan14/roads-fees-system/utilities"
)

func main() {
	utilities.InitLogger()
	fmt.Printf("This is the invoicer service\n")
	listenAddr := flag.String("listen-addr", ":3100", "HTTP listen address")
	flag.Parse()
	var (
		store          = NewMemoryStore()
		svc            = NewInvoiceAggregator(store)
		svcWithLogging = NewLogMiddleware(svc)
	)

	err := makeHTTPTransport(*listenAddr, svcWithLogging)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		return
	}
}

func HandleAggregate(svc Aggregator) http.HandlerFunc {
	// Handle the aggregation of invoices
	// This is where you would implement the logic to aggregate invoices
	// and return the result to the client
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload : " + err.Error()})
			return
		}

		if err := svc.AggregateDistances(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

	}
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Printf("HTTP Transport running on port: %s\n", listenAddr)
	http.HandleFunc("POST /aggregate", HandleAggregate(svc))
	return http.ListenAndServe(listenAddr, nil)
}

func writeJSON(rw http.ResponseWriter, status int, data any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(data)
}
