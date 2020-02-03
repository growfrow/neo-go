package emit

import (
	"encoding/binary"
	"testing"

	"github.com/CityOfZion/neo-go/pkg/io"
	"github.com/CityOfZion/neo-go/pkg/vm/opcode"
	"github.com/stretchr/testify/assert"
)

func TestEmitInt(t *testing.T) {
	t.Run("1-byte int", func(t *testing.T) {
		buf := io.NewBufBinWriter()
		Int(buf.BinWriter, 10)
		result := buf.Bytes()
		assert.EqualValues(t, opcode.PUSH10, result[0])
	})

	t.Run("2-byte int", func(t *testing.T) {
		buf := io.NewBufBinWriter()
		Int(buf.BinWriter, 100)
		result := buf.Bytes()
		assert.EqualValues(t, opcode.PUSHBYTES1, result[0])
		assert.EqualValues(t, 100, result[1])
	})

	t.Run("4-byte int", func(t *testing.T) {
		buf := io.NewBufBinWriter()
		Int(buf.BinWriter, 1000)
		result := buf.Bytes()
		assert.EqualValues(t, opcode.PUSHBYTES2, result[0])
		assert.EqualValues(t, 1000, binary.LittleEndian.Uint16(result[1:3]))
	})
}

func TestEmitBool(t *testing.T) {
	buf := io.NewBufBinWriter()
	Bool(buf.BinWriter, true)
	Bool(buf.BinWriter, false)
	result := buf.Bytes()
	assert.Equal(t, opcode.Opcode(result[0]), opcode.PUSH1)
	assert.Equal(t, opcode.Opcode(result[1]), opcode.PUSH0)
}

func TestEmitString(t *testing.T) {
	buf := io.NewBufBinWriter()
	str := "City Of Zion"
	String(buf.BinWriter, str)
	assert.Equal(t, buf.Len(), len(str)+1)
	assert.Equal(t, buf.Bytes()[1:], []byte(str))
}

func TestEmitSyscall(t *testing.T) {
	syscalls := []string{
		"Neo.Runtime.Log",
		"Neo.Runtime.Notify",
		"Neo.Runtime.Whatever",
	}

	buf := io.NewBufBinWriter()
	for _, syscall := range syscalls {
		Syscall(buf.BinWriter, syscall)
		result := buf.Bytes()
		assert.Equal(t, opcode.Opcode(result[0]), opcode.SYSCALL)
		assert.Equal(t, result[1], uint8(len(syscall)))
		assert.Equal(t, result[2:], []byte(syscall))
		buf.Reset()
	}
}

func TestEmitCall(t *testing.T) {
	buf := io.NewBufBinWriter()
	Call(buf.BinWriter, opcode.JMP, 100)
	result := buf.Bytes()
	assert.Equal(t, opcode.Opcode(result[0]), opcode.JMP)
	label := binary.LittleEndian.Uint16(result[1:3])
	assert.Equal(t, label, uint16(100))
}
