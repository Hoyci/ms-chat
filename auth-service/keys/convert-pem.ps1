# Script para converter todos os arquivos .pem em base64 (.b64)
Get-ChildItem -Path . -Filter "*.pem" | ForEach-Object {
    $inputPath = $_.FullName
    $outputPath = "$($inputPath).b64"
    
    Write-Host "Convertendo: $($_.Name) -> $(Split-Path $outputPath -Leaf)"

    $base64 = [Convert]::ToBase64String([IO.File]::ReadAllBytes($inputPath))
    $base64 | Out-File -Encoding ascii $outputPath
}

Write-Host "`n✅ Conversão concluída!"