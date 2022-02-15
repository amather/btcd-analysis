package utils

import (
	"github.com/btcsuite/btcd/wire"
)

const (
	h0 uint32 = 0x6a09e667
	h1 uint32 = 0xbb67ae85
	h2 uint32 = 0x3c6ef372
	h3 uint32 = 0xa54ff53a
	h4 uint32 = 0x510e527f
	h5 uint32 = 0x9b05688c
	h6 uint32 = 0x1f83d9ab
	h7 uint32 = 0x5be0cd19
)

var (
	k = [...]uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
	}
	zerodata = [64]byte{}
)

func rotl(a uint32, b uint32) uint32 {
	return (((a) << (b)) | ((a) >> (32 - (b))))
}
func rotr(a uint32, b uint32) uint32 {
	return (((a) >> (b)) | ((a) << (32 - (b))))
}

func ch(x uint32, y uint32, z uint32) uint32 {
	return (((x) & (y)) ^ ((^(x)) & (z)))
}
func maj(x uint32, y uint32, z uint32) uint32 {
	return (((x) & (y)) ^ ((x) & (z)) ^ ((y) & (z)))
}

func ep0(x uint32) uint32 {
	return (rotr(x, 2) ^ rotr(x, 13) ^ rotr(x, 22))
}
func ep1(x uint32) uint32 {
	return (rotr(x, 6) ^ rotr(x, 11) ^ rotr(x, 25))
}
func sig0(x uint32) uint32 {
	return (rotr(x, 7) ^ rotr(x, 18) ^ ((x) >> 3))
}
func sig1(x uint32) uint32 {
	/*
		r1 := rotr(x, 17)
		r2 := rotr(x, 19)
		r3 := x >> 10
		x1 := r1 ^ r2
		x2 := x1 ^ r3
		return x2
	*/
	return (rotr(x, 17) ^ rotr(x, 19) ^ ((x) >> 10))
}

type Asha256Result struct {
	AValues  [64]uint32
	BValues  [64]uint32
	CValues  [64]uint32
	DValues  [64]uint32
	EValues  [64]uint32
	FValues  [64]uint32
	GValues  [64]uint32
	HValues  [64]uint32
	T1Values [64]uint32
	T2Values [64]uint32
	KWValues [64]uint32
}

func (res *Asha256Result) update(i int, a uint32, b uint32, c uint32, d uint32, e uint32, f uint32, g uint32, h uint32) {
	res.AValues[i] = a
	res.BValues[i] = b
	res.CValues[i] = c
	res.DValues[i] = d
	res.EValues[i] = e
	res.FValues[i] = f
	res.GValues[i] = g
	res.HValues[i] = h

	/*
		res.AValues[i] = SwapEndianess(res.AValues[i])
		res.BValues[i] = SwapEndianess(res.BValues[i])
		res.CValues[i] = SwapEndianess(res.CValues[i])
		res.DValues[i] = SwapEndianess(res.DValues[i])
		res.EValues[i] = SwapEndianess(res.EValues[i])
		res.FValues[i] = SwapEndianess(res.FValues[i])
		res.GValues[i] = SwapEndianess(res.GValues[i])
		res.HValues[i] = SwapEndianess(res.HValues[i])
	*/
	res.AValues[i] = SwapEndianess(a)
	res.BValues[i] = SwapEndianess(b)
	res.CValues[i] = SwapEndianess(c)
	res.DValues[i] = SwapEndianess(d)
	res.EValues[i] = SwapEndianess(e)
	res.FValues[i] = SwapEndianess(f)
	res.GValues[i] = SwapEndianess(g)
	res.HValues[i] = SwapEndianess(h)
}

type Asha256 struct {
	data    [64]byte
	datalen uint32
	bitlen  uint64
	state   [8]uint32
	result  *Asha256Result
}

