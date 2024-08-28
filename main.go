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

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	gremlingo "github.com/apache/tinkerpop/gremlin-go/v3/driver"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	ctx := context.Background()
	db_endpoint := os.Getenv("NEPTUNE_ENDPOINT")
	connString := "wss://" + db_endpoint + ":8181/gremlin"

	// Signing Request
	const emptyStringSHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`
	req, err := http.NewRequest(http.MethodGet, connString, strings.NewReader(""))
	if err != nil {
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to load AWS SDK config: %w", err))
	}
	cr, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to retrieve AWS credentials: %w", err))
	}
	signer := v4.NewSigner()

	gen := func() gremlingo.AuthInfoProvider {
		err := signer.SignHTTP(ctx, cr, req, emptyStringSHA256, "db-neptune-ig", "us-west-2", time.Now())

		if err != nil {
			fmt.Println(err)
		}
		return gremlingo.HeaderAuthInfo(req.Header)
	}

	auth := gremlingo.NewDynamicAuth(gen)

	fmt.Println("Connecting to Neptune endpoint: ", connString)
	fmt.Println("Using auth: ", auth)

	driverRemoteConnection, err := gremlingo.NewDriverRemoteConnection(connString,
		func(settings *gremlingo.DriverRemoteConnectionSettings) {
			settings.TraversalSource = "g"
			settings.AuthInfo = auth
			//settings.TlsConfig = &tls.Config{InsecureSkipVerify: true} Use this only if you're on a Mac running Go 1.18+ doing local dev. See https://github.com/golang/go/issues/51991

		})

	if err != nil {
		fmt.Println(err)
		return
	}
	// Cleanup
	defer driverRemoteConnection.Close()

	// Creating graph traversal
	g := gremlingo.Traversal_().WithRemote(driverRemoteConnection)

	// Add a vertex with properties to the graph with the terminal step Iterate()
	promise := g.AddV("gremlin").Property("language", "go").Iterate()

	// The returned promised is a go channel to wait for all submitted steps to finish execution and return error.
	err = <-promise
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get the value of the property
	result, err := g.V().HasLabel("gremlin").Values("language").ToList()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print the result
	for _, r := range result {
		fmt.Println(r.GetString())
	}
}
