package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"math"

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

type MovingAverage struct {
	data []int
	size uint
	idx  uint
}

type MovingAverageState struct {
	Mean   float64
	StdDev float64
}

func (ma *MovingAverage) Init(size int) {
	ma.data = make([]int, size)
	ma.size = uint(size)
	ma.idx = 0
}

func (ma *MovingAverage) Add(value int) {
	ma.data[ma.idx] = value
}

func (ma *MovingAverage) Next() {
	ma.idx = (ma.idx + 1) % ma.size
}

func (ma *MovingAverage) State() *MovingAverageState {

	sum := float64(0)
	for _, v := range ma.data {
		sum += float64(v)
	}
	mean := (sum / float64(ma.size))

	//variances := make([]float64, ma.size)
	variance := float64(0)
	for _, v := range ma.data {
		//variances[i] = math.Pow(float64(v)-mean, 2)
		variance += math.Pow(float64(v)-mean, 2)
	}
	variance /= float64(ma.size)
	std_dev := math.Pow(variance, 0.5)

	return &MovingAverageState{
		Mean:   mean,
		StdDev: std_dev,
	}
}

func GetHeaderSliceReduced(header *wire.BlockHeader) []byte {
	var word []byte = make([]byte, 4)
	var hdr []byte = make([]byte, 0, 80)

	/*
		// Version
		binary.LittleEndian.PutUint32(word, uint32(header.Version))
		hdr = append(hdr, word...)
	*/

	// Prev Block
	hdr = append(hdr, header.PrevBlock.CloneBytes()...)

	// Merkle Root
	mrkl := header.MerkleRoot.CloneBytes()
	hdr = append(hdr, mrkl...)

	// Timestamp
	ts := header.Timestamp.Unix()
	binary.LittleEndian.PutUint32(word, uint32(ts))
	hdr = append(hdr, word...)

	// nBits
	binary.LittleEndian.PutUint32(word, uint32(header.Bits))
	hdr = append(hdr, word...)

	// Nonce
	binary.LittleEndian.PutUint32(word, uint32(header.Nonce))
	hdr = append(hdr, word...)

	return hdr
}

func ExpectedBitcountMean(header *wire.BlockHeader) {

	/*
		// header.Version
		slc := GetHeaderSliceReduced(header)

		totalFairBits := 0

		// PrevBlock
		totalFairBits += 256

		// MerkleRoot
		totalFairBits += 256

		// Timestamp
		/// bits  seconds
		///   1        2
		///   2        4
		///   3        8
		///   ...
		///   8		 256 ( 4.26min)
		///   9		 512 ( 8.53min)
		///  10  	1024 (17.06min)
		//header.Timestamp.Unix();
		totalFairBits += (32 - 10)

		// nBits
		nBitsSlice := make([]byte, 4)
		binary.LittleEndian.PutUint32(nBitsSlice, header.Bits)
		nBitsCount := Bitcount(nBitsSlice)
		beBits := SwapEndianess(header.Bits)


		//totalFairBits +=

		//header.Nonce
	*/
}
