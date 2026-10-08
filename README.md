# Complyance Go SDK

The Go SDK for the Complyance Unified ingestion API.

## Installation

```sh
go get github.com/complyance-io/complyance-go-sdk/v3@v3.0.5-beta
```

## Raw XML and JSON submission

The example requires `GETS_API_KEY`, `GETS_API_ORIGIN`, `GETS_XML_FILE`, and `GETS_JSON_FILE`. Supply complete Unify request envelopes in the files. Choose exactly one mode on each run: `xml-file`, `xml-string`, `xml-bulk-file`, `xml-bulk-string`, `json`, `json-legacy`, `json-omitted`, or `json-bulk`. Bulk modes require a bulk request fixture.

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    sdk "github.com/complyance-io/complyance-go-sdk/v3/pkg"
)

func main() {
    if len(os.Args) != 2 { panic("Provide one submission mode") }
    key := os.Getenv("GETS_API_KEY")
    origin := os.Getenv("GETS_API_ORIGIN")
    if key == "" || origin == "" { panic("Set GETS_API_KEY and GETS_API_ORIGIN") }
    client := sdk.NewAPIClient(key, sdk.EnvironmentSandbox, sdk.NewNoRetryConfig())
    var result interface{}
    var err error
    switch os.Args[1] {
    case "xml-file":
        result, err = client.SubmitXmlSingleFile(os.Getenv("GETS_XML_FILE"), origin)
    case "xml-bulk-file":
        result, err = client.SubmitXmlBulkFile(os.Getenv("GETS_XML_FILE"), origin)
    case "xml-string", "xml-bulk-string":
        encoded, readErr := os.ReadFile(os.Getenv("GETS_XML_FILE"))
        if readErr != nil { panic(readErr) }
        xml := string(encoded)
        if os.Args[1] == "xml-string" {
            result, err = client.SubmitXmlSingle([]byte(xml), origin)
        } else {
            result, err = client.SubmitXmlBulk([]byte(xml), origin)
        }
    case "json", "json-legacy", "json-omitted", "json-bulk":
        encoded, readErr := os.ReadFile(os.Getenv("GETS_JSON_FILE"))
        if readErr != nil { panic(readErr) }
        switch os.Args[1] {
        case "json": result, err = client.SubmitJSONSingle(encoded, origin, true)
        case "json-legacy": result, err = client.SubmitJSONSingle(encoded, origin, false)
        case "json-omitted": result, err = client.SubmitLegacyJSONSingle(encoded, origin)
        case "json-bulk": result, err = client.SubmitJSONBulk(encoded, origin)
        }
    default: panic("Unknown mode")
    }
    if err != nil { panic(err) }
    encoded, err := json.Marshal(result)
    if err != nil { panic(err) }
    fmt.Println(string(encoded))
}
```

These raw methods post to the supplied HTTP origin at `/api/v3/unify`, preserve the supplied body bytes, and send once without queueing or retrying. XML string inputs use Go's UTF-8 string bytes. JSON single accepts an explicit `new-api` boolean; `SubmitLegacyJSONSingle` omits that header. JSON bulk always selects the revamped contract.

Single revamped JSON/XML responses retain the returned object, including persisted document validation results at HTTP 422. Bulk responses retain summary counts and indexed result/error objects, including validation failures with a returned document ID. HTTP failures are returned as errors.

## Existing API

Existing configured `PushToUnify` methods, policy handling, queues, and retry behavior remain available. Retrieval uses `GetDocumentStatus(documentID)` after configuring the SDK. `GetSubmissionStatus` and `GetStatus` are deprecated and return an error.

## License

MIT. See [LICENSE](LICENSE).
