# ION Challenge and API Package in Go

##### Challenge Package
Microsoft has a Proof-of-Work / HashCash style spam prevention on their open API for interfacing with their ION Node.
The challenge package adheres to that API and handles the proof of work

##### API Package
This helps format and submit operations to the ION API, it accepts an optional challenge.

Example Usage:
```go
package main

import (
	"fmt"

	"github.com/13x-tech/ion-api-go/pkg/api"
)

func main() {
	longFormURI := "did:ion:test:EiAnKD8-jfdd0MDcZUjAbRgaThBrMxPTFOxcnfJhI7Ukaw:eyJkZWx0YSI6eyJwYXRjaGVzIjpbeyJhY3Rpb24iOiJyZXBsYWNlIiwiZG9jdW1lbnQiOnsicHVibGljS2V5cyI6W3siaWQiOiJzaWdfNzJiZDE2ZDYiLCJwdWJsaWNLZXlKd2siOnsiY3J2Ijoic2VjcDI1NmsxIiwia3R5IjoiRUMiLCJ4IjoiS2JfMnVOR3Nyd1VOdkh2YUNOckRGdW14VXlQTWZZd3kxNEpZZmphQUhmayIsInkiOiJhSFNDZDVEOFh0RUxvSXBpN1A5eDV1cXBpeEVxNmJDenQ0QldvUVk1UUFRIn0sInB1cnBvc2VzIjpbImF1dGhlbnRpY2F0aW9uIiwiYXNzZXJ0aW9uTWV0aG9kIl0sInR5cGUiOiJFY2RzYVNlY3AyNTZrMVZlcmlmaWNhdGlvbktleTIwMTkifV0sInNlcnZpY2VzIjpbeyJpZCI6ImxpbmtlZGRvbWFpbnMiLCJzZXJ2aWNlRW5kcG9pbnQiOnsib3JpZ2lucyI6WyJodHRwczovL3d3dy52Y3NhdG9zaGkuY29tLyJdfSwidHlwZSI6IkxpbmtlZERvbWFpbnMifV19fV0sInVwZGF0ZUNvbW1pdG1lbnQiOiJFaUR4SWxJak9xQk5NTGZjdzZndWpHNEdFVDM3UjBIRWM2Z20xclNZTjlMOF9RIn0sInN1ZmZpeERhdGEiOnsiZGVsdGFIYXNoIjoiRWlBLXV3TWo3RVFheURmWTRJS3pfSE9LdmJZQ05td19Tb1lhUmhOcWhFSWhudyIsInJlY292ZXJ5Q29tbWl0bWVudCI6IkVpQ0czQ1M5RFJpeU1JRVoxRl9sSjZnRVRMZWVHREwzZnpuQUViMVRGdFZXNEEifX0#sig_72bd16d6"

	response, err := SubmitIONRequest(longFormURI)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", response)
}

func SubmitIONRequest(longFormURI string) (string, error) {

	suffixData, delta, err := api.ParseLongForm(longFormURI)
	if err != nil {
		return "", fmt.Errorf("invalid long form uri: %w", err)
	}

	ionAPI, err := api.New(
		ion.WithEndpoint("https://{ion-node-url}/operations"),
	)
	if err != nil {
		return "", fmt.Errorf("could not create api: %w", err)
	}

	response, err := ionAPI.Submit(
		ion.CreateOperation(suffixData, delta),
	)
	if err != nil {
		return "", fmt.Errorf("could not submit request: %w", err)
	}

	return string(response), nil
}
```


With Challenge Usage:
```go
	api.New(
		ion.WithEndpoint("https://{ion-node-url}/operations"),
    ion.WithChallengeEndpoint("https://{ion-node-url}/proof-of-work-challenge"),
	)
```