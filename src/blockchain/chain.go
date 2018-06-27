package blockchain

import (
	"fmt"
	"sync"
)

// Blockchain type
type Blockchain struct {
	sync.Mutex
	Blocks []*Block
}

// Get a block from the block chain
func (bc *Blockchain) get(index int) *Block {
	if index == -1 {
		index = len(bc.Blocks) - 1
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
	return bc.calculateDifficulty()
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
	if !block.HasDifficulty(bc.calculateDifficulty()) {
		return fmt.Errorf("block has low difficulty: %d < %d",
			block.Difficulty, bc.calculateDifficulty())
	}
	return nil
}

func (bc *Blockchain) calculateDifficulty() uint64 {
	fmt.Printf("Calculate difficulty\n")
	defer func() {
		fmt.Printf("Calculated!\n")
	}()
	last := bc.get(-1)
	if last == nil {
		return 0
	}
	fmt.Printf("last[%d]=%d\n", last.Index, last.Difficulty)
	prev := bc.get(int(last.Index - 1))
	if prev == nil {
		return 0
	}
	blocktime := last.Timestamp - prev.Timestamp
	fmt.Printf("prev[%d]=%d blocktime=%d\n", prev.Index,
		prev.Difficulty, blocktime)
	if (blocktime) <= 90 {
		return prev.Difficulty + 1
	} else if (blocktime) > 120 && prev.Difficulty > 0 {
		return prev.Difficulty - 1
	}
	return prev.Difficulty
}

// Add a block to the blockchain
func (bc *Blockchain) Add(block *Block) error {
	bc.Lock()
	defer bc.Unlock()
	if err := bc.isValid(block); err != nil {
		return err
	}
	bc.Blocks = append(bc.Blocks, block)
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
	return nil
}

// GetCumulativeDifficulty gets the cumulative difficulty of a block by index
func (bc *Blockchain) GetCumulativeDifficulty(index int) (int, error) {
	bc.Lock()
	defer bc.Unlock()
	return bc.getCumulativeDifficulty(index)
}

func (bc *Blockchain) getCumulativeDifficulty(index int) (int, error) {
	cum := 0
	for i := 0; i <= index; i++ {
		block := bc.get(i)
		if block == nil {
			return -1, fmt.Errorf("no such block: %d", i)
		}
		cum += int(block.Difficulty)
	}
	return cum, nil
}

// Replace part of the blockchain
func (bc *Blockchain) Replace(bcn *Fragment) error {
	bc.Lock()
	defer bc.Unlock()
	if err := bcn.IsValid(); err != nil {
		return err
	}
	cum, err := bcn.calculateCumulativeDifficulty(bc)
	if err != nil {
		return err
	}
	last := bc.get(-1)
	ccum, err := bc.getCumulativeDifficulty(int(last.Index))
	if err != nil {
		return err
	}
	if last == nil || int(cum) > ccum {
		newblocks := bc.Blocks[:bcn.Start]
		for i := bcn.Start; i <= bcn.End; i++ {
			newblocks = append(newblocks, bcn.Blocks[i])
		}
		bc.Blocks = newblocks
		fmt.Println("replaced chain")
		return nil
	}
	return fmt.Errorf("invalid chain fragment")
}

// Fragment holds a fragment of a blockchain
type Fragment struct {
	Start  uint64
	End    uint64
	Blocks map[uint64]*Block
}

// GetFragment from the Blockchain
func (bc *Blockchain) GetFragment(start, end uint64) *Fragment {
	fragment := Fragment{
		Start:  start,
		End:    end,
		Blocks: map[uint64]*Block{},
	}
	for i := fragment.Start; i <= fragment.End; i++ {
		fragment.Blocks[i] = bc.Get(int(i))
	}
	return &fragment
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
		block, exists := f.Blocks[i]
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
func (f *Fragment) calculateCumulativeDifficulty(bc *Blockchain) (uint64, error) {
	prev := bc.get(int(f.Start) - 1)
	if prev == nil {
		return 0, fmt.Errorf("chain is from the future: %d",
			int(f.Start)-1)
	}
	cumulativeDifficulty, err := bc.getCumulativeDifficulty(int(f.Start) - 1)
	if err != nil {
		return 0, err
	}
	for i := f.Start; i <= f.End; i++ {
		block, exists := f.Blocks[i]
		if !exists {
			return 0, fmt.Errorf("missing block")
		}
		cumulativeDifficulty += int(block.Difficulty)
		prev = block
	}
	return uint64(cumulativeDifficulty), nil
}

// Treenode represents a block in a tree in the simplest form possible
type Treenode struct {
	Index    uint64
	Hash     Hash
	PrevHash Hash
}

// Tree holds the utmost basic data to match a blockchain
type Tree []*Treenode

// IsValid does a simple Tree check
func (t *Tree) IsValid() bool {
	return true
}
