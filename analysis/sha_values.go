package analysis

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func AnalyzeSHAValues(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[shavalues] starting.")

	vars := []string{"a", "b", "c", "d", "e", "f", "g", "h", "t1", "t2", "kw"}
	varsCount := len(vars)

	last_rounds := 4
	var files_be []*os.File = make([]*os.File, varsCount*last_rounds)

	for i := 0; i < last_rounds; i++ {
		for j, v := range vars {
			fname := fmt.Sprintf("shavalues_be_round%d_%s.dat", (64 - i), v)
			f, err := os.Create(path.Join(config.Path, fname))
			if err != nil {
				panic(err)
			}
			//log.Printf("files_be[%d] = %s", (i*varsCount)+j, fmt.Sprintf("shavalues_be_round%d_%s.dat", (64-i), v))
			files_be[(i*varsCount)+j] = f
			defer f.Close()
		}
	}

	//var startBlock int64 = 600_000
	var startBlock int64 = 715_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	for i, b := range blks {

		bHash, res := utils.CalcHashCustom(b)
		for j := 63; j >= 0; j-- {

			if j < 64-last_rounds {
				break
			}

			fidx := (63 - j) * varsCount

			//log.Printf("writing to files_be[%d] = %s", fidx+0, fmt.Sprintf("%d %d\n", int(startBlock)+i, res.AValues[j]))
			files_be[fidx+0].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.AValues[j]))
			files_be[fidx+1].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.BValues[j]))
			files_be[fidx+2].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.CValues[j]))
			files_be[fidx+3].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.DValues[j]))
			files_be[fidx+4].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.EValues[j]))
			files_be[fidx+5].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.FValues[j]))
			files_be[fidx+6].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.GValues[j]))
			files_be[fidx+7].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.HValues[j]))
			files_be[fidx+8].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.T1Values[j]))
			files_be[fidx+9].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.T2Values[j]))
			files_be[fidx+10].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.KWValues[j]))
		}

		if i%10_000 == 0 {
			/*
				second := &utils.Asha256{}
				second.Init()
				second.Update(bHash, 32)
				secondHash := make([]byte, 32)
				second.Final(&secondHash)
			*/
			log.Printf("[shavalues] current block: %d (%v), hash: %s", int(startBlock)+i, b.Timestamp, hex.EncodeToString(utils.ReverseArray(bHash)))
		}

	}
}
