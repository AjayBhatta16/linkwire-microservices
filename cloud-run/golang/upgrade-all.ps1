$functionsDir = Join-Path $PSScriptRoot "functions"

Get-ChildItem -Path $functionsDir -Directory | ForEach-Object {
    $folderName = $_.Name
    $folderPath = $_.FullName

    Write-Host "Upgrading dependencies for function: $folderName"
    
    cd $folderPath
    go get github.com/AjayBhatta16/linkwire-golang-shared@latest
    go mod tidy
}

cd ..