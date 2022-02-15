package analysis

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/amather/shanalyzer/arch"
	"github.com/amather/shanalyzer/bitcoin"
	"github.com/amather/shanalyzer/sha256"
	"github.com/btcsuite/btcd/rpcclient"
)

func AnalyzeHeaderBitHits(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[bithits] starting.")

	//var startBlock int64 = 600_000
	var startBlock int64 = 715_000
	var endBlock int64 = 715_001
	//blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	blks, err := client.GetBlockHeaderRange(startBlock, endBlock)
	if err != nil {
		panic(err)
	}

	totalBlocks := (endBlock - startBlock) + 1
	var files_be []*os.File = make([]*os.File, totalBlocks)

	for i, j := startBlock, 0; i <= endBlock; i++ {

		fname := fmt.Sprintf("header_bithits_%d.dat", i)
		f, err := os.Create(path.Join(config.Path, fname))
		if err != nil {
			panic(err)
		}
		//log.Printf("files_be[%d] = %s", (i*varsCount)+j, fmt.Sprintf("shavalues_be_round%d_%s.dat", (64-i), v))
		files_be[j] = f
		j++
		defer f.Close()
	}

	for i, b := range blks {
		log.Printf("[bithits] current block: %d (%v)", int(startBlock)+i, b.Timestamp)

		var buf bytes.Buffer
		b.Serialize(&buf)
		hdrwords := bitcoin.FromBlockHeaderBuffer(buf.Bytes())

		// double hash
		ctx := sha256.NewSha256Context()
		ctx.Update(hdrwords, 80)
		h1, _ := ctx.Final()
		ctx = sha256.NewSha256Context()
		ctx.Update(sha256.SwapResultEndianess(h1[:], false), 32)
		h2, _ := ctx.Final()

		for j := 0; j < 8; j++ {
			word := h2[j]
			for k := 0; k < 32; k++ {
				bit := word.Bits[k]

				// parse
				pr := arch.NewParserResult()
				arch.ParseBitAll(bit, 0, pr, false)

				for _, ib := range pr.InputBits {
					pos := ib.InputPos
					hits := ib.Hits

					f := files_be[i]
					f.WriteString(fmt.Sprintf("%d %d\n", pos, hits))
				}

				arch.ParseBit(bit, 0, pr, true)

			}
		}

		/*
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
		*/

		/*
			if i%10_000 == 0 {
				//log.Printf("[bithits] current block: %d (%v), hash: %s", int(startBlock)+i, b.Timestamp, hex.EncodeToString(utils.ReverseArray(bHash)))
			}
		*/

	}
	log.Printf("[bithits] ending.")
}
