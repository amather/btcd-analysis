package utils

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/btcsuite/btcd/wire"
)

func SwapEndianess(value uint32) uint32 {
	return ((value & 0xff) << 24) | ((value & 0xff00) << 8) | ((value & 0xff0000) >> 8) | ((value & 0xff000000) >> 24)
}

func ReverseArray(arr []byte) []byte {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func GetHeaderSlice(header *wire.BlockHeader) []byte {
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

func CalcHash(header *wire.BlockHeader) []byte {

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
	hdr := GetHeaderSlice(header)

	first := sha256.New()
	first.Write(hdr)
	return ReverseArray(first.Sum(nil))

	/*
		second := sha256.New()
		second.Write(first.Sum(nil))

		log.Printf("second hash: %x", reverseArray(second.Sum(nil)))
		return nil
	*/
}

func Bitcount(input []byte) uint32 {

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
