package mkbfv

import (
	"github.com/ldsec/lattigo/v2/bfv"
	"github.com/ldsec/lattigo/v2/ring"
)

// MKCiphertext is type for a bfv ciphertext in a multi key setting
// it contains a bfv ciphertexts along with the list of participants in the right order
type MKCiphertext struct {
	ciphertexts *bfv.Ciphertext
	peerIDs     []uint64
}

// PadCiphers pad two ciphertext corresponding to a different set of parties
// to make their dimension match. ci = 0 if the participant i is not involved in the ciphertext
// the peerIDs list is also updated to match the new dimension
func PadCiphers(c1, c2, c1Out, c2Out *MKCiphertext, params *bfv.Parameters) {

	ringQ := GetRingQ(params)
	allPeers := MergeSlices(c1.peerIDs, c2.peerIDs)
	k := len(allPeers)

	res1 := make([]*ring.Poly, k+1) // + 1 for c0
	res2 := make([]*ring.Poly, k+1)

	// put c0 in
	res1[0] = c1.ciphertexts.Element.Value()[0]
	res2[0] = c2.ciphertexts.Element.Value()[0]

	// copy ciphertext values if participant involved
	// else put a 0 polynomial
	for i, peer := range allPeers {
		index := i + 1

		index1 := Contains(c1.peerIDs, peer)
		index2 := Contains(c2.peerIDs, peer)

		if index1 >= 0 {
			res1[index] = c1.ciphertexts.Element.Value()[index1+1]
		} else {
			res1[index] = ringQ.NewPoly()
		}

		if index2 >= 0 {
			res2[index] = c2.ciphertexts.Element.Value()[index2+1]
		} else {
			res2[index] = ringQ.NewPoly()
		}

	}

	// set ciphertext values
	c1Cipher := new(bfv.Element)
	c2Cipher := new(bfv.Element)

	c1Cipher.SetValue(res1)
	c2Cipher.SetValue(res2)

	c1Out.ciphertexts = c1Cipher.Ciphertext()
	c2Out.ciphertexts = c2Cipher.Ciphertext()

	// set peer Ids
	c1Out.peerIDs = allPeers
	c2Out.peerIDs = allPeers
}
