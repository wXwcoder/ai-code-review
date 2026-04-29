@echo off
REM ==============================================================================
REM SVN Pre-commit Hook - AI Code Review (Windows Batch)
REM ==============================================================================

REM 设置 PowerShell 脚本路径（和批处理文件同目录
set SCRIPT_DIR=%~dp0
set PS_SCRIPT=%SCRIPT_DIR%pre-commit.ps1

REM 参数
set REPOS=%1
set TXN=%2

REM 检查 PowerShell 是否可用
where powershell >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] PowerShell 不可用
    exit 0
)

REM 调用 PowerShell 脚本
powershell -ExecutionPolicy Bypass -File "%PS_SCRIPT%" -Repos "%REPOS%" -Txn "%TXN%"

REM 退出码
if %ERRORLEVEL% EQU 0 (
    exit 0
) else (
    exit 1
)
