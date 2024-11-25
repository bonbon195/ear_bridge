[На русском](README_ru.md)

# EarBridge

Winner of the "White Noise" exhibition 2024 in Sverdlovsk Oblast

## About the Program

This is a program for listening to audio on multiple headphones simultaneously. Written in Go, it uses the Gio for a GUI. For working with audio devices, it uses [malgo](https://github.com/gen2brain/malgo).

## How to Build

### Windows
Run the script in Powershell with the command `./scripts/build.ps1`

### Other Platforms
1. Install [Go](https://go.dev/)
2. Install gogio with `go install gioui.org/cmd/gogio@latest`
3. Compile with `gogio -target <your_OS> -o build/<filename> .`

