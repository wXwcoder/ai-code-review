@echo off
setlocal
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
set GO111MODULE=on
set VERSION=1

cd web
call npm run build
cd ..


echo Compiling Windows version of AICodeReview...
go build -ldflags="-s -w -X main.VERSION=0.0.%VERSION%" -o AICodeReview.exe
if %errorlevel% neq 0 (
    echo Compilation failed!
    exit /b %errorlevel%
)
echo Compilation completed, output file: AICodeReview.exe
set ZIP_DIR=AICodeReview.%VERSION%
if exist %ZIP_DIR% rmdir /s /q %ZIP_DIR%
mkdir %ZIP_DIR%
copy AICodeReview.exe %ZIP_DIR%\ >nul
copy start.bat %ZIP_DIR%\ >nul
xcopy dist %ZIP_DIR%\dist /E /I /Y >nul

if not exist output mkdir output
if exist output\AICodeReview.%VERSION%.zip del output\AICodeReview.%VERSION%.zip
powershell -Command "Compress-Archive -Path '%ZIP_DIR%' -DestinationPath 'output\AICodeReview.%VERSION%.zip' -CompressionLevel Optimal" >nul
rmdir /s /q %ZIP_DIR%
endlocal