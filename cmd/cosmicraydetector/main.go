package main

import (
  "fmt"
  "strconv"
  "time"
)

const bytesPerMib = 1048576
const bytesPerUint64 = 8


type flip struct {
    Value uint64
    Binary  string
    DeltaBits  int
    Duration time.Duration // how long the bit was stored before it was flipped
    Time time.Time // when the bit was flipped
}

// TODO:
// better structure is map of index to array of flip
// then when it happens compare it to the previous one

func main() {
  fmt.Println("Starting cosmic ray detector")

  reqMibSize := 1000
  startTime := time.Now()
  memBlock := getMemoryBlock(reqMibSize)
  flips := make(map[int]flip)
  delaySecs := 60
  infiniteLoop(delaySecs, memBlock, flips, startTime)
}

func getMemoryBlock(reqMibSize int) []uint64 {
  blockLength := reqMibSize * bytesPerMib / bytesPerUint64
  memBlock := make([]uint64, blockLength)
  byteSize := blockLength * bytesPerUint64
  mibSize := byteSize / bytesPerMib
  fmt.Printf("memBlock: %T, %d bytes, %d MiB\n", memBlock, byteSize, mibSize)
  return memBlock
}

func infiniteLoop(delaySecs int, memBlock []uint64, flips map[int]flip, startTime time.Time) {
  for {
    time.Sleep(time.Duration(delaySecs) * time.Second)
    checkBitFlip(memBlock, flips, startTime)
    fmt.Printf("Slept %d seconds: len(flips)=%d\n", delaySecs, len(flips))
  }
}

func checkBitFlip(memBlock []uint64, flips map[int]flip, startTime time.Time) {
  for i, val := range memBlock {
    if val != 0 {
      binary := strconv.FormatInt(int64(val), 2)
      bits := 0
      for _, char := range binary {
        if char == '1' {
          bits += 1
        }
      }
      t := time.Now()
      duration := t.Sub(startTime)
      f := flip{
        Value: val,
        Binary: binary,
        DeltaBits: bits,
        Duration: duration,
        Time: t,
      }
      fmt.Printf("%v", f)
      flips[i] = f
    }
  }
}