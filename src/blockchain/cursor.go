package blockchain

// Cursor is a cursor for recursing through the block chain
type Cursor struct {
	Blockchain *Blockchain
	Index int
}

// Cursor returns a cursor for the blockchain
func (bc *Blockchain) Cursor(index int) *Cursor {
	bc.Lock()
	defer bc.Unlock()
	if index > len(bc.Blocks)+1 || index < 0 {
		return nil
	}
	return &Cursor{
		Blockchain: bc,
		Index: index,
	}
}

// Next block in the cursor
func (c *Cursor) Next() (*Block) {
	block := c.Blockchain.Get(c.Index)
	if block == nil {
		return nil
	}
	c.Index++
	return block
}
