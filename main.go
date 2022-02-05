// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/amather/btcd-analysis/nonce"
	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func swapEndianess(value uint32) uint32 {
	return ((value & 0xff) << 24) | ((value & 0xff00) << 8) | ((value & 0xff0000) >> 8) | ((value & 0xff000000) >> 24)
}

func reverseArray(arr []byte) []byte {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func getHeaderSlice(header *wire.BlockHeader) []byte {
	var word []byte = make([]byte, 4)
	var hdr []byte = make([]byte, 0, 80)
	binary.LittleEndian.PutUint32(word, uint32(header.Version))
	hdr = append(hdr, word...)
	hdr = append(hdr, header.PrevBlock.CloneBytes()...)
	//mrkl := reverseArray(header.MerkleRoot.CloneBytes())
	mrkl := header.MerkleRoot.CloneBytes()
	hdr = append(hdr, mrkl...)
	ts := header.Timestamp.Unix()
	binary.LittleEndian.PutUint32(word, uint32(ts))
	hdr = append(hdr, word...)
	binary.LittleEndian.PutUint32(word, uint32(header.Bits))
	hdr = append(hdr, word...)
	binary.LittleEndian.PutUint32(word, uint32(header.Nonce))
	hdr = append(hdr, word...)
	//hdr = append(hdr, byte(swapEndianess(([]byte)header.Version)))

	return hdr
}

func calcHash(header *wire.BlockHeader) []byte {

	/*
		var word []byte = make([]byte, 4)
		var hdr []byte = make([]byte, 0, 80)
		binary.LittleEndian.PutUint32(word, uint32(header.Version))
		hdr = append(hdr, word...)
		hdr = append(hdr, header.PrevBlock.CloneBytes()...)
		//mrkl := reverseArray(header.MerkleRoot.CloneBytes())
		mrkl := header.MerkleRoot.CloneBytes()
		hdr = append(hdr, mrkl...)
		binary.LittleEndian.PutUint32(word, uint32(header.Timestamp.Unix()))
		hdr = append(hdr, word...)
		binary.LittleEndian.PutUint32(word, uint32(header.Bits))
		hdr = append(hdr, word...)
		binary.LittleEndian.PutUint32(word, uint32(header.Nonce))
		hdr = append(hdr, word...)
		//hdr = append(hdr, byte(swapEndianess(([]byte)header.Version)))
	*/
	hdr := getHeaderSlice(header)

	first := sha256.New()
	first.Write(hdr)
	return reverseArray(first.Sum(nil))

	/*
		second := sha256.New()
		second.Write(first.Sum(nil))

		log.Printf("second hash: %x", reverseArray(second.Sum(nil)))
		return nil
	*/
}

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

func bitcount(input []byte) uint32 {

	var res uint32 = 0

	for i := 0; i < len(input); i++ {
		b := input[i]
		for j := 0; j < 8; j++ {
			mask := uint8(1 << j)
			val := uint32(b&mask) >> j

			res += val
		}
	}

	return res
}

func writeNonce(client *rpcclient.Client) {
	f_nonce, err := os.Create("nonce.dat")
	if err != nil {
		panic(err)
	}
	f_hash_1, err := os.Create("hash_1.dat")
	if err != nil {
		panic(err)
	}
	f_hash_2, err := os.Create("hash_2.dat")
	if err != nil {
		panic(err)
	}
	f_hash_3, err := os.Create("hash_3.dat")
	if err != nil {
		panic(err)
	}
	f_hash_4, err := os.Create("hash_4.dat")
	if err != nil {
		panic(err)
	}
	f_hash_5, err := os.Create("hash_5.dat")
	if err != nil {
		panic(err)
	}
	f_hash_6, err := os.Create("hash_6.dat")
	if err != nil {
		panic(err)
	}
	f_hash_7, err := os.Create("hash_7.dat")
	if err != nil {
		panic(err)
	}
	f_hash_8, err := os.Create("hash_8.dat")
	if err != nil {
		panic(err)
	}
	f_bitcount, err := os.Create("bitcount.dat")
	if err != nil {
		panic(err)
	}
	f_bitcount_hdr, err := os.Create("bitcount_header.dat")
	if err != nil {
		panic(err)
	}
	defer f_nonce.Close()
	defer f_hash_1.Close()
	defer f_hash_2.Close()
	defer f_hash_3.Close()
	defer f_hash_4.Close()
	defer f_hash_5.Close()
	defer f_hash_6.Close()
	defer f_hash_7.Close()
	defer f_hash_8.Close()
	defer f_bitcount.Close()
	defer f_bitcount_hdr.Close()

	// TODO: 10 block moving average (10BMA) bitcount
	// TODO: nBits changes
	// TODO: bitcount - nonce (biased due to HW)

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current block height: %d", blockCount-1)

	var startBlock int64 = 620_000
	//var startBlock int64 = 200_000
	//var startBlock int64 = 0
	blockCount = 720_000
	blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	//blks, err := client.GetBlockHeaderRange(startBlock, blockCount)
	if err != nil {
		panic(err)
	}

	var bitSum = make([]uint32, 256)

	for i, b := range blks {
		//log.Printf("Block: %v, hash: %v", b.Timestamp, b.BlockHash().String())

		if i%10_000 == 0 {
			log.Printf("current block: %d (%v)", int(startBlock)+i, b.Timestamp)
		}

		var nonce_be uint32 = swapEndianess(b.Nonce)
		//nonce_be = ((b.Nonce & 0xff) << 24) | ((b.Nonce & 0xff00) << 8) | ((b.Nonce & 0xff0000) >> 8) | ((b.Nonce & 0xff000000) >> 24)
		f_nonce.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, nonce_be))
		//f.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, b.Nonce))

		bitCountHeader := bitcount(getHeaderSlice(b))
		f_bitcount_hdr.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, bitCountHeader))

		//bitCount := bitcount(calcHash(b))
		//f_bitcount.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, bitCount))
		f_hash := calcHash(b)
		bitsum(bitSum, f_hash)
		bitCount := bitcount(f_hash)
		f_bitcount.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, bitCount))
		f_hash_1.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[0:4])))
		f_hash_2.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[4:8])))
		f_hash_3.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[8:12])))
		f_hash_4.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[12:16])))
		f_hash_5.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[16:20])))
		f_hash_6.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[20:24])))
		f_hash_7.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[24:28])))
		f_hash_8.WriteString(fmt.Sprintf("%d %d\n", int(startBlock)+i, binary.LittleEndian.Uint32(f_hash[28:32])))
	}

	f_bitsum, err := os.Create("bitsum.dat")
	if err != nil {
		panic(err)
	}
	defer f_bitsum.Close()
	for i := 0; i < len(bitSum); i++ {
		f_bitsum.WriteString(fmt.Sprintf("%d %d\n", i, bitSum[i]))
	}
}

