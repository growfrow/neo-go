package transaction

import (
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/CityOfZion/neo-go/pkg/crypto/hash"
	"github.com/CityOfZion/neo-go/pkg/util"
)

// Witness contains 2 scripts.
type Witness struct {
	InvocationScript   []byte
	VerificationScript []byte
}

// DecodeBinary implements the payload interface.
func (w *Witness) DecodeBinary(r io.Reader) error {
	br := util.BinReader{R: r}

	w.InvocationScript = br.ReadBytes()
	w.VerificationScript = br.ReadBytes()
	return br.Err
}

// EncodeBinary implements the payload interface.
func (w *Witness) EncodeBinary(writer io.Writer) error {
	bw := util.BinWriter{W: writer}

	bw.WriteBytes(w.InvocationScript)
	bw.WriteBytes(w.VerificationScript)

	return bw.Err
}

// MarshalJSON implements the json marshaller interface.
func (w *Witness) MarshalJSON() ([]byte, error) {
	data := map[string]string{
		"invocation":   hex.EncodeToString(w.InvocationScript),
		"verification": hex.EncodeToString(w.VerificationScript),
	}

	return json.Marshal(data)
}

// Size returns the size in bytes of the Witness.
func (w *Witness) Size() int {
	return util.GetVarSize(w.InvocationScript) + util.GetVarSize(w.VerificationScript)
}

// ScriptHash returns the hash of the VerificationScript.
func (w Witness) ScriptHash() util.Uint160 {
	return hash.Hash160(w.VerificationScript)
}
