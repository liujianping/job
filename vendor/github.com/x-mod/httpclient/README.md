httpclient
===

smooth package for http client:

- http request builder
- http response processor
- http transport provider
- http client invoker


````go

import "github.com/x-mod/httpclient"

client := httpclient.New(
    //Request Builder
    httpclient.Request(
        httpclient.NewRequestBuilder(
            httpclient.URL("https://url"),
            httpclient.Method("GET"),
        ),
    ),
    //Response Processor
    httpclient.Response(
        httpclient.NewDumpResponse(),
    ),
    //Client Transport
    httpclient.Transport(
        //...
    ),
)
//Do 
rsp, err := client.Do(context.TODO())

//or Execute http request with response processor
if err := client.Execute(context.TODO()); err != nil {
    ...
}

````