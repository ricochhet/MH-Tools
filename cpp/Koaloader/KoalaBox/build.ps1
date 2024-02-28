# Usage:
#   build.ps1 32 Debug
#   build.ps1 32 Release
#   build.ps1 64 Debug
#   build.ps1 64 Release

$target = $args[0]

$arch = $args[1]
if ($arch -notmatch '^(32|64)$')
{
    throw "Invalid architecute. Expected: '32' or '64'. Got: '$arch'"
}
$platform = If ($arch -eq '32')
{
    'Win32'
}
Else
{
    'x64'
}

$config = $args[2]
if ($config -notmatch '^(Debug|Release|RelWithDebInfo)$')
{
    throw "Invalid architecute. Expected: 'Debug' or 'Release' or 'RelWithDebInfo'. Got: '$config'"
}

$Env:BUILD_DIR = "build\$arch"

function Build-Project {
    cmake -G "Visual Studio 17 2022" -A $platform -B "$Env:BUILD_DIR" "$Env:CMAKE_OPTIONS" "-DMODULE=$Env:MODULE"

    cmake --build "$Env:BUILD_DIR" --target $target --config $config
}
