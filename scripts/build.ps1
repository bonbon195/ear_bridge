if (!(Get-Command -Name "gogio" -ErrorAction SilentlyContinue)) {
    go install gioui.org/cmd/gogio@latest
}
gogio -target windows -o build/ear_bridge.exe .
Remove-Item *.syso