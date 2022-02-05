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

	var startBlock int64 = 600_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	f_bitcount_header, err := os.Create(path.Join(config.Path, "bitcount_header.dat"))
	if err != nil {
		panic(err)
	}
	defer f_bitcount_header.Close()

	for i, b := range blks {

		if i%10_000 == 0 {
			log.Printf("[bitcount_header] current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

		bitCount := utils.Bitcount(utils.GetHeaderSlice(b))
		f_bitcount_header.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, bitCount))
	}
}
