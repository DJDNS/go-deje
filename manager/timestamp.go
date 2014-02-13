package manager

import (
	"github.com/campadrenalin/go-deje/model"
	"sort"
	"strconv"
)

type TimestampManager struct {
	ObjectManager
}

func NewTimestampManager() TimestampManager {
	om := NewObjectManager()
	return TimestampManager{om}
}

func (tm *TimestampManager) Register(timestamp model.Timestamp) {
	tm.register(timestamp)
}

func (tm *TimestampManager) Unregister(timestamp model.Timestamp) {
	tm.unregister(timestamp)
}

type Uint64Slice []uint64

func (s Uint64Slice) Len() int           { return len(s) }
func (s Uint64Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s Uint64Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Uint64Slice) Sort()              { sort.Sort(s) }

type chan_ts chan model.Timestamp

func (tm *TimestampManager) emitBlock(c chan_ts, block ManageableSet) {
	// Sort keys within block
	keys := make([]string, len(block))
	i := 0
	for str, _ := range block {
		keys[i] = str
		i++
	}
	sort.Strings(keys)

	// Output to chan
	for _, key := range keys {
		ts := block[key]
		c <- ts.(model.Timestamp)
	}
}

func (tm *TimestampManager) emitTimestamps(c chan_ts, bh []uint64) {
	// Iterate through blocks
	for _, h := range bh {
		block := tm.GetGroup(strconv.FormatUint(uint64(h), 10))
		tm.emitBlock(c, block)
	}
}

func (tm *TimestampManager) sortedBlocks() (Uint64Slice, error) {
	blocks := tm.ObjectManager.by_group

	// Get list of block heights
	block_heights := make(Uint64Slice, len(blocks))
	i := 0
	for h := range blocks {
		int_height, err := strconv.ParseUint(h, 10, 64)
		if err != nil {
			return nil, err
		}
		block_heights[i] = int_height
		i++
	}

	// Sort and return
	block_heights.Sort()
	return block_heights, nil
}

func (tm *TimestampManager) Iter() (<-chan model.Timestamp, error) {
	c := make(chan model.Timestamp)
	sorted_blocks, err := tm.sortedBlocks()
	if err != nil {
		return c, err
	}

	go func() {
		defer close(c)
		tm.emitTimestamps(c, sorted_blocks)
	}()
	return c, nil
}