func main() {
	// Only override the handlers for notifications you care about.
	// Also note most of these handlers will only be called if you register
	// for notifications.  See the documentation of the rpcclient
	// NotificationHandlers type for more details about each handler.
	ntfnHandlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*btcutil.Tx) {
			log.Printf("Block connected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			log.Printf("Block disconnected: %v (%d) %v",
				header.BlockHash(), height, header.Timestamp)
		},
	}

	// Connect to local btcd RPC server using websockets.
	btcdHomeDir := btcutil.AppDataDir("btcd", false)
	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Fatal(err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         "localhost:8334",
		Endpoint:     "ws",
		User:         "btc",
		Pass:         "emigax07tip03top",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &ntfnHandlers)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &utils.Config{
		Path:     "out",
		FullPath: "",
	}

	// create output directory if it doesn't exist
	if _, err := os.Stat(cfg.Path); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(cfg.Path, fs.ModeDir)
			if err != nil {
				panic(err)
			}
		}
	}
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg.FullPath = path.Join(pwd, cfg.Path)

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[main] current block height: %d", blockCount-1)

	nonce.Analyze(cfg, client, blockCount)

	//writeNonce(client)

	/*
		blkh, err := client.GetBlockHash(721553)
		if err != nil {
			log.Fatal(err)
		}
		blkhdr, err := client.GetBlock(blkh)
		if err != nil {
			log.Fatal(err)
		}
		//l := fmt.Sprintf("blk 400_000: version=%v, hash=0x%v, nonce=0x%x", blkhdr.Header.Version, blkhdr.Header.BlockHash(), blkhdr.Header.Nonce)
		l := fmt.Sprintf("blk 721_553: hash=0x%v, nonce=%d", blkhdr.Header.BlockHash(), blkhdr.Header.Nonce)
		log.Println(l)
	*/

	/*
		// Register for block connect and disconnect notifications.
		if err := client.NotifyBlocks(); err != nil {
			log.Fatal(err)
		}
		log.Println("NotifyBlocks: Registration Complete")

		// Get the current block count.
		blockCount, err := client.GetBlockCount()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Block count: %d", blockCount)
	*/

	// For this example gracefully shutdown the client after 10 seconds.
	// Ordinarily when to shutdown the client is highly application
	// specific.
	log.Println("Client shutdown in 1 minute...")
	time.AfterFunc(time.Second*60, func() {
		log.Println("Client shutting down...")
		client.Shutdown()
		log.Println("Client shutdown complete.")
	})

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
