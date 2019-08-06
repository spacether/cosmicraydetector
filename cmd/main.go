package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime/debug"
	"time"
	"unicode/utf8"
)

const bytesPerMib = 1048576
const bytesPerUint64 = 8
const uint64FmtString = "%08b"

type flip struct {
	Value          uint64
	Binary         string
	NumChangedBits int
	ChangedBits    string
	Duration       time.Duration // how long the value was stored before it was changed
	Time           time.Time     // when the bit flip happened
}

var blockSize = flag.Int("blockSize", 1000, "memory block size in MiB")

func main() {
	fmt.Println("Starting cosmic ray detector")
	// the SetGCPercent input is int %, and triggers garbage collection
	// the default is 100%, so if our program uses 1 GiB, we don't want to wait
	// for 2 GiB of usage, instead wait for 1,100 of usage
	debug.SetGCPercent(10)

	flag.Parse()
	reqMibSize := *blockSize
	startTime := time.Now()
	memBlock := getMemoryBlock(reqMibSize)
	flips := make([]flip, 0)
	indexedFlips := make(map[int][]*flip)

	// handle interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanUp(&memBlock, &flips, &indexedFlips)
		os.Exit(0)
	}()

	t1 := time.Now()
	checkBitFlip(memBlock, &flips, indexedFlips, startTime)
	t2 := time.Now()
	elapsed := t2.Sub(t1)
	delaySecs := 60 * (int(math.Ceil(elapsed.Seconds()/60.0)) + 1)
	infiniteLoop(delaySecs, memBlock, &flips, indexedFlips, startTime)
}

func cleanUp(memBlock *[]uint64, flips *[]flip, indexedFlips *map[int][]*flip) {
	fmt.Printf("\nCleaning up\n")
	*memBlock = make([]uint64, 0)
	*flips = make([]flip, 0)
	*indexedFlips = make(map[int][]*flip)
	debug.FreeOSMemory()
	fmt.Printf("Cleanup Done\n")
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
		t1 := time.Now()
		checkBitFlip(memBlock, flips, indexedFlips, startTime)
		t2 := time.Now()
		elapsed := t2.Sub(t1)
		fmt.Printf("Slept %d seconds: time=%v inspection_duration=%v len(flips)=%d\n", delaySecs, t2, elapsed, len(*flips))
		remainder := time.Duration(delaySecs)*time.Second - elapsed
		time.Sleep(remainder)
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
				Value:          val,
				Binary:         binary,
				ChangedBits:    changedBits,
				NumChangedBits: bits,
				Duration:       duration,
				Time:           t,
			}
			fmt.Printf("%v\n", f)
			*flips = append(*flips, f)
			indexedFlips[i] = append(indexedFlips[i], &f)
		}
	}
}