func (ctx *Asha256) transform() {
	var a, b, c, d, e, f, g, h uint32
	var t1, t2 uint32
	var m [64]uint32

	for i, j := 0, 0; i < 16; {
		m[i] = (uint32(ctx.data[j]) << 24) | (uint32(ctx.data[j+1]) << 16) | (uint32(ctx.data[j+2]) << 8) | (uint32(ctx.data[j+3]))
		i++
		j += 4
	}
	for i := 16; i < 64; i++ {
		m[i] = sig1(m[i-2]) + m[i-7] + sig0(m[i-15]) + m[i-16]
	}

	a = ctx.state[0]
	b = ctx.state[1]
	c = ctx.state[2]
	d = ctx.state[3]
	e = ctx.state[4]
	f = ctx.state[5]
	g = ctx.state[6]
	h = ctx.state[7]

	for i := 0; i < 64; i++ {

		t1 = h + ep1(e) + ch(e, f, g) + k[i] + m[i]
		t2 = ep0(a) + maj(a, b, c)
		h = g
		g = f
		f = e
		e = d + t1
		d = c
		c = b
		b = a
		a = t1 + t2

		ctx.result.T1Values[i] = t1
		ctx.result.T2Values[i] = t2
		ctx.result.KWValues[i] = SwapEndianess(k[i] + m[i])
		ctx.result.update(i, a, b, c, d, e, f, g, h)

	}

	ctx.state[0] += a
	ctx.state[1] += b
	ctx.state[2] += c
	ctx.state[3] += d
	ctx.state[4] += e
	ctx.state[5] += f
	ctx.state[6] += g
	ctx.state[7] += h
}

func (h *Asha256) Init() {
	h.state[0] = h0
	h.state[1] = h1
	h.state[2] = h2
	h.state[3] = h3
	h.state[4] = h4
	h.state[5] = h5
	h.state[6] = h6
	h.state[7] = h7
	h.result = &Asha256Result{}

}
func (h *Asha256) Update(data []byte, len int) {

	for i := 0; i < len; i++ {

		h.data[h.datalen] = data[i]
		h.datalen++
		if h.datalen == 64 {
			h.transform()
			h.bitlen += 512
			h.datalen = 0
		}
	}
}
func (h *Asha256) Final(hash *[]byte) *Asha256Result {

	var i uint32 = h.datalen

	if h.datalen < 56 {
		h.data[i] = 0x80
		i++
		for i < 56 {
			h.data[i] = 0x00
			i++
		}
	} else {
		h.data[i] = 0x80
		i++
		for i < 64 {
			h.data[i] = 0x0
			i++
		}
		h.transform()
		copy(h.data[:], zerodata[:])
	}

	h.bitlen += uint64(h.datalen) * 8
	h.data[63] = byte(h.bitlen)
	h.data[62] = byte(h.bitlen >> 8)
	h.data[61] = byte(h.bitlen >> 16)
	h.data[60] = byte(h.bitlen >> 24)
	h.data[59] = byte(h.bitlen >> 32)
	h.data[58] = byte(h.bitlen >> 40)
	h.data[57] = byte(h.bitlen >> 48)
	h.data[56] = byte(h.bitlen >> 56)
	h.transform()

	for i := 0; i < 4; i++ {
		(*hash)[i] = byte((h.state[0] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+4] = byte((h.state[1] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+8] = byte((h.state[2] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+12] = byte((h.state[3] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+16] = byte((h.state[4] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+20] = byte((h.state[5] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+24] = byte((h.state[6] >> (24 - i*8)) & 0x000000ff)
		(*hash)[i+28] = byte((h.state[7] >> (24 - i*8)) & 0x000000ff)
	}

	return h.result
}

func CalcHashCustom(header *wire.BlockHeader) ([]byte, *Asha256Result) {

	hdr := GetHeaderSlice(header)
	final := make([]byte, 32)

	testHash := &Asha256{}
	testHash.Init()
	testBytes := make([]byte, 32)
	testHash.Final(&testBytes)

	firstHash := &Asha256{}
	firstHash.Init()
	firstHash.Update(hdr, len(hdr))
	//res := firstHash.Final(&final)
	firstHash.Final(&final)

	secondHash := &Asha256{}
	secondHash.Init()
	secondHash.Update(final, 32)
	res := secondHash.Final(&final)

	return final, res

	/*
		second := sha256.New()
		second.Write(first.Sum(nil))

		log.Printf("second hash: %x", reverseArray(second.Sum(nil)))
		return nil
	*/
}
