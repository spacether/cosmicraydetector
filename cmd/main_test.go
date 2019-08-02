package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryBlock(t *testing.T) {
	reqMibSizes := []int{10, 100, 1000}
	for _, reqMibSize := range reqMibSizes {
		memBlock := getMemoryBlock(reqMibSize)
		byteSize, mibSize := byteMibSize(memBlock)
		if mibSize != reqMibSize {
			t.Errorf("mibSize was incorrect, got: %d, want: %d.", mibSize, reqMibSize)
		}
		reqByteSize := reqMibSize * bytesPerMib
		if byteSize != reqByteSize {
			t.Errorf("byteSize was incorrect, got: %d, want: %d.", byteSize, reqByteSize)
		}
	}
}

func TestCheckBitFlip(t *testing.T) {
	memBlock := getMemoryBlock(100)
	flips := make([]flip, 0)
	indexedFlips := make(map[int][]*flip)
	startTime := time.Now()
	checkBitFlip(memBlock, &flips, indexedFlips, startTime)
	reqLen := 0
	if len(flips) != reqLen {
		t.Errorf("flips length is incorrect, got: %d, want: %d.", len(flips), reqLen)
	}

	rows := []struct {
		flipIndex     int
		value         uint64
		expectedFlip  flip
		flipsLen      int
		flipsLenAtInd int
	}{
		{0, 1, flip{Value: 1, Binary: "00000001", ChangedBits: "_______X", NumChangedBits: 1}, 1, 1},
		{0, 3, flip{Value: 3, Binary: "00000011", ChangedBits: "______X_", NumChangedBits: 1}, 2, 2},
		{1, 7, flip{Value: 7, Binary: "00000111", ChangedBits: "_____XXX", NumChangedBits: 3}, 3, 1},
	}

	for _, row := range rows {
		// detects a single bit flip
		i := row.flipIndex
		memBlock[row.flipIndex] = row.value
		checkBitFlip(memBlock, &flips, indexedFlips, startTime)
		if len(indexedFlips[i]) != row.flipsLenAtInd {
			t.Errorf("indexedFlips[%d] length is incorrect, got: %d, want: %d.", i, len(indexedFlips[i]), row.flipsLenAtInd)
		}
		if len(flips) != row.flipsLen {
			t.Errorf("flips length is incorrect, got: %d, want: %d.", len(flips), row.flipsLen)
		}

		f := *indexedFlips[i][len(indexedFlips[i])-1]
		expected := row.expectedFlip
		expected.Duration = f.Duration
		expected.Time = f.Time
		assert.Equal(t, f, expected, "flip stored at index %d is incorrect.", i)
	}
}
