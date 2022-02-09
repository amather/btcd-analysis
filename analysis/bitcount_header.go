package analysis

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func AnalyzeBitcountHeader(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[bitcount_header] starting.")

	var startBlock int64 = 200_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	f_bitcount_header, err := os.Create(path.Join(config.Path, "bitcount_header.dat"))
	if err != nil {
		panic(err)
	}
	defer f_bitcount_header.Close()

	maSize := 100
	ma := utils.MovingAverage{}
	ma.Init(maSize)

	for i, b := range blks {

		if i%10_000 == 0 {
			log.Printf("[bitcount_header] current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

		bitCount := utils.Bitcount(utils.GetHeaderSlice(b))

		ma.Add(int(bitCount))
		ma.Next()

		maMean := 0
		maLower := 0
		maUpper := 0
		if i >= maSize {
			state := ma.State()
			maMean = int(state.Mean)
			maLower = int(state.Mean - state.StdDev)
			maUpper = int(state.Mean + state.StdDev)
		}
		//version := (b.Version & 0x1fffe000)
		version := (uint32(b.Version) & uint32(0x1ffe000)) >> 13
		//version := (uint32(b.Version) & uint32(0xe0001fff))

		f_bitcount_header.WriteString(fmt.Sprintf("%d %d %d %d %d %x\n", int(startBlock)+i, bitCount, maMean, maLower, maUpper, version))
	}
}
