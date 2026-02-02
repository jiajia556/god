# 获取当前目录下所有子目录
$subDirectories = Get-ChildItem -Directory

# 遍历每个子目录
foreach ($dir in $subDirectories) {
    Write-Host "处理目录: $($dir.FullName)"
    
    # 获取该目录中的所有文件
    $files = Get-ChildItem -Path $dir.FullName -File
    
    # 遍历目录中的每个文件
    foreach ($file in $files) {
        # 构建新的文件名
        $newName = $file.BaseName + ".go.tmpl"
        $newPath = Join-Path $file.Directory.FullName $newName

        # 重命名文件
        try {
            Rename-Item -Path $file.FullName -NewName $newName
            Write-Host "  重命名: $($file.Name) -> $newName"
        }
        catch {
            Write-Host "  错误: 无法重命名 $($file.Name) - $($_.Exception.Message)" -ForegroundColor Red
        }
    }
}

Write-Host "处理完成!" -ForegroundColor Green