package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/13x-tech/ion-api-go/pkg/challange"
)

func ParseLongForm(uri string) (SuffixData, Delta, error) {
	splitURI := strings.Split(uri, ":")
	if len(splitURI) < 2 {
		return SuffixData{}, Delta{}, fmt.Errorf("invalid long form uri: %s", uri)
	}

	splitLength := len(splitURI)

	didSuffix := splitURI[splitLength-2]
	longFormData := splitURI[splitLength-1]

	longFormDataBytes, err := base64.RawURLEncoding.DecodeString(longFormData)
	if err != nil {
		return SuffixData{}, Delta{}, fmt.Errorf("failed to base64 decode suffix data: %w", err)
	}

	var longFormDataStruct struct {
		SuffixData SuffixData `json:"suffixData"`
		Delta      Delta      `json:"delta"`
	}

	if err := json.Unmarshal(longFormDataBytes, &longFormDataStruct); err != nil {
		return SuffixData{}, Delta{}, fmt.Errorf("failed to unmarshal suffix data: %w", err)
	}

	testSuffix, err := longFormDataStruct.SuffixData.URI()
	if err != nil {
		return SuffixData{}, Delta{}, fmt.Errorf("failed to create suffix uri: %w", err)
	}

	if testSuffix != didSuffix {
		return SuffixData{}, Delta{}, fmt.Errorf("suffix uri does not match: %s != %s", testSuffix, didSuffix)
	}

	return longFormDataStruct.SuffixData, longFormDataStruct.Delta, nil
}

type API struct {
	endpoint  string
	challange *challange.Challange
}

type Options func(a *API) error

func WithEndpoint(endpoint string) Options {
	return func(a *API) error {
		a.endpoint = endpoint
		return nil
	}
}

func WithChallange(endpoint string) Options {
	return func(a *API) error {
		ch, err := challange.New(
			challange.WithEndpoint(endpoint),
		)
		if err != nil {
			return fmt.Errorf("challange error: %w", err)
		}
		a.challange = ch
		return nil
	}
}

func New(opts ...Options) (*API, error) {
	a := new(API)
	for _, opt := range opts {
		opt(a)
	}
	if a.endpoint == "" {
		return nil, fmt.Errorf("invalid endpoint")
	}

	return a, nil
}

func (a *API) Submit(i interface{}) ([]byte, error) {

	var requestJSON []byte
	var err error

	switch v := i.(type) {
	case Create:
		requestJSON, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("could not marshal create request: %w", err)
		}
	case Update:
		requestJSON, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("could not marshal update request: %w", err)
		}
	case Deactivate:
		requestJSON, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("could not marshal deactivate request: %w", err)
		}
	case Recover:
		requestJSON, err = json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("could not marshal recover request: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid interface")
	}

	req, err := http.NewRequest("POST", a.endpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	if a.challange != nil {
		nonce, answer, err := a.challange.Get(requestJSON)
		if err != nil {
			return nil, fmt.Errorf("could not fetch challange: %w", err)
		}
		req.Header.Add("Challenge-Nonce", nonce)
		req.Header.Add("Answer-Nonce", answer)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %w", err)
	}

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %w", err)
	}

	return response, nil
}
