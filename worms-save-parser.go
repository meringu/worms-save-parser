package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	Null3                   [4]byte // No idea
	TeamDrawCount           uint32
	TeamDeathmatchDrawCount uint32
	Null4                   [4]byte // No idea
	TeamKills               uint32
	TeamDeathmatchKills     uint32
	Null5                   [4]byte // No idea
	TeamDeaths              uint32
	TeamDeathmatchDeaths    uint32
	Null6                   [649]byte // Zeros
}

type team_stripped struct {
	TeamName  string
	WormNames [4]string
	TeamLossCount           uint32
	TeamDeathmatchLossCount uint32
	TeamWinCount            uint32
	TeamDeathmatchWinCount  uint32
	TeamDrawCount           uint32
	TeamDeathmatchDrawCount uint32
	TeamKills               uint32
	TeamDeathmatchKills     uint32
	TeamDeaths              uint32
	TeamDeathmatchDeaths    uint32
}

type save struct {
	Null0 [512]byte
	Teams [9]team
	Null1 [192]byte
}

type save_stripped struct {
	Teams [9]team_stripped
}

type card struct {
	Null0   [128]byte
	Headers [15]header
	Null1   [6144]byte
	Saves   [15]save
}

func byte_arr_to_str(bytes [17]byte) string {
	var new_str [17]byte
	for i := 0; i <= 17; i++ {
		if bytes[i] == 0 {
			break
		}
		new_str[i] = bytes[i]
  }
  return strings.TrimRight(fmt.Sprintf("%s", new_str), "\x00")
}

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Sprintf("usage: %s [MCR_FILE]", os.Args[0]))
	}

	mcr := os.Args[1]

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
		}
	}

	if wormsSlot == -1 {
		panic("No worms save found")
	}

	s := c.Saves[wormsSlot]

	if binary.Size(s) != 8192 {
		panic(fmt.Sprintf("Expected save to be 8192 bytes; got %d", binary.Size(s)))
	}

	ss := save_stripped{}

	for t := 0; t < 9; t++ {
		ss.Teams[t].TeamName = byte_arr_to_str(s.Teams[t].TeamName)
		for i := 0; i < 4; i++ {
			ss.Teams[t].WormNames[i] = byte_arr_to_str(s.Teams[t].WormNames[i])
		}
		ss.Teams[t].TeamLossCount = s.Teams[t].TeamLossCount
		ss.Teams[t].TeamDeathmatchLossCount = s.Teams[t].TeamDeathmatchLossCount
		ss.Teams[t].TeamWinCount = s.Teams[t].TeamWinCount
		ss.Teams[t].TeamDeathmatchWinCount = s.Teams[t].TeamDeathmatchWinCount
		ss.Teams[t].TeamDrawCount = s.Teams[t].TeamDrawCount
		ss.Teams[t].TeamDeathmatchDrawCount = s.Teams[t].TeamDeathmatchDrawCount
		ss.Teams[t].TeamKills = s.Teams[t].TeamKills
		ss.Teams[t].TeamDeathmatchKills = s.Teams[t].TeamDeathmatchKills
		ss.Teams[t].TeamDeaths = s.Teams[t].TeamDeaths
		ss.Teams[t].TeamDeathmatchDeaths = s.Teams[t].TeamDeathmatchDeaths
	}
	json, _ := json.MarshalIndent(ss, "", "  ")
  fmt.Println(string(json))

}
