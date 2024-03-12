New-Item -ItemType Directory -Path build -Force

Write-Host "Building MHWArchiveManager"
Set-Location "cmd\main"
go build -o "..\..\build\MHWArchiveManager.exe" -ldflags="-s -w -H=windowsgui -extldflags=-static" .

Write-Host "Building MHWArchiveManager.Gui"
Set-Location "..\gui_only"
go build -o "..\..\build\MHWArchiveManager.Gui.exe" -ldflags="-s -w -H=windowsgui -extldflags=-static" .

Write-Host "Building MHWArchiveManager.Cli"
Set-Location "..\cli_only"
go build -o "..\..\build\MHWArchiveManager.Cli.exe" -ldflags="-s -w -extldflags=-static" .

Set-Location "..\..\"

Copy-item -Path ".\scripts\MHW-BuildEffectRemovalMod.ps1" -Destination ".\build\"
Copy-item -Path ".\scripts\MHW-BuildEffectRemovalMod2.ps1" -Destination ".\build\"
