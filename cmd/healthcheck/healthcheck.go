/*
Healthcheck is a simple program that sends an HTTP request to the local host (self) to a configured port number.
It's used in environment where you need a simple probe for health checks (e.g., an empty container in docker).
The probe URL is http://localhost:3000/liveness . Only the port can be changed.

Usage:

	healthcheck [flags]

The flags are:

	-port <1-65535>
		Change the port where the request is sent.

Return values (exit codes):

	0
		The request was successful (HTTP 200 or HTTP 204)

	> 0
		The request was not successful (connection error or unexpected HTTP status code)
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	port := flag.Int("port", 3000, "HTTP port for healthcheck")
	flag.Parse()

	// Context con timeout (best practice per healthcheck)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://localhost:%d/liveness", *port)

	// Creazione richiesta con context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating request:", err)
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 5 * time.Second, // ulteriore sicurezza
	}

	res, err := client.Do(req)
	if err != nil {
		// Errore nella richiesta (timeout, rifiuto, ecc.)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		fmt.Fprintln(os.Stderr, "Healthcheck request not OK:", res.Status)
		os.Exit(1)
	}

	os.Exit(0)
}