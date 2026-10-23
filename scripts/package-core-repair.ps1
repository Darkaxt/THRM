param(
    [Parameter(Mandatory)][string]$OutputDirectory,
    [Parameter(Mandatory)][string]$WinresPath
)
$ErrorActionPreference = 'Stop'
$root = Split-Path $PSScriptRoot
$output = [IO.Path]::GetFullPath($OutputDirectory)
if (Test-Path $output) { throw 'Use a new output directory.' }
New-Item -ItemType Directory $output | Out-Null
$revision = git -C $root rev-parse HEAD
if ($LASTEXITCODE) { throw 'Cannot resolve source revision.' }
$upstreamUrl = 'https://github.com/TIANLI0/THRM/releases/download/v3.7.3/THRM-windows-portable.zip'
$upstreamHash = '1e8a7d1f25cc702260b4fed61abd83d0d68277f0c79bd6e7ba957fa067175c21'
$download = Join-Path $output 'upstream.zip'
Invoke-WebRequest $upstreamUrl -OutFile $download
if ((Get-FileHash $download -Algorithm SHA256).Hash -ne $upstreamHash) { throw 'Upstream digest mismatch.' }
$unpacked = Join-Path $output 'upstream'
Expand-Archive $download $unpacked
$gui = @(Get-ChildItem $unpacked -Recurse -Filter THRM.exe)
if ($gui.Count -ne 1) { throw 'Unexpected upstream portable layout.' }
$package = Join-Path $output 'THRM-3.7.3.1-windows-x64'
New-Item -ItemType Directory $package | Out-Null
Copy-Item $gui[0].FullName $package
Copy-Item (Join-Path $gui[0].DirectoryName 'bridge') $package -Recurse
& (Join-Path $PSScriptRoot 'build-core-repair.ps1') -OutputDirectory (Join-Path $output 'core') -WinresPath $WinresPath
Copy-Item (Join-Path $output 'core/THRM Core.exe') $package
Copy-Item (Join-Path $root 'LICENSE') $package
Copy-Item (Join-Path $root 'docs/GUI-REOPEN-REPAIR.md') (Join-Path $package 'REPAIR.md')
$files = @(Get-ChildItem $package -File -Recurse | ForEach-Object {
    @{path=[IO.Path]::GetRelativePath($package,$_.FullName); sha256=(Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLowerInvariant()}
})
@{version='3.7.3.1'; coreSource=$revision; upstreamSource='2131102526c8772211896ce23dd24a93065a37d3'; upstreamVersion='3.7.3'; upstreamUrl=$upstreamUrl; upstreamSha256=$upstreamHash; files=$files} |
    ConvertTo-Json -Depth 8 | Set-Content (Join-Path $package 'provenance.json') -Encoding utf8NoBOM
$archive = Join-Path $output 'THRM-3.7.3.1-windows-x64.zip'
Compress-Archive -Path (Join-Path $package '*') -DestinationPath $archive -CompressionLevel Optimal
"$((Get-FileHash $archive -Algorithm SHA256).Hash.ToLowerInvariant())  $([IO.Path]::GetFileName($archive))" |
    Set-Content (Join-Path $output 'SHA256SUMS.txt') -Encoding ascii
Get-Item $archive
