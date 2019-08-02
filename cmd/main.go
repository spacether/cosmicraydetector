package main

import (
	"fmt"
	"time"
	"unicode/utf8"
)

const bytesPerMib = 1048576
const bytesPerUint64 = 8
const uint64FmtString = "%08b"

type flip struct {
	Value     uint64
	Binary    string
	NumChangedBits int
	ChangedBits string
	Duration  time.Duration // how long the value was stored before it was changed
	Time      time.Time     // when the bit flip happened
}

func main() {
	fmt.Println("Starting cosmic ray detector")

	reqMibSize := 1000
	startTime := time.Now()
	memBlock := getMemoryBlock(reqMibSize)
	flips := make([]flip, 0)
	indexedFlips := make(map[int][]*flip)
	delaySecs := 60
	infiniteLoop(delaySecs, memBlock, &flips, indexedFlips, startTime)
}

func getMemoryBlock(reqMibSize int) []uint64 {
	blockLength := reqMibSize * bytesPerMib / bytesPerUint64
	memBlock := make([]uint64, blockLength)
	byteSize, mibSize := byteMibSize(memBlock)
	fmt.Printf("memBlock: %T, %d bytes, %d MiB\n", memBlock, byteSize, mibSize)
	return memBlock
}

func byteMibSize(memBlock []uint64) (byteSize int, mibSize int) {
	byteSize = len(memBlock) * bytesPerUint64
	mibSize = byteSize / bytesPerMib
	return
}

func infiniteLoop(delaySecs int, memBlock []uint64, flips *[]flip, indexedFlips map[int][]*flip, startTime time.Time) {
	for {
		time.Sleep(time.Duration(delaySecs) * time.Second)
		checkBitFlip(memBlock, flips, indexedFlips, startTime)
		fmt.Printf("Slept %d seconds: len(flips)=%d\n", delaySecs, len(*flips))
	}
}

func checkBitFlip(memBlock []uint64, flips *[]flip, indexedFlips map[int][]*flip, startTime time.Time) {
	for i, val := range memBlock {
		fs, ok := indexedFlips[i]
		var oldVal uint64
		if ok {
			f := *fs[len(fs)-1]
			oldVal = f.Value
		}
		oldBinary := fmt.Sprintf(uint64FmtString, oldVal)
		if val != oldVal {
			binary := fmt.Sprintf(uint64FmtString, val)
			bits := 0
			var changedBits string
			for i, r := range binary {
				oldR, _ := utf8.DecodeRuneInString(oldBinary[i:])
				if r != oldR {
					bits++
					changedBits += "X"
				} else {
					changedBits += "_"
				}
			}
			t := time.Now()
			duration := t.Sub(startTime)
			f := flip{
				Value:     val,
				Binary:    binary,
				ChangedBits: changedBits,
				NumChangedBits: bits,
				Duration:  duration,
				Time:      t,
			}
			fmt.Printf("%v\n", f)
			*flips = append(*flips, f)
			indexedFlips[i] = append(indexedFlips[i], &f)
		}
	}
}
