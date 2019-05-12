httpclient
===
[![GoDoc](https://godoc.org/github.com/x-mod/httpclient?status.svg)](https://godoc.org/github.com/x-mod/httpclient) [![Go Report Card](https://goreportcard.com/badge/github.com/x-mod/httpclient)](https://goreportcard.com/report/github.com/x-mod/httpclient) 

More smooth package for http operations as a client:

- http request builder
- http response processor
- http client extension

### http.Request Builder

````go

import "github.com/x-mod/httpclient"

requestBuilder := httpclient.NewRequestBuilder(
        httpclient.URL("https://url"),
        httpclient.Method("GET"),
        httpclient.Query("key", "value"),
        httpclient.Header("key", "value"),
        httpclient.BasicAuth("user", "pass"),
        httpclient.Credential(*tlsconfig),
        httpclient.Body(
            httpclient.JSON(map[string]interface{}{
                "a": "hello",
                "b": true,
                "c": 1,
            }),
        ),
    )

req, err := requestBuilder.Get()

````

### http.Response Processor

````go
//ResponseProcessor interface
type ResponseProcessor interface {
	Process(context.Context, *http.Response) error
}
````

Implement your own ResponseProcessor

### http.Client Extension

extend the http.Client with more useful interfaces:

````go
import "github.com/x-mod/httpclient"

client := httpclient.New(
    httpclient.MaxConnsPerHost(16),
    httpclient.Retry(3),
    httpclient.Response(
        httpclient.NewDumpResponse(),
    ),
)

//get standard http.Client
c := client.GetClient()
//get standard http.Transport
tr := client.GetTransport()

//extension fn
err := client.Execute(context.TODO())
err := client.ExecuteRequest(context.TODO(), request)

````