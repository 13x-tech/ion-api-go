package api

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

type Delta struct {
	Patches          []map[string]interface{} `json:"patches"`
	UpdateCommitment string                   `json:"updateCommitment"`
}
