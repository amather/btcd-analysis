package nonce

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func Analyze(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[nonce] starting.")

	var startBlock int64 = 0
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	f_nonce_le, err := os.Create(path.Join(config.Path, "nonce_le.dat"))
	if err != nil {
		panic(err)
	}
	defer f_nonce_le.Close()

	f_nonce_be, err := os.Create(path.Join(config.Path, "nonce_be.dat"))
	if err != nil {
		panic(err)
	}
	defer f_nonce_be.Close()

	for i, b := range blks {

		if i%10_000 == 0 {
			log.Printf("[nonce] current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

		var nonce_le uint32 = b.Nonce
		var nonce_be uint32 = utils.SwapEndianess(b.Nonce)

		f_nonce_le.WriteString(fmt.Sprintf("%d %d\r\n", int(i)+1, nonce_le))
		f_nonce_be.WriteString(fmt.Sprintf("%d %d\r\n", int(i)+1, nonce_be))
	}
}
