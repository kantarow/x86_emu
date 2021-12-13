package main

import (
    "fmt"
    "io/ioutil"
    "os"
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

const MEMORY_SIZE = (1024 * 1024)

type Emulator struct {
    registers [REGISTERS_COUNT]uint32 // 汎用レジスタ
    eflags uint32                     // EFLAGSレジスタ
    memory []uint8                    // メモリ（バイト列）
    eip uint32                        // プログラムカウンタ
}

func main() {
    emu := create_emu(MEMORY_SIZE, 0x0000, 0x7c00)

    if len(os.Args) != 2 {
        fmt.Println("usage: x86_emu filename\n")
        os.Exit(1)
    }

    buf, err := ioutil.ReadFile(os.Args[1])

    if err != nil {
        fmt.Println(os.Stderr, err)
        os.Exit(1)
    }

    copy(emu.memory[:], buf)

    for emu.eip < MEMORY_SIZE {
        code := uint8(emu.get_code8(0))
        fmt.Printf("EIP = %X, Code = %X\n", emu.eip, code)
        emu.call_instruction(code)

        if emu.eip == 0x00 {
            fmt.Println("End of program")
            break
        }
    }

    emu.dump_registers()
}

func create_emu(mem_size int, eip uint32, esp uint32) *Emulator {
    emu := &Emulator{eip: eip, memory: make([]uint8, mem_size)}
    emu.registers[ESP] = esp

    return emu
}

func (emu *Emulator) dump_registers() {
    for i := 0; i < REGISTERS_COUNT; i++ {
        fmt.Printf("%s = %08x\n", registers_name[i], emu.registers[i])
    }
}

func (emu *Emulator) mov_r32_imm32() {
    reg := emu.get_code8(0) - 0xB8
    value := emu.get_code32(1)
    emu.registers[reg] = value
    emu.eip += 5
}

func (emu *Emulator) short_jump() {
    diff := emu.get_sign_code8(1)
    emu.eip = emu.eip + uint32(diff) + 2
}

func (emu *Emulator) get_code8(index int) uint8 {
    return emu.memory[int(emu.eip) + index]
}

func (emu *Emulator) get_sign_code8(index int) int8 {
    return int8(emu.memory[int(emu.eip) + index])
}

func (emu *Emulator) get_code32(index int) uint32 {
    var ret uint32 = 0

    for i := 0; i < 4; i++ {
        ret |= uint32(emu.get_code8(index + i) << (i * 8))
    }

    return ret
}

func (emu *Emulator) call_instruction(code uint8) {
    switch {
    case 0xB8 <= code && code <= 0xB8 + REGISTERS_COUNT - 1:
        emu.mov_r32_imm32()
    case code == 0xEB:
        emu.short_jump()
    }
}
