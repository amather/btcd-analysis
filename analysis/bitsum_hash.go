package analysis

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func bitsum(arr []uint32, hash []byte) {
	if len(hash) != 32 {
		return
	}

	for i := 0; i < 32; i++ {
		b := hash[i]
		for j := 0; j < 8; j++ {
			mask := uint8(1 << j)
			val := uint32(b&mask) >> j

			idx := (i * 8) + j
			arr[idx] += val
		}
	}
}

func AnalyzeBitsumHash(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[bitsum_hash] starting.")

	f_bitsum, err := os.Create(path.Join(config.Path, "bitsum_hash.dat"))
	if err != nil {
		panic(err)
	}
	defer f_bitsum.Close()

	var startBlock int64 = blockCount - 100_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	var bitSum = make([]uint32, 256)

	for i, b := range blks {

		if i%10_000 == 0 {
			log.Printf("[bitsum_hash] current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

		bHashLe := utils.CalcHash(b)
		bitsum(bitSum, bHashLe)
	}

	for i := 0; i < len(bitSum); i++ {
		f_bitsum.WriteString(fmt.Sprintf("%d %d\n", i, bitSum[i]))
	}
}
