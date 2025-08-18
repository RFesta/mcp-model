param(
  [Parameter(Mandatory=$true)][string]$NewName
)

$Old = "modelo-mcp"

# Substitui conteúdo
Get-ChildItem -Recurse -File | ForEach-Object {
  (Get-Content $_.FullName) -replace $Old, $NewName | Set-Content $_.FullName
}

# Renomeia pasta cmd
if (Test-Path "cmd\$Old") { Rename-Item "cmd\$Old" $NewName }

Write-Host "OK: base renomeada para $NewName"
