# ==============================================================================
# SVN Pre-commit Hook - AI Code Review (PowerShell Version)
# ==============================================================================
# Description: 在提交代码前调用 AI 进行代码审查
#              如果评分 <= 6 分，阻止提交
#              评分 > 6 分，允许提交
# ==============================================================================

param(
    [Parameter(Mandatory=$true)]
    [string]$Repos,
    
    [Parameter(Mandatory=$true)]
    [string]$Txn
)

# 配置参数
$ReviewApiUrl = "http://localhost:8080/api/reviews/svn/pre-commit"
$RepoName = "my-repo"  # 修改为你的仓库名称
$Timeout = 30  # API 超时时间（秒）

# 颜色输出
$ErrorForegroundColor = 'Red'
$SuccessForegroundColor = 'Green'
$WarningForegroundColor = 'Yellow'

# ==============================================================================
# 辅助函数
# ==============================================================================

function Log-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message"
}

function Log-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $ErrorForegroundColor
}

function Log-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $SuccessForegroundColor
}

function Log-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $WarningForegroundColor
}

# ==============================================================================
# 主程序
# ==============================================================================

Log-Info "开始 AI 代码审查..."

# 获取 SVN 路径（假设 svnlook 在 PATH 中
$svnlook = "svnlook"
if (!(Get-Command $svnlook -ErrorAction SilentlyContinue) {
    $svnlook = Join-Path $env:VISUALSVN_SERVER "bin\svnlook.exe"
    if (!(Test-Path $svnlook)) {
        Log-Error "找不到 svnlook 命令"
        exit 0  # 允许提交，因为是配置问题
    }
}

try {
    # 获取提交信息
    $commitMsg = & $svnlook log -t $Txn $Repos 2>&1 | Out-String
    $author = & $svnlook author -t $Txn $Repos 2>&1
    $changed = & $svnlook changed -t $Txn $Repos 2>&1
    $diffContent = & $svnlook diff -t $Txn $Repos 2>&1 | Out-String
    
    # 检查是否有变更（空提交直接通过
    if ([string]::IsNullOrWhiteSpace($diffContent)) {
        Log-Success "没有代码变更，跳过审查"
        exit 0
    }

    # 检查是否有 "[skip-review]" 标记，如果有则跳过
    if ($commitMsg -match "\[skip-review\]") {
        Log-Success "检测到 [skip-review] 标记，跳过代码审查"
        exit 0
    }

    # 构建 JSON 数据
    $payload = @{
        repo = $RepoName
        diff = $diffContent
        author = $author
        message = $commitMsg
    }
    $jsonPayload = $payload | ConvertTo-Json -Depth 10

    Log-Info "调用代码审查 API (超时时间 ${Timeout}s)..."
    
    try {
        $response = Invoke-WebRequest -Uri $ReviewApiUrl -Method Post -Body $jsonPayload -ContentType "application/json" -TimeoutSec $Timeout
        $responseBody = $response.Content | ConvertFrom-Json
        $httpStatus = $response.StatusCode
    }
    catch [System.Net.WebException] {
        if ($_.Exception.Status -eq 'Timeout') {
            Log-Error "API 调用超时（${Timeout}秒）"
        } else {
            Log-Error "API 调用失败：$_"
        }
        Log-Error "代码审查服务异常，本次提交不做审查"
        exit 0  # 允许提交，因为是服务端问题
    }
    catch {
        Log-Error "API 调用失败：$_"
        Log-Error "代码审查服务异常，本次提交不做审查"
        exit 0  # 允许提交，因为是服务端问题
    }

    # 输出审查结果
    Write-Host ""
    Write-Host "========================================="
    Write-Host "           AI 代码审查结果"
    Write-Host "========================================="
    Write-Host ""

    if ($responseBody.allowed) {
        Log-Success $responseBody.message
        if ($responseBody.url) {
            Log-Info "查看详细报告：$($responseBody.url)"
        }
        Write-Host ""
        exit 0
    }
    else {
        Log-Error $responseBody.message
        Write-Host ""
        Write-Host "审查报告：" -ForegroundColor $WarningForegroundColor
        Write-Host "-----------------------------------------"
        Write-Host $responseBody.report
        Write-Host "-----------------------------------------"
        Write-Host ""
        if ($responseBody.url) {
            Log-Info "查看完整报告：$($responseBody.url)"
            Write-Host ""
        }
        Log-Error "代码未通过审查，请修复问题后重新提交！"
        Write-Host ""
        exit 1
    }
}
catch {
    Log-Error "发生错误：$_"
    Log-Error "代码审查服务异常，本次提交不做审查"
    exit 0  # 允许提交，因为是服务端问题
}
