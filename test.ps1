$GO_FILE = "test.go"
Write-Output "EXECUTEI"

$NUM_EXEC = 3
$processes = @()

for ($i = 1; $i -le $NUM_EXEC; $i++) {
    Write-Output "Running $i"
    $process = Start-Process -NoNewWindow -FilePath "go" -ArgumentList "run", $GO_FILE -PassThru
    $processes += $process
}

foreach ($process in $processes) {
    $process.WaitForExit()
}

Write-Output "Done"