# thank you b7
$Name="autotyper"
$zipFilePath = ".\bin\extension.zip"

echo "Building for Linux..."
$env:GOOS="linux"; go build -o bin/${Name}-linux .
echo "Building for Mac..."
$env:GOOS="darwin"; go build -o bin/${Name}-mac .
echo "Building for Windows..."
$env:GOOS="windows"; go build -o bin/${Name}-win.exe -ldflags="-H=windowsgui" .
echo "Build complete."


# Check if the zip file exists and delete it if it does
if (Test-Path $zipFilePath) {
    echo "Removing existing extension.zip..."
    Remove-Item $zipFilePath
}

echo "Creating extension.zip..."
Compress-Archive -Path .\bin\* -DestinationPath $zipFilePath
echo "Done"