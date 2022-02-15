// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/amather/btcd-analysis/analysis"
	"github.com/amather/btcd-analysis/utils"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

func main() {

	/*
		var h = &utils.Asha256{}
		h.Init()
		h.Update(make([]byte, 0), 0)
		res := make([]byte, 32)
		h.Final(&res)
	*/

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

	//analysis.AnalyzeNonce(cfg, client, blockCount)
	//analysis.AnalyzeHashDist(cfg, client, blockCount)
	//analysis.AnalyzeBitsumHash(cfg, client, blockCount)
	//analysis.AnalyzeBitcountHash(cfg, client, blockCount)
	//analysis.AnalyzeBitcountHeader(cfg, client, blockCount)
	//analysis.AnalyzeSHAValues(cfg, client, blockCount)
	analysis.AnalyzeHeaderBitHits(cfg, client, blockCount)

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
	time.AfterFunc(time.Second*20, func() {
		log.Println("Client shutting down...")
		client.Shutdown()
		log.Println("Client shutdown complete.")
	})

	// Wait until the client either shuts down gracefully (or the user
	// terminates the process with Ctrl+C).
	client.WaitForShutdown()
}
