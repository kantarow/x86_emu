package main

import (
    "fmt"
    "io/ioutil"
    "os"
    "x86_emu/emulator"
)

const MEMORY_SIZE = (1024 * 1024)

func main() {
    emu := emulator.NewEmulator(MEMORY_SIZE, 0x7c00, 0x7c00)

    if len(os.Args) != 2 {
        fmt.Println("usage: x86_emu filename\n")
        os.Exit(1)
    }

    buf, err := ioutil.ReadFile(os.Args[1])

    if err != nil {
        fmt.Println(os.Stderr, err)
        os.Exit(1)
    }


    emu.Exec(buf, 0x7c00)
}
