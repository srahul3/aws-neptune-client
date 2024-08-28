package main

import (
	"fmt"
	"os"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
)

func main() {

	// get metune endpoint from environment variable
	endpoint := os.Getenv("NEPTUNE_ENDPOINT")

	// Creating the connection to the server.
	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection("wss://"+endpoint+":8182/gremlin",
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "g"
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	// Cleanup
	defer driverRemoteConnection.Close()

	// Creating graph traversal
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	// Perform traversal
	results, err := g.V().Limit(2).ToList()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print results
	for _, r := range results {
		fmt.Println(r.GetString())
	}
}
