package main

import (
	"encoding/binary"
	"io/ioutil"
	"math/rand"
	"os"
)

// 8 byte = 64 bit

type Cipher struct {
	p, q, phi, d, e, n uint32
}

func gcd(u, n uint32) uint32 {
	n1, n2 := u, n
	q := n1 / n2
	r := n1 - q*n2
	if r == 0 {
		return n2
	} else {
		return gcd(n2, r)
	}
}

// n2*b2 = 1 mod n1
// n1 = n
// initial b1 =0 b2=1
func exgcd(n1, n2, b1, b2 int32) int32 {
	q := n1 / n2
	r := n1 - q*n2
	if r != 0 {
		return exgcd(n2, r, b2, b1-q*b2)
	}
	if n2 == 1 {
		return b2
	} else {
		println("No Inverse")
		return 0
	}
}

/*
	d = x ^ r mod n
	initial d = 1
*/
func SwiftModBigNum(x, r, d, n uint32) uint32 {
	a, b, c := x, r, d
	if b == 0 {
		return c
	} else if b%uint32(2) == 0 {
		return SwiftModBigNum((a*a)%n, b/2, c, n)
	}
	return SwiftModBigNum(a, b-1, (c*a)%n, n)
}

func MillerRabin(n uint32) bool {
	if n < 3 {
		return n == 2
	}
	times := 20
	for i := 1; i <= times; i++ {
		a := uint32(rand.Int31())%(n-2) + 2
		r := SwiftModBigNum(a, n-1, uint32(1), n)
		if r != 1 {
			return false
		}
	}
	return true
}
func encrypt(cipher *Cipher, Plain []byte) []uint32 {
	cipher.phi = (cipher.p - 1) * (cipher.q - 1)
	cipher.n = cipher.p * cipher.q
	if gcd(cipher.e, cipher.phi) == 1 {
		inv := exgcd(int32(cipher.phi), int32(cipher.e), int32(0), int32(1))
		if inv < 0 {
			inv = inv + int32(cipher.n)
		}
		cipher.d = uint32(inv)
	}
	CipherCode := make([]uint32, 0)
	for i := 0; i < len(Plain); i++ {
		CipherCode = append(CipherCode, SwiftModBigNum(uint32(Plain[i]), cipher.e, 1, cipher.n))
	}
	return CipherCode
}
func decrpyt(cipher Cipher, CipherCode []uint32) []uint8 {
	decrpytcode := make([]uint8, 0)
	for i := 0; i < len(CipherCode); i++ {
		decrpytcode = append(decrpytcode, uint8(SwiftModBigNum(uint32(CipherCode[i]), cipher.d, 1, cipher.n)))
	}
	return decrpytcode
}

func main() {
	Plaintext := []byte("I LOVE NANJING UNIVERSITY OF AERONAUTICS AND ASTRONAUTICS")
	cipher := Cipher{p: 191, q: 7, e: 11}
	if MillerRabin(cipher.p) || MillerRabin(cipher.q) {
		CipherCode := encrypt(&cipher, Plaintext)
		Ciphertext := make([]byte, 0)
		for i := 0; i < len(CipherCode); i++ {
			bytes := make([]byte, 4)
			binary.BigEndian.PutUint32(bytes, CipherCode[i])
			Ciphertext = append(Ciphertext, bytes...)
		}
		e := ioutil.WriteFile("ciphertext", Ciphertext, 0644)
		decrptyText := decrpyt(cipher, CipherCode)
		e = ioutil.WriteFile("decrptytext", decrptyText, 0644)
		if e != nil {
			println("file err")
			os.Exit(-1)
		}
	}
	println("7^563mod561=")
	println(SwiftModBigNum(uint32(7), uint32(563), uint32(1), uint32(561)))
}
