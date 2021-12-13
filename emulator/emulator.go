package emulator

import (
    "fmt"
)

const (
    EAX = iota
    ECX
    EDX
    EBX
    ESP
    EBP
    ESI
    EDI
    REGISTERS_COUNT
)

var registers_name = [...] string{"EAX", "ECX", "EDX", "EBX", "ESP", "EBP", "ESI", "EDI"}

type Emulator struct {
    registers [REGISTERS_COUNT]uint32 // 汎用レジスタ
    eflags uint32                     // EFLAGSレジスタ
    memory []uint8                    // メモリ（バイト列）
    eip uint32                        // プログラムカウンタ
    memory_size uint32                // メモリサイズ
}


func NewEmulator(memory_size uint32, eip uint32, esp uint32) *Emulator {
    emu := &Emulator{eip: eip, memory: make([]uint8, memory_size)}
    emu.memory_size = memory_size
    emu.registers[ESP] = esp

    return emu
}

func (emu *Emulator) Exec(program []byte, load_address uint) {
    copy(emu.memory[load_address:], program)

    for emu.eip < emu.memory_size {
        code := uint8(emu.getCode8(0))
        fmt.Printf("EIP = %X, Code = %X\n", emu.eip, code)
        emu.callInstruction(code)

        if emu.eip == 0 {
            fmt.Println("End of program")
            break
        }
    }

    emu.dumpRegisters()
}

func (emu *Emulator) dumpRegisters() {
    for i := 0; i < REGISTERS_COUNT; i++ {
        fmt.Printf("%s = %08x\n", registers_name[i], emu.registers[i])
    }
}

func (emu *Emulator) movR32Imm32() {
    reg := emu.getCode8(0) - 0xB8
    value := emu.getCode32(1)
    emu.registers[reg] = value
    emu.eip += 5
}

func (emu *Emulator) shortJump() {
    diff := emu.getSignCode8(1)
    emu.eip = emu.eip + uint32(diff) + 2
}

func (emu *Emulator) nearJump() {
    diff := emu.getSignCode32(1)
    emu.eip = emu.eip + uint32(diff) + 5
}

func (emu *Emulator) getCode8(index int) uint8 {
    return emu.memory[int(emu.eip) + index]
}

func (emu *Emulator) getSignCode8(index int) int8 {
    return int8(emu.memory[int(emu.eip) + index])
}

func (emu *Emulator) getCode32(index int) uint32 {
    var ret uint32 = 0

    for i := 0; i < 4; i++ {
        code := emu.getCode8(index + i)
        ret |= uint32(code) << (i * 8)
    }

    return ret
}

func (emu *Emulator) getSignCode32(index int) int32 {
    return int32(emu.getCode32(index))
}

func (emu *Emulator) callInstruction(code uint8) {
    switch {
    case 0xB8 <= code && code <= 0xB8 + REGISTERS_COUNT - 1:
        emu.movR32Imm32()
    case code == 0xE9:
        emu.nearJump()
    case code == 0xEB:
        emu.shortJump()
    }
}
