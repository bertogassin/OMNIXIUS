# Загрузка сайта OMNIXIUS на OVH по FTP
# Запуск: правый клик -> "Выполнить с PowerShell" или: powershell -ExecutionPolicy Bypass -File upload-ovh.ps1

# ========== ЗАПОЛНИ СВОИ ДАННЫЕ (из панели OVH -> хостинг -> FTP) ==========
$FTP_HOST = "ftp.cluster0XX.hosting.ovh.net"   # например ftp.cluster032.hosting.ovh.net
$FTP_USER = "твой_логин_ftp"
$FTP_PASS = "твой_пароль_ftp"
# ============================================================================

$localRoot = $PSScriptRoot
$remotePath = "/www"

function Upload-FtpDirectory {
    param([string]$localDir, [string]$remoteDir)
    $items = Get-ChildItem $localDir -Force
    foreach ($item in $items) {
        $remoteItem = "$remoteDir/$($item.Name)"
        if ($item.PSIsContainer) {
            if ($item.Name -eq ".git" -or $item.Name -eq "node_modules") { continue }
            try {
                $makeDir = [System.Net.WebRequest]::Create("ftp://$FTP_HOST$remoteItem")
                $makeDir.Method = [System.Net.WebRequestMethods+Ftp]::MakeDirectory
                $makeDir.Credentials = New-Object System.Net.NetworkCredential($FTP_USER, $FTP_PASS)
                $makeDir.GetResponse() | Out-Null
            } catch { }
            Upload-FtpDirectory -localDir $item.FullName -remoteDir $remoteItem
        } else {
            $uri = [System.Uri]"ftp://$FTP_HOST$remoteItem"
            $request = [System.Net.FtpWebRequest]::Create($uri)
            $request.Method = [System.Net.WebRequestMethods+Ftp]::UploadFile
            $request.Credentials = New-Object System.Net.NetworkCredential($FTP_USER, $FTP_PASS)
            $request.UseBinary = $true
            $request.UsePassive = $true
            $fileContent = [System.IO.File]::ReadAllBytes($item.FullName)
            $request.ContentLength = $fileContent.Length
            $requestStream = $request.GetRequestStream()
            $requestStream.Write($fileContent, 0, $fileContent.Length)
            $requestStream.Close()
            $response = $request.GetResponse()
            Write-Host "  OK: $($item.Name)"
            $response.Close()
        }
    }
}

Write-Host "OMNIXIUS -> OVH FTP" -ForegroundColor Cyan
Write-Host "Host: $FTP_HOST" -ForegroundColor Gray
if ($FTP_USER -match "твой_логин" -or $FTP_PASS -match "твой_пароль") {
    Write-Host ""
    Write-Host "Ошибка: открой upload-ovh.ps1 и впиши свои FTP логин и пароль из панели OVH." -ForegroundColor Yellow
    Write-Host "Панель OVH -> Web Cloud -> Хостинги -> твой хостинг -> FTP-SSH." -ForegroundColor Gray
    pause
    exit 1
}

Write-Host "Загружаю файлы..." -ForegroundColor Green
# Файлы в корне
foreach ($f in @("index.html", "contact.html", "404.html")) {
    $path = Join-Path $localRoot $f
    if (Test-Path $path) {
        $uri = [System.Uri]"ftp://$FTP_HOST$remotePath/$f"
        $request = [System.Net.FtpWebRequest]::Create($uri)
        $request.Method = [System.Net.WebRequestMethods+Ftp]::UploadFile
        $request.Credentials = New-Object System.Net.NetworkCredential($FTP_USER, $FTP_PASS)
        $request.UseBinary = $true
        $request.UsePassive = $true
        $fileContent = [System.IO.File]::ReadAllBytes($path)
        $request.ContentLength = $fileContent.Length
        $requestStream = $request.GetRequestStream()
        $requestStream.Write($fileContent, 0, $fileContent.Length)
        $requestStream.Close()
        $request.GetResponse() | Out-Null
        Write-Host "  OK: $f"
    }
}
# Папки css и js
Upload-FtpDirectory -localDir (Join-Path $localRoot "css") -remoteDir "$remotePath/css"
Upload-FtpDirectory -localDir (Join-Path $localRoot "js")  -remoteDir "$remotePath/js"

Write-Host ""
Write-Host "Готово. Сайт загружен на хостинг OVH." -ForegroundColor Green
Write-Host "Открой свой домен в браузере." -ForegroundColor Gray
pause
