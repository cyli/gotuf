package data

import (
	"encoding/json"
	"time"

	cjson "github.com/tent/canonical-json-go"
)

type SignedSnapshot struct {
	Signatures []Signature
	Signed     Snapshot
	Dirty      bool
}

type Snapshot struct {
	Type    string    `json:"_type"`
	Version int       `json:"version"`
	Expires time.Time `json:"expires"`
	Meta    Files     `json:"meta"`
}

func NewSnapshot() (*SignedSnapshot, error) {
	return &SignedSnapshot{
		Signatures: make([]Signature, 0),
		Signed: Snapshot{
			Type:    TUFTypes["snapshot"],
			Version: 0,
			Expires: DefaultExpires("snapshot"),
			Meta:    make(Files),
		},
	}, nil
}

func (sp *SignedSnapshot) hashForRole(role string) []byte {
	return sp.Signed.Meta[role].Hashes["sha256"]
}

func (sp SignedSnapshot) ToSigned() (*Signed, error) {
	s, err := cjson.Marshal(sp.Signed)
	if err != nil {
		return nil, err
	}
	signed := json.RawMessage{}
	err = signed.UnmarshalJSON(s)
	if err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(sp.Signatures))
	copy(sigs, sp.Signatures)
	return &Signed{
		Signatures: sigs,
		Signed:     signed,
	}, nil
}

func (sp *SignedSnapshot) AddMeta(role string, meta FileMeta) {
	sp.Signed.Meta[role] = meta
	sp.Dirty = true
}

func SnapshotFromSigned(s *Signed) (*SignedSnapshot, error) {
	sp := Snapshot{}
	err := json.Unmarshal(s.Signed, &sp)
	if err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(s.Signatures))
	copy(sigs, s.Signatures)
	return &SignedSnapshot{
		Signatures: sigs,
		Signed:     sp,
	}, nil
}
