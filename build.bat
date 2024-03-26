@echo off
setlocal EnableExtensions

set GOOS=linux
go build -o ProxyJuice .

if "%errorlevel%" neq "0" (
    echo ERROR: Linux build failed!
    goto :eof
)

echo Builded ProxyJuice
echo Compressing ProxyJuice

if "%1" == "final" (
    upx.exe --lzma ProxyJuice
    echo Compressed ProxyJuice
)

echo Done
