package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type header struct {
	SaveType    [1]byte
	Null0       [3]byte
	SaveSize    [3]byte
	Null1       [1]byte
	Pointer1    [1]byte
	Pointer2    [1]byte
	Region      [2]byte
	ProductCode [10]byte
	Identifier  [8]byte
	Null2       [97]byte
	XOR         [1]byte
}

type team struct {
	Null0     [4]byte
	TeamName  [17]byte
	WormNames [4][17]byte

	Null1                   [38]byte // Team weapon, CPU etc
	TeamLossCount           uint32
	TeamDeathmatchLossCount uint32
	Null2                   [4]byte // No idea
	TeamWinCount            uint32
	TeamDeathmatchWinCount  uint32
	Null3                   [16]byte // Draw count and something else ?
	TeamKills               uint32
	TeamDeathmatchKills     uint32
	Null4                   [4]byte // No idea
	TeamDeaths              uint32
	TeamDeathmatchDeaths    uint32
	Null5                   [649]byte // Zeros
}

type save struct {
	Null0 [512]byte
	Teams [9]team
	Null1 [192]byte
}

type card struct {
	Null0   [128]byte
	Headers [15]header
	Null1   [6144]byte
	Saves   [15]save
}

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("usage: %s [MCR_FILE]", os.Args[0]))
	}

	mcr := os.Args[1]

	fmt.Printf("Parsing: \"%s\"\n", mcr)

	file, err := os.Open(mcr)
	if err != nil {
		panic(err)
	}

	c := card{}

	rawMemoryCard := make([]byte, binary.Size(c))
	_, err = file.Read(rawMemoryCard)
	if err != nil {
		panic(err)
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, rawMemoryCard)
	if err != nil {
		panic(err)
	}

	err = binary.Read(buf, binary.BigEndian, &c)
	if err != nil {
		panic(err)
	}

	wormsIdentifier := "waoption"
	var wormsIdentifierBytes [8]byte
	copy(wormsIdentifierBytes[:], wormsIdentifier)

	wormsSlot := -1

	for i := 0; i < 15; i++ {
		if c.Headers[i].Identifier == wormsIdentifierBytes {
			wormsSlot = i
			fmt.Printf("Found Worms save at slot %d\n", i)
		}
	}

	if wormsSlot == -1 {
		panic("No worms save found")
	}

	s := c.Saves[wormsSlot]

	if binary.Size(s) != 8192 {
		panic(fmt.Sprintf("Expected save to be 8192 bytes; got %d", binary.Size(s)))
	}

	for t := 0; t < 9; t++ {
		fmt.Printf("TeamName:                %s\n", s.Teams[t].TeamName)
		for i := 0; i < 4; i++ {
			fmt.Printf("Worm %d:                  %s\n", i, s.Teams[t].WormNames[i])
		}

		fmt.Printf("TeamLossCount:           %d\n", s.Teams[t].TeamLossCount)
		fmt.Printf("TeamDeathmatchLossCount: %d\n", s.Teams[t].TeamDeathmatchLossCount)
		fmt.Printf("TeamWinCount:            %d\n", s.Teams[t].TeamWinCount)
		fmt.Printf("TeamDeathmatchWinCount:  %d\n", s.Teams[t].TeamDeathmatchWinCount)
		fmt.Printf("TeamKills:               %d\n", s.Teams[t].TeamKills)
		fmt.Printf("TeamDeathmatchKills:     %d\n", s.Teams[t].TeamDeathmatchKills)
		fmt.Printf("TeamDeaths:              %d\n", s.Teams[t].TeamDeaths)
		fmt.Printf("TeamDeathmatchDeaths:    %d\n", s.Teams[t].TeamDeathmatchDeaths)

		fmt.Printf("\n")
	}

}
