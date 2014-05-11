// based on https://github.com/crowsonkb/base58/blob/master/base58.go
// implements modified base54 used in ICE

package ice

import (
    "fmt"
    "math"
    "math/big"
    )

// The 58-character encoding alphabet.
const Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/"

// The radix of the base58 encoding system.
const Radix = len(Alphabet)

// Bits of entropy per base 54 digit.
var BitsPerDigit = math.Log2(float64(Radix))

var invAlphabet map[rune]*big.Int
var radixBig = big.NewInt(int64(Radix))

func init() {
    invAlphabet = make(map[rune]*big.Int, Radix)
    for index, value := range Alphabet {
        invAlphabet[value] = big.NewInt(int64(index))
    }
}

type CorruptInputError int

func (err CorruptInputError) Error() string {
    return fmt.Sprintf("illegal icechar data at input byte %d", err)
}

// DecodeInt returns the big.Int represented by the icechar string s.
func decodeInt(s string) (*big.Int, error) {
    n := new(big.Int)
    for index, digit := range s {
        n.Mul(n, radixBig)
        value, ok := invAlphabet[digit]
        if !ok {
          return nil, corruptInputError(index)
        }
        n.Add(n, value)
    }
    return n, nil
}

// Decode returns the bytes represented by the base58 string s.
func iceCharDecode(s string) ([]byte, error) {
    var zeros int
    for i := 0; i < len(s) && s[i] == Alphabet[0]; i++ {
        zeros++
    }
    n, err := decodeInt(s)
    if err != nil {
        return nil, err
    }
    return append(make([]byte, zeros), n.Bytes()...), nil
}

// EncodeInt encodes the big.Int n using icechar.
func encodeInt(n *big.Int) string {
    n = new(big.Int).Set(n)
    buf := make([]byte, 0, maxEncodedLen(n.BitLen()))
    remainder := new(big.Int)
    for n.Sign() == 1 {
        n.DivMod(n, radixBig, remainder)
        buf = append(buf, Alphabet[remainder.Int64()])
    }
    bufReverse := make([]byte, len(buf))
    for index, value := range buf {
        bufReverse[len(buf)-index-1] = value
    }
    return string(bufReverse)
}

// Encode encodes src using icechar.
func iceCharEncode(src []byte) string {
    var zeros int
    for i := 0; i < len(src) && src[i] == 0; i++ {
        zeros++
    }
    n := new(big.Int).SetBytes(src[zeros:])
    buf := append(make([]byte, zeros), encodeInt(n)...)
    for i := 0; i < zeros; i++ {
        buf[i] = Alphabet[0]
    }
    return string(buf)
}

// MaxEncodedLen returns the maximum length in bytes of an encoding of n source
// bits.
func maxEncodedLen(n int) int {
    return int(math.Ceil(float64(n) / BitsPerDigit))
}
