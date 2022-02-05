package analysis

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
)

func AnalyzeHashDist(config *utils.Config, client *rpcclient.Client, blockCount int64) {

	log.Printf("[hashdist] starting.")

	var files_le []*os.File = make([]*os.File, 0, 8)
	var files_be []*os.File = make([]*os.File, 0, 8)
	for i := 0; i < 8; i++ {
		fname := fmt.Sprintf("hashdist_le_%d.dat", i+1)
		f, err := os.Create(path.Join(config.Path, fname))
		if err != nil {
			panic(err)
		}
		files_le = append(files_le, f)
		defer f.Close()

		fname = fmt.Sprintf("hashdist_be_%d.dat", i+1)
		f, err = os.Create(path.Join(config.Path, fname))
		if err != nil {
			panic(err)
		}
		files_be = append(files_be, f)
		defer f.Close()
	}

	var startBlock int64 = 600_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	for i, b := range blks {

		if i%10_000 == 0 {
			log.Printf("[hashdist] current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

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
	}
}
