param(
    [string]$dll = "version"
)

cd ../Koaloader/
.\build.ps1 64 Release $dll
cmake --build "./build/64/$dll" --target Koaloader --config Release