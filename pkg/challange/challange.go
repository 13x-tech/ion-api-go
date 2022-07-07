package challenge

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type Response struct {
	Nonce    string `json:"challengeNonce"`
	Duration int    `json:"validDurationInMinutes"`
	Target   string `json:"largestAllowedHash"`
}

type Challenge struct {
	endpoint string
}

type Options func(*Challenge)

func WithEndpoint(endpoint string) Options {
	return func(c *Challenge) {
		c.endpoint = endpoint
	}
}

func New(opts ...Options) (*Challenge, error) {
	c := new(Challenge)
	for _, opt := range opts {
		opt(c)
	}

	if c.endpoint == "" {
		return nil, fmt.Errorf("invalid endpoint")
	}

	return c, nil
}

func (c Challenge) Get(request []byte) (string, string, error) {

	res, err := http.Get(c.endpoint)
	if err != nil {
		return "", "", fmt.Errorf("could not get challenge: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("could not read response: %w", err)
	}

	var r Response
	if err := json.Unmarshal(body, &r); err != nil {
		return "", "", fmt.Errorf("could not unmarshal response: %w", err)
	}

	t, err := hex.DecodeString(r.Target)
	if err != nil {
		return "", "", fmt.Errorf("could not decode target: %w", err)
	}

	challenge, err := hex.DecodeString(r.Nonce)
	if err != nil {
		return "", "", fmt.Errorf("could not decode nonce: %w", err)
	}

	answer, err := genHash(request, challenge, t, r.Duration)
	if err != nil {
		return "", "", fmt.Errorf("could not get answer: %w", err)
	}

	return r.Nonce, answer, nil
}

func genRandomHexString() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(500))
	if err != nil {
		return "", err
	}

	b := make([]byte, n.Int64())

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	s := hex.EncodeToString([]byte(hex.EncodeToString(b)))
	if len(s) > 1000 {
		s = s[:1000]
	}

	return s, nil
}

func genHash(request, challenge, target []byte, maxDuration int) (string, error) {
	start := time.Now()

	for {
		nonce, err := genRandomHexString()
		if err != nil {
			return "", err
		}

		n, err := hex.DecodeString(nonce)
		if err != nil {
			return "", fmt.Errorf("could not decode string: %w", err)
		}

		password := []byte(fmt.Sprintf("%s%s", n, request))
		work := argon2.IDKey(
			password,
			challenge,
			1,
			1000,
			1,
			32,
		)

		if strings.Compare(hex.EncodeToString(target), hex.EncodeToString(work)) > 0 && time.Since(start) < time.Duration(maxDuration)*time.Minute {
			return nonce, nil
		}

		if time.Since(start) > time.Duration(maxDuration)*time.Minute {
			return "", fmt.Errorf("duration exceeded")
		}
	}
}
