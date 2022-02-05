package analysis

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func AnalyzeBitcountHash(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[bitcount_hash] starting.")

	var startBlock int64 = 600_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	f_bitcount_hash, err := os.Create(path.Join(config.Path, "bitcount_hash.dat"))
	if err != nil {
		panic(err)
	}
	defer f_bitcount_hash.Close()

	for i, b := range blks {

		if i%10_000 == 0 {
			log.Printf("[bitcount_hash] current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

		bitCount := utils.Bitcount(utils.CalcHash(b))
		f_bitcount_hash.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, bitCount))
	}
}
