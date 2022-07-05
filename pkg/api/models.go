package api

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/gowebpki/jcs"
	mh "github.com/multiformats/go-multihash"
)

type Operation struct {
	Type string `json:"type"`
}

type UpdateOp struct {
	DIDSuffix   string `json:"didSuffix"`
	RevealValue string `json:"revealValue"`
	SignedData  string `json:"signedData"`
}

func CreateOperation(suffixData SuffixData, delta Delta) Create {
	return Create{
		Operation: Operation{
			Type: "create",
		},
		SuffixData: suffixData,
		Delta:      delta,
	}
}

type Create struct {
	Operation
	SuffixData SuffixData `json:"suffixData"`
	Delta      Delta      `json:"delta"`
}

func (c Create) ShortFormURI(method string) (string, error) {
	suffix, err := c.SuffixData.URI()
	if err != nil {
		return "", fmt.Errorf("could not get short form uri: %w", err)
	}

	return fmt.Sprintf("did:%s:%s", method, suffix), nil
}

func (c Create) LongFormURI(method string) (string, error) {
	suffix, err := c.SuffixData.URI()
	if err != nil {
		return "", fmt.Errorf("could not get short form uri: %w", err)
	}
	encodedSuffixData, err := c.EncodedSuffixData()
	if err != nil {
		return "", fmt.Errorf("could not get encoded suffix data: %w", err)
	}

	return fmt.Sprintf("did:%s:%s:%s", method, suffix, encodedSuffixData), nil
}

func (c Create) EncodedSuffixData() (string, error) {

	marshalStruct := struct {
		Delta      Delta      `json:"delta"`
		SuffixData SuffixData `json:"suffixData"`
	}{
		Delta:      c.Delta,
		SuffixData: c.SuffixData,
	}

	didData, err := json.Marshal(marshalStruct)
	if err != nil {
		return "", fmt.Errorf("failed to marshal DID: %w", err)
	}

	jsonData, err := jcs.Transform(didData)
	if err != nil {
		return "", fmt.Errorf("failed to transform DID: %w", err)
	}

	encodedSuffixData := base64.RawURLEncoding.EncodeToString(jsonData)

	return encodedSuffixData, nil
}

func UpdateOperation(suffix, reveal, signature string, delta Delta) Update {
	return Update{
		Operation: Operation{
			Type: "update",
		},
		UpdateOp: UpdateOp{
			DIDSuffix:   suffix,
			RevealValue: reveal,
			SignedData:  signature,
		},
		Delta: delta,
	}
}

type Update struct {
	Operation
	UpdateOp
	Delta Delta `json:"delta"`
}

func RecoverOperation(suffix, reveal, signature string, delta Delta) Recover {
	return Recover{
		Operation: Operation{
			Type: "recover",
		},
		UpdateOp: UpdateOp{
			DIDSuffix:   suffix,
			RevealValue: reveal,
			SignedData:  signature,
		},
		Delta: delta,
	}
}

type Recover Update

func DeactivateOperation(suffix, reveal, signature string) Deactivate {
	return Deactivate{
		Operation: Operation{
			Type: "deactivate",
		},
		UpdateOp: UpdateOp{
			DIDSuffix:   suffix,
			RevealValue: reveal,
			SignedData:  signature,
		},
	}
}

type Deactivate struct {
	Operation
	UpdateOp
}

type SuffixData struct {
	Type               string `json:"type,omitempty"`
	DeltaHash          string `json:"deltaHash"`
	RecoveryCommitment string `json:"recoveryCommitment"`
	AnchorOrigin       string `json:"anchorOrigin,omitempty"`
}

func (s SuffixData) URI() (string, error) {
	// Short Form DID URI
	// https://identity.foundation/sidetree/spec/#short-form-did

	bytes, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal suffix data: %w", err)
	}

	jcsBytes, err := jcs.Transform(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to transform bytes: %w", err)
	}

	h256 := sha256.Sum256(jcsBytes)
	hash, err := mh.Encode(h256[:], mh.SHA2_256)
	if err != nil {
		return "", fmt.Errorf("failed to create hash: %w", err)
	}
	encoder := base64.RawURLEncoding
	return encoder.EncodeToString(hash), nil
}

type Delta struct {
	Patches          []map[string]interface{} `json:"patches"`
	UpdateCommitment string                   `json:"updateCommitment"`
}
