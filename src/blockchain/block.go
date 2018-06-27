package blockchain

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// Hash represents a hash
type Hash [64]byte

// MarshalJSON for hashes
func (h *Hash) MarshalJSON() ([]byte,error) {
	hash := make([]byte, hex.EncodedLen(len(h[:])))
	hex.Encode(hash, h[:])
	return []byte(fmt.Sprintf("\"%s\"", hash)), nil
}

// UnmarshalJSON for hashes
func (h *Hash) UnmarshalJSON(buf []byte) error {
	val := buf[1:len(buf)-1]
	if hex.DecodedLen(len(val)) > 64 {
		return fmt.Errorf("long buffer: %d/%d", len(val), hex.DecodedLen(len(val)))
	}
	hex.Decode(h[:], val)
	return nil
}

// Inc (rement) a Hash's value
func (h *Hash) Inc() {
	for i := 0; i < 64; i++ {
		if h[i] <= 0xfe {
			h[i]++
			break
		}
		h[i] = 0x00
	}
}

// Block in the block chain
type Block struct {
	Hash                 Hash
	Index                uint64
	Timestamp            uint64
	Difficulty           uint64
	PrevHash             Hash
	Transactions         []Tx
	Nonce                Hash
}

// GenerateHash generates the block hash
func (block *Block) GenerateHash() [64]byte {
	var header [(8*3)]byte
	var ret [64]byte
	h := sha3.New512()

	binary.LittleEndian.PutUint64(header[0:], uint64(block.Index))
	binary.LittleEndian.PutUint64(header[0:], uint64(block.Timestamp))
	binary.LittleEndian.PutUint64(header[16:], uint64(block.Difficulty))

	h.Write(header[:])
	h.Write(block.PrevHash[:])
	h.Write(block.Nonce[:])
	hashed := h.Sum(nil)
	copy(ret[:], hashed[:])
	return ret
}

// Generate a block
func (block *Block) Generate(timestamp, difficulty uint64, nonce [64]byte) (*Block, error) {

	var newBlock Block

	newBlock.Index = block.Index + 1
	newBlock.Timestamp = timestamp
	newBlock.PrevHash = block.Hash
	newBlock.Nonce = nonce
	newBlock.Difficulty = difficulty
	newBlock.Hash = newBlock.GenerateHash()

	return &newBlock, nil
}

// IsValid block
func (block *Block) IsValid(newBlock *Block) error {
	if block.Index+1 != newBlock.Index {
		return fmt.Errorf("non sequential block index")
	}
	if block.Timestamp > newBlock.Timestamp {
		return fmt.Errorf("bad timestamp")
	}
	if block.Hash != newBlock.PrevHash {
		return fmt.Errorf("bad previous hash")
	}
	if newBlock.GenerateHash() != newBlock.Hash {
		return fmt.Errorf("bad hash")
	}
	return nil
}

// CountDiff for a block
func (h *Hash) CountDiff() uint64 {
	var d uint64
	for i := 0; i < 64; i++ {
		for b := uint8(0); b < 8; b++ {
			mask := uint8(1) << (7-b)
			bits := uint8(h[i])
			//fmt.Printf("%08b & %08b = %08b shift=%d (%d) d=%d\n",
			//	bits, mask, bits & mask, (7-b), b, d)
			if (bits & mask) != 0 {
				return d
			}
			d++
		}
	}
	return d
}

// HasDifficulty checks the difficulty for adding the block
func (block *Block) HasDifficulty(difficulty uint64) bool {
	d := block.Hash.CountDiff()
	return d >= difficulty
}
