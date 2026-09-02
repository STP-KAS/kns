$ErrorActionPreference = "Stop"
$silverc = "C:\Users\Remco\tools\silverc\silverc.exe"
if (-not (Test-Path $silverc)) { throw "missing official silverc.exe" }
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
& $silverc "$root\contracts\v1\KasName.sil" --constructor-args "$root\contracts\v1\ctor-KasName.json" -o "$root\contracts\v1\KasName.json"
& $silverc "$root\contracts\v1\KaChatPayTimeout.sil" --constructor-args "$root\contracts\v1\ctor-KaChatPayTimeout.json" -o "$root\contracts\v1\KaChatPayTimeout.json"
& $silverc "$root\contracts\v1\WorkCredit.sil" --constructor-args "$root\contracts\v1\ctor-WorkCredit.json" -o "$root\contracts\v1\WorkCredit.json"
Write-Output "compiled KasName + KaChatPayTimeout + WorkCredit"
