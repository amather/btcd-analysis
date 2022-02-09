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

	last_rounds := 4
	var files_be []*os.File = make([]*os.File, 8*last_rounds)

	for i := 0; i < last_rounds; i++ {
		for j, v := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
			fname := fmt.Sprintf("shavalues_be_round%d_%s.dat", (64 - i), v)
			f, err := os.Create(path.Join(config.Path, fname))
			if err != nil {
				panic(err)
			}
			//log.Printf("files_be[%d] = %s", (i*8)+j, fmt.Sprintf("shavalues_be_round%d_%s.dat", (64-i), v))
			files_be[(i*8)+j] = f
			defer f.Close()
		}
	}

	/*
		var files_le []*os.File = make([]*os.File, 0, 8)
		var files_be []*os.File = make([]*os.File, 0, 8)
		for i := 0; i < 8; i++ {
			fname := fmt.Sprintf("shavalues_le_%d.dat", i+1)
			f, err := os.Create(path.Join(config.Path, fname))
			if err != nil {
				panic(err)
			}
			files_le = append(files_le, f)
			defer f.Close()

			fname = fmt.Sprintf("shavalues_be_%d.dat", i+1)
			f, err = os.Create(path.Join(config.Path, fname))
			if err != nil {
				panic(err)
			}
			files_be = append(files_be, f)
			defer f.Close()
		}
	*/

	var startBlock int64 = 600_000
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

			fidx := (63 - j) * 8

			//log.Printf("writing to files_be[%d] = %s", fidx+0, fmt.Sprintf("%d %d\n", int(startBlock)+i, res.AValues[j]))
			files_be[fidx+0].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.AValues[j]))
			files_be[fidx+1].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.BValues[j]))
			files_be[fidx+2].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.CValues[j]))
			files_be[fidx+3].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.DValues[j]))
			files_be[fidx+4].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.EValues[j]))
			files_be[fidx+5].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.FValues[j]))
			files_be[fidx+6].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.GValues[j]))
			files_be[fidx+7].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, res.HValues[j]))
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

		/*
			bHashLe := utils.CalcHash(b)
				files_le[0].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[0:4])))
				files_le[1].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[4:8])))
				files_le[2].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[8:12])))
				files_le[3].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[12:16])))
				files_le[4].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[16:20])))
				files_le[5].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[20:24])))
				files_le[6].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[24:28])))
				files_le[7].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashLe[28:32])))

				bHashBe := utils.ReverseArray(bHashLe)
				files_be[0].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[0:4])))
				files_be[1].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[4:8])))
				files_be[2].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[8:12])))
				files_be[3].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[12:16])))
				files_be[4].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[16:20])))
				files_be[5].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[20:24])))
				files_be[6].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[24:28])))
				files_be[7].WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(bHashBe[28:32])))
		*/
	}
}
