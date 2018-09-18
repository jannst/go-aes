package matrix

import (
	"go-aes/math"
	s "go-aes/tables"
)

func (m *Matrix) Encrypt(keys []byte, rounds int) {
	if len(keys)%m.Size() != 0 {
		panic("m.size does not devide len(keys)")
	}

	m.AddRoundKey(keys[:m.Size()])

	for round := 1; round < rounds; round++ {
		m.SubBytes()
		m.ShiftRows()
		m.MixColumns()
		m.AddRoundKey(keys[round*m.Size() : (round+1)*m.Size()])
	}

	m.SubBytes()
	m.ShiftRows()
	m.AddRoundKey(keys[rounds*m.Size() : (rounds+1)*m.Size()])
}

func (m *Matrix) SubBytes() {
	s.SubWord(m.data)
}

func (m *Matrix) ShiftRows() {
	buf := make([]byte, len(m.data))
	var dst int
	copy(buf, m.data)

	for row := 1; row < m.height; row++ {
		for col := 0; col < m.nk; col++ {
			dst = math.Modulo(col-row, m.nk)
			m.data[dst*m.height+row] = buf[col*m.height+row]
		}
	}
}

func (m *Matrix) MixColumns() {
	buf := [4]byte{}
	for i := 0; i < m.Size(); i += 4 {
		buf[0] = s.Mul2(m.data[i]) ^ s.Mul3(m.data[i+1]) ^ m.data[i+2] ^ m.data[i+3]
		buf[1] = m.data[i] ^ s.Mul2(m.data[i+1]) ^ s.Mul3(m.data[i+2]) ^ m.data[i+3]
		buf[2] = m.data[i] ^ m.data[i+1] ^ s.Mul2(m.data[i+2]) ^ s.Mul3(m.data[i+3])
		buf[3] = s.Mul3(m.data[i]) ^ m.data[i+1] ^ m.data[i+2] ^ s.Mul2(m.data[i+3])
		m.data[i] = buf[0]
		m.data[i+1] = buf[1]
		m.data[i+2] = buf[2]
		m.data[i+3] = buf[3]
	}
}

func (m *Matrix) AddRoundKey(key []byte) {
	if m.Size() != len(key) {
		panic("input is not of same size as matrix")
	}

	for i := 0; i < len(m.data); i++ {
		m.data[i] ^= key[i]
	}
}
