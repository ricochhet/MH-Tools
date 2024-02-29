windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.syso
Move-Item -Path .\MHWArchiveManager.syso -Destination .\cmd\main -Force

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.Cli.syso
Move-Item -Path .\MHWArchiveManager.Cli.syso -Destination .\cmd\cli_only -Force

windres MHWArchiveManager.rc -O coff -o MHWArchiveManager.Gui.syso
Move-Item -Path .\MHWArchiveManager.Gui.syso -Destination .\cmd\gui_only -Force