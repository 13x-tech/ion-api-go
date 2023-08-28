package tests

import (
	"fmt"
	"testing"

	"github.com/13x-tech/ion-api-go/pkg/api"
	"github.com/stretchr/testify/require"
)

func TestDIDRegistration(t *testing.T) {
	const (
		didAPIEndpoint       = "https://beta.ion.msidentity.com/api/v1.0"
		didChallengeEndpoint = "https://beta.ion.msidentity.com/api/v1.0/proof-of-work-challenge"
		discoverEndpoint     = "https://beta.discover.did.microsoft.com/1.0"
		longformDID          = `did:ion:EiAKkYIJO8KmruJVa0YPXZd9JsapoSDzm0jqSz8XLL88xA:eyJkZWx0YSI6eyJwYXRjaGVzIjpbeyJhY3Rpb24iOiJyZXBsYWNlIiwiZG9jdW1lbnQiOnsicHVibGljS2V5cyI6W3siaWQiOiJrZXktMSIsInB1YmxpY0tleUp3ayI6eyJjcnYiOiJzZWNwMjU2azEiLCJrdHkiOiJFQyIsIngiOiJmZ1ZKWWRQMkdhUnptU3ZpT0FVOHBTdk9uTlZPQk1LQWViWGY1aU5WRlNNIiwieSI6InFZeldLc3VULS1qT01mSmtSWGVDOUp3dkYxVWsxZ2JDZEFQMTlHbDQ2WjQifSwicHVycG9zZXMiOlsiYXV0aGVudGljYXRpb24iXSwidHlwZSI6Ikpzb25XZWJLZXkyMDIwIn1dLCJzZXJ2aWNlcyI6W3siaWQiOiJ6aW9uX2R3biIsInNlcnZpY2VFbmRwb2ludCI6eyJub2RlcyI6WyJodHRwczovL2R3bi56aW9uLmZ5aSJdfSwidHlwZSI6IkRlY2VudHJhbGl6ZWRXZWJOb2RlIn1dfX1dLCJ1cGRhdGVDb21taXRtZW50IjoiRWlEMXVzaWVGMWVmSGRtT1oydERGQmhDb0g2emFfN0ZlbS1Uc3hmWGg4QVRyZyJ9LCJzdWZmaXhEYXRhIjp7ImRlbHRhSGFzaCI6IkVpQXRuMWNidGFfVXRWS3dpcHBNQWJXQUFOOFUxZFh2a1VZS2dPTk9FMmNWcXciLCJyZWNvdmVyeUNvbW1pdG1lbnQiOiJFaUFZdWZrWnRuRE1PdXk4YTRpMDJBdEw0VHI1eDJzV1F0aXc1TGxvd1diOFNnIn19`
	)

	SubmitIONRequest := func(longFormURI string, withChallenge bool) (string, error) {
		var (
			err    error
			ionAPI *api.API
		)
		const didOperationsEndpoint = didAPIEndpoint + "/operations"

		suffixData, delta, err := api.ParseLongForm(longFormURI)
		if err != nil {
			return "", fmt.Errorf("invalid long form uri: %w", err)
		}

		opts := []api.Options{
			api.WithEndpoint(didOperationsEndpoint),
		}

		if withChallenge {
			opts = append(opts, api.WithChallenge(didChallengeEndpoint))
		}

		ionAPI, err = api.New(opts...)
		if err != nil {
			return "", fmt.Errorf("could not create api: %w", err)
		}
		response, err := ionAPI.Submit(
			api.CreateOperation(suffixData, delta),
		)
		if err != nil {
			return "", fmt.Errorf("could not submit request: %w", err)
		}
		return string(response), nil
	}

	t.Run("test simple registration", func(tt *testing.T) {
		longFormURI := `did:ion:EiAKkYIJO8KmruJVa0YPXZd9JsapoSDzm0jqSz8XLL88xA:eyJkZWx0YSI6eyJwYXRjaGVzIjpbeyJhY3Rpb24iOiJyZXBsYWNlIiwiZG9jdW1lbnQiOnsicHVibGljS2V5cyI6W3siaWQiOiJrZXktMSIsInB1YmxpY0tleUp3ayI6eyJjcnYiOiJzZWNwMjU2azEiLCJrdHkiOiJFQyIsIngiOiJmZ1ZKWWRQMkdhUnptU3ZpT0FVOHBTdk9uTlZPQk1LQWViWGY1aU5WRlNNIiwieSI6InFZeldLc3VULS1qT01mSmtSWGVDOUp3dkYxVWsxZ2JDZEFQMTlHbDQ2WjQifSwicHVycG9zZXMiOlsiYXV0aGVudGljYXRpb24iXSwidHlwZSI6Ikpzb25XZWJLZXkyMDIwIn1dLCJzZXJ2aWNlcyI6W3siaWQiOiJ6aW9uX2R3biIsInNlcnZpY2VFbmRwb2ludCI6eyJub2RlcyI6WyJodHRwczovL2R3bi56aW9uLmZ5aSJdfSwidHlwZSI6IkRlY2VudHJhbGl6ZWRXZWJOb2RlIn1dfX1dLCJ1cGRhdGVDb21taXRtZW50IjoiRWlEMXVzaWVGMWVmSGRtT1oydERGQmhDb0g2emFfN0ZlbS1Uc3hmWGg4QVRyZyJ9LCJzdWZmaXhEYXRhIjp7ImRlbHRhSGFzaCI6IkVpQXRuMWNidGFfVXRWS3dpcHBNQWJXQUFOOFUxZFh2a1VZS2dPTk9FMmNWcXciLCJyZWNvdmVyeUNvbW1pdG1lbnQiOiJFaUFZdWZrWnRuRE1PdXk4YTRpMDJBdEw0VHI1eDJzV1F0aXc1TGxvd1diOFNnIn19`
		response, err := SubmitIONRequest(longFormURI, false)
		require.NoError(tt, err)
		fmt.Printf("Response: %s\n", response)
	})

	t.Run("test inbuilt registration with challenge", func(tt *testing.T) {
		longFormURI := `did:ion:EiAKkYIJO8KmruJVa0YPXZd9JsapoSDzm0jqSz8XLL88xA:eyJkZWx0YSI6eyJwYXRjaGVzIjpbeyJhY3Rpb24iOiJyZXBsYWNlIiwiZG9jdW1lbnQiOnsicHVibGljS2V5cyI6W3siaWQiOiJrZXktMSIsInB1YmxpY0tleUp3ayI6eyJjcnYiOiJzZWNwMjU2azEiLCJrdHkiOiJFQyIsIngiOiJmZ1ZKWWRQMkdhUnptU3ZpT0FVOHBTdk9uTlZPQk1LQWViWGY1aU5WRlNNIiwieSI6InFZeldLc3VULS1qT01mSmtSWGVDOUp3dkYxVWsxZ2JDZEFQMTlHbDQ2WjQifSwicHVycG9zZXMiOlsiYXV0aGVudGljYXRpb24iXSwidHlwZSI6Ikpzb25XZWJLZXkyMDIwIn1dLCJzZXJ2aWNlcyI6W3siaWQiOiJ6aW9uX2R3biIsInNlcnZpY2VFbmRwb2ludCI6eyJub2RlcyI6WyJodHRwczovL2R3bi56aW9uLmZ5aSJdfSwidHlwZSI6IkRlY2VudHJhbGl6ZWRXZWJOb2RlIn1dfX1dLCJ1cGRhdGVDb21taXRtZW50IjoiRWlEMXVzaWVGMWVmSGRtT1oydERGQmhDb0g2emFfN0ZlbS1Uc3hmWGg4QVRyZyJ9LCJzdWZmaXhEYXRhIjp7ImRlbHRhSGFzaCI6IkVpQXRuMWNidGFfVXRWS3dpcHBNQWJXQUFOOFUxZFh2a1VZS2dPTk9FMmNWcXciLCJyZWNvdmVyeUNvbW1pdG1lbnQiOiJFaUFZdWZrWnRuRE1PdXk4YTRpMDJBdEw0VHI1eDJzV1F0aXc1TGxvd1diOFNnIn19`
		response, err := SubmitIONRequest(longFormURI, true)
		require.NoError(tt, err)
		fmt.Printf("Response: %s\n", response)
	})
}
