param(
    [Parameter(Mandatory)][string]$OutputDirectory,
    [Parameter(Mandatory)][string]$WinresPath,
    [string]$Version = '3.7.3.1'
)
$ErrorActionPreference = 'Stop'
$projectRoot = Split-Path $PSScriptRoot
$output = [IO.Path]::GetFullPath($OutputDirectory)
New-Item -ItemType Directory -Force -Path $output | Out-Null
Push-Location $projectRoot
try {
    $resource = Get-Content 'cmd/core/winres/winres.json' -Raw | ConvertFrom-Json
    Copy-Item (Join-Path $projectRoot 'cmd/core/winres/icon.png') (Join-Path $output 'icon.png')
    $resource.RT_GROUP_ICON.APP.'0000' = @('icon.png')
    $resource.RT_MANIFEST.'#1'.'0409'.identity.version = $Version
    $versionInfo = $resource.RT_VERSION.'#1'.'0000'
    $versionInfo.fixed.file_version = $Version
    $versionInfo.fixed.product_version = $Version
    $versionInfo.info.'0409'.FileVersion = $Version
    $versionInfo.info.'0409'.ProductVersion = $Version
    $resourcePath = Join-Path $output 'core-winres.json'
    $resource | ConvertTo-Json -Depth 20 | Set-Content $resourcePath -Encoding utf8NoBOM
    & $WinresPath make --arch amd64 --in $resourcePath --out cmd/core/rsrc
    if ($LASTEXITCODE) { throw 'Windows resource generation failed.' }
    $linker = "-s -w -H=windowsgui -X github.com/TIANLI0/THRM/internal/version.BuildVersion=$Version"
    go build -trimpath -ldflags $linker -o (Join-Path $output 'THRM Core.exe') ./cmd/core/
    if ($LASTEXITCODE) { throw 'Core repair build failed.' }
} finally { Pop-Location }
