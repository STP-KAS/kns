# Bring project delusional local sites back up.
# KNS :8080  Gramlane :8081  Kaspa Till :8082
$ErrorActionPreference = "Stop"
Get-Process kns,gramlane,kastill -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Milliseconds 400

Start-Process -FilePath "C:\Users\Remco\kns\kns.exe" -WorkingDirectory "C:\Users\Remco\kns" -WindowStyle Minimized
Start-Process -FilePath "C:\Users\Remco\Documents\kaspa\superapp\gramlane.exe" -WorkingDirectory "C:\Users\Remco\Documents\kaspa\superapp" -WindowStyle Minimized
Start-Process -FilePath "C:\Users\Remco\Documents\kaspa\superappstablesalternative\kastill.exe" -WorkingDirectory "C:\Users\Remco\Documents\kaspa\superappstablesalternative" -WindowStyle Minimized
Start-Sleep 2

function Ping($url, $name) {
  try {
    $r = Invoke-WebRequest $url -UseBasicParsing -TimeoutSec 4
    Write-Output ("OK  {0}  {1}" -f $name, $url)
  } catch {
    Write-Output ("DOWN {0}  {1}" -f $name, $url)
  }
}
Ping "http://localhost:8080/" "KNS"
Ping "http://localhost:8081/" "Gramlane"
Ping "http://localhost:8082/" "Kaspa Till"
Ping "http://localhost:8080/234" "KNS #234"
Ping "http://localhost:8081/234" "Gramlane #234"
Ping "http://localhost:8082/234" "Till #234"
Write-Output "If a line says DOWN, run the .exe in that folder once and look at the console error."
