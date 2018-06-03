package blockchain

import (
	"fmt"
	"math/big"
	"sync"
)

// Blockchain type
type Blockchain struct {
	sync.Mutex
	Difficulty uint64
	CumulativeDifficulty *big.Int `json:"-"`
	Blocks []*Block
}

// Get a block from the block chain
func (bc *Blockchain) get(index int) *Block {
	if index == -1 {
		index = len(bc.Blocks)-1
	}
	if index >= len(bc.Blocks) || index < 0 {
		return nil
	}
	return bc.Blocks[index]
}

// Get a block from the block chain
func (bc *Blockchain) Get(index int) *Block {
	bc.Lock()
	defer bc.Unlock()
	return bc.get(index)
}

// GetAll blocks from the block chain
func (bc *Blockchain) GetAll() []*Block {
	bc.Lock()
	defer bc.Unlock()
	blocks := make([]*Block, len(bc.Blocks))
	for idx, block := range bc.Blocks {
		blocks[idx] = block
	}
	return blocks
}

func min(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

// GetDifficulty for the block chain
func (bc *Blockchain) GetDifficulty() uint64 {
	bc.Lock()
	defer bc.Unlock()
	return bc.Difficulty
}

// IsValid block to add to the chain
func (bc *Blockchain) isValid(block *Block) error {
	//fmt.Printf("DIFF: %d >= %d\n", block.Nonce.CountDiff(), bc.Difficulty)
	last := bc.get(-1)
	if last == nil {
		return fmt.Errorf("no previous block")
	}
	if err := last.IsValid(block); err != nil {
		return err
	}
	if !block.HasDifficulty(bc.Difficulty) {
		return fmt.Errorf("block has low difficulty")
	}
	return nil
}

func (bc *Blockchain) calculateDifficulty() uint64 {
	last := bc.get(-1)
	if last == nil {
		return 0
	}
	prev := bc.get(int(last.Index-1))
	if (last.Timestamp - prev.Timestamp) <= 90 {
		return bc.Difficulty+1
	} else if (last.Timestamp - prev.Timestamp) > 120 && bc.Difficulty > 0 {
		return bc.Difficulty-1
	}
	return bc.Difficulty
}

// Add a block to the blockchain
func (bc *Blockchain) Add(block *Block) error {
	bc.Lock()
	defer bc.Unlock()
	if err := bc.isValid(block); err != nil {
		return err
	}
	bc.Blocks = append(bc.Blocks, block)
	block.CumulativeDifficulty = big.NewInt(0)
	for _, blk := range bc.Blocks {
		block.CumulativeDifficulty.Add(
			block.CumulativeDifficulty,
			big.NewInt(int64(blk.Difficulty)),
		)
	}
	/*
	block.CumulativeDifficulty.Add(
		bc.CumulativeDifficulty,
		big.NewInt(int64(bc.Difficulty)),
	)
	*/
	bc.CumulativeDifficulty.Add(
		bc.CumulativeDifficulty,
		big.NewInt(int64(bc.Difficulty)),
	)
	bc.Difficulty = bc.calculateDifficulty()
	return nil
}

// Genesis add thes genesis block to a new Blockchain
func (bc *Blockchain) Genesis(block *Block) error {
	bc.Lock()
	defer bc.Unlock()
	if len(bc.Blocks) > 0 {
		return fmt.Errorf("genesis block already exists")
	}
	bc.Blocks = append(bc.Blocks, block)
	bc.CumulativeDifficulty = big.NewInt(0)
	block.CumulativeDifficulty = big.NewInt(0)
	return nil
}

// Replace part of the blockchain
func (bc *Blockchain) Replace(bcn *Fragment) error {
	bc.Lock()
	defer bc.Unlock()
	if err := bcn.IsValid(); err != nil {
		return err
	}
	if cum, err := bcn.calculateCumulativeDifficulty(bc); err != nil {
		return err
	} else if cum == nil {
		return fmt.Errorf("no cumulative difficulty")
	} else {
		if cum.Cmp(bc.CumulativeDifficulty) == 1 {
			newblocks := bc.Blocks[:bcn.Start]
			for i := bcn.Start; i <= bcn.End; i++ {
				newblocks = append(newblocks, bcn.Blocks[i])
			}
			bc.Blocks = newblocks
			bc.Difficulty = bc.get(-1).Difficulty
			bc.Difficulty = bc.calculateDifficulty()
			bc.CumulativeDifficulty =
				big.NewInt(0).Add(cum, big.NewInt(0))
			return nil
		}
	}
	return fmt.Errorf("invalid chain fragment")
}

// Fragment holds a fragment of a blockchain
type Fragment struct {
	CumulativeDifficulty *big.Int `json:"-"`
	Start uint64
	End uint64
	Blocks map[uint64]*Block
}

// IsValid fragment?
func (f *Fragment) IsValid() error {
	if f.Start > f.End {
		return fmt.Errorf("ouroborous error, smeg head")
	}
	if f.Start <= 0 {
		return fmt.Errorf("do not replace genesys")
	}
	for i := f.Start; i <= f.End; i++ {
		block, exists := f.Blocks[i];
		if !exists {
			return fmt.Errorf("missing block: %d", i)
		}
		if i > f.Start && i < f.End-1 {
			next, exists := f.Blocks[i+1]
			if !exists {
				return fmt.Errorf("missing next block: %d", i)
			}
			if next.PrevHash != block.Hash {
				return fmt.Errorf("broken chain, no one to calls us smeg head")
			}
			if next.Timestamp < block.Timestamp {
				return fmt.Errorf(
					"coincidentality broken, make your lunch before you eat it: %d < %d",
					next.Timestamp, block.Timestamp)
			}
		}
	}
	return nil
}

// CalculateCumulativeDifficulty for this Fragment
func (f *Fragment) calculateCumulativeDifficulty(bc *Blockchain) (*big.Int, error) {
	prev := bc.get(int(f.Start)-1)
	if prev == nil {
		return nil, fmt.Errorf("chain is from the future: %d",
			int(f.Start)-1)
	}
	if prev.CumulativeDifficulty == nil {
		return nil, fmt.Errorf("no cumulative from the start")
	}
	for i := f.Start; i <= f.End; i++ {
		block, exists := f.Blocks[i];
		if !exists {
			return nil, fmt.Errorf("missing block")
		}
		if prev.CumulativeDifficulty == nil {
			return nil, fmt.Errorf("no cumulative")
		}
		block.CumulativeDifficulty = big.NewInt(0).Add(
			prev.CumulativeDifficulty,
			big.NewInt(int64(block.Difficulty)),
		)
		prev = block
	}
	return f.Blocks[f.End].CumulativeDifficulty, nil
}

// Treenode represents a block in a tree in the simplest form possible
type Treenode struct {
	Index     uint64
	Hash      Hash
	PrevHash  Hash
}

// Tree holds the utmost basic data to match a blockchain
type Tree []*Treenode

// IsValid does a simple Tree check
func (t *Tree) IsValid() bool {
	return true
}
