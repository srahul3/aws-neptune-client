// package main

// import (
// 	"fmt"
// 	"os"

// 	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
// )

// func main() {

// 	// get metune endpoint from environment variable
// 	endpoint := os.Getenv("NEPTUNE_ENDPOINT")

// 	// Creating the connection to the server.
// 	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection("wss://"+endpoint+":8182/gremlin",
// 		func(settings *gremlingo.DriverRemoteConnectionSettings) {
// 			settings.TraversalSource = "g"
// 		})
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Cleanup
// 	defer driverRemoteConnection.Close()

// 	// Creating graph traversal
// 	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

// 	// Perform traversal
// 	results, err := g.V().Limit(2).ToList()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Print results
// 	for _, r := range results {
// 		fmt.Println(r.GetString())
// 	}
// }

// package main

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
// 	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
// 	"github.com/aws/aws-sdk-go-v2/config"
// )

// func main() {
// 	ctx := context.Background()
// 	db_endpoint := os.Getenv("NEPTUNE_ENDPOINT")
// 	connString := "wss://" + db_endpoint + ":8181/gremlin"

// 	// Signing Request
// 	const emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
// 	req, err := http.NewRequest(http.MethodGet, connString, strings.NewReader(""))
// 	if err != nil {
// 		return
// 	}

// 	cfg, err := config.LoadDefaultConfig(ctx)
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("unable to load AWS SDK config: %w", err))
// 	}
// 	cr, err := cfg.Credentials.Retrieve(ctx)
// 	if err != nil {
// 		fmt.Println(fmt.Errorf("unable to retrieve AWS credentials: %w", err))
// 	}
// 	signer := v4.NewSigner()

// 	gen := func() gremlingo.AuthInfoProvider {
// 		err := signer.SignHTTP(ctx, cr, req, emptyStringSHA256, "db-neptune-ig", "us-west-2", time.Now())

// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		return gremlingo.HeaderAuthInfo(req.Header)
// 	}

// 	auth := gremlingo.NewDynamicAuth(gen)

// 	fmt.Println("Connecting to Neptune endpoint: ", connString)
// 	fmt.Println("Using auth: ", auth)

// 	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection(connString,
// 		func(settings *gremlingo.DriverRemoteConnectionSettings) {
// 			settings.TraversalSource = "g"
// 			settings.AuthInfo = auth
// 			//settings.TlsConfig = &tls.Config{InsecureSkipVerify: true} Use this only if you're on a Mac running Go 1.18+ doing local dev. See https://github.com/golang/go/issues/51991

// 		})

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Cleanup
// 	defer driverRemoteConnection.Close()

// 	// Creating graph traversal
// 	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

// 	// Add a vertex with properties to the graph with the terminal step Iterate()
// 	promise := g.AddV("gremlin").Property("language", "go").Iterate()

// 	// The returned promised is a go channel to wait for all submitted steps to finish execution and return error.
// 	err = <-promise
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Get the value of the property
// 	result, err := g.V().HasLabel("gremlin").Values("language").ToList()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	// Print the result
// 	for _, r := range result {
// 		fmt.Println(r.GetString())
// 	}
// }

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	uri := "bolt+ssc://" + os.Getenv("NEPTUNE_ENDPOINT") + ":8182/opencypher"
	fmt.Println("Connecting to: ", uri)

	// Create a driver without authentication
	// driver, err := neo4j.NewDriverWithContext(uri, neo4j.NoAuth(), func(config *neo4j.Config) {
	// 	config.Encrypted = true
	// 	config.TrustStrategy = neo4j.TrustSystemCertificates // or TrustAllCertificates if Neptune uses self-signed certificates
	// })
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.NoAuth())
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}
	defer driver.Close(context.Background())

	// Open a new session
	session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
	defer session.Close(context.Background())

	// Example query
	cypherQuery := "MATCH (n) RETURN n LIMIT 10"

	// Run the query
	result, err := session.Run(context.Background(), cypherQuery, nil)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	// Process the result
	for result.Next(context.Background()) {
		record := result.Record()
		fmt.Println(record.Values)
	}

	// Check for errors in the result
	if err = result.Err(); err != nil {
		log.Fatalf("Result error: %v", err)
	}

}
