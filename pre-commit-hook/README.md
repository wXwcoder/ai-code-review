# SVN Pre-commit Hook - AI 代码审查

## 功能介绍

这是一个 SVN pre-commit 钩子，使用 AI 在代码审查工具，根据评分决定是否允许提交。

- **评分规则**：
  - 评分 > 6 分：允许提交
  - 评分 <= 6 分：阻止提交

## 安装配置

### 1. 部署代码审查服务

首先确保你已经部署了 `git-svn-reviewbot-go` 服务，并且服务正在运行。

### 2. 超时配置

钩子脚本中已经配置了超时时间，默认 30 秒：

- Linux/MacOS: 编辑 `pre-commit` 中的 `TIMEOUT=30`
- Windows: 编辑 `pre-commit.ps1` 中的 `$Timeout = 30`

**注意**：超时情况下会允许提交，不阻止开发者工作。

### 3. 复制钩子到 SVN 仓库

#### Linux/MacOS

在 SVN 仓库的 `hooks` 目录中：

```bash
# 1. 进入仓库 hooks 目录
cd /path/to/svn/repo/hooks

# 2. 复制 pre-commit 脚本
cp /path/to/git-svn-reviewbot-go/pre-commit-hook/pre-commit ./pre-commit

# 3. 设置执行权限
chmod +x pre-commit

# 4. 编辑脚本，修改配置
vim pre-commit
```

编辑脚本中的配置：
```bash
REVIEW_API_URL="http://localhost:8080/api/reviews/svn/pre-commit"
REPO_NAME="my-repo"  # 修改为你的仓库名称
```

#### Windows

```batch
# 1. 进入仓库 hooks 目录
cd C:\path\to\svn\repo\hooks

# 2. 复制脚本
copy C:\path\to\git-svn-reviewbot-go\pre-commit-hook\pre-commit.bat pre-commit.bat
copy C:\path\to\git-svn-reviewbot-go\pre-commit-hook\pre-commit.ps1 pre-commit.ps1

# 3. 编辑 pre-commit.ps1 配置
notepad pre-commit.ps1
```

编辑 PowerShell 脚本中的配置：
```powershell
$ReviewApiUrl = "http://localhost:8080/api/reviews/svn/pre-commit"
$RepoName = "my-repo"  # 修改为你的仓库名称
```

## 使用方法

### 正常提交（会自动审查）

```bash
# 正常提交，AI 会自动审查
svn commit -m "修改用户登录逻辑"
```

### 跳过审查

在提交信息中添加 `[skip-review]` 可以跳过代码审查：

```bash
svn commit -m "[skip-review] 紧急修复，临时提交"
```

## API 文档

### POST /api/reviews/svn/pre-commit

请求 Body：
```json
{
  "repo": "my-repo",
  "diff": "svn diff 内容",
  "author": "username",
  "message": "commit message"
}
```

响应：
```json
{
  "allowed": true,
  "score": 8,
  "report": "完整的 Markdown 报告",
  "url": "http://localhost:8080/api/reviews/svn?repo=xxx&revision=xxx",
  "message": "代码审查通过！评分：8/10 分"
}
```

## 评分标准

| 分数区间 | 含义 |
|---------|-----|
| 10 | 完美，没有任何问题 |
| 8-9 | 优秀，只有很小的改进空间 |
| 6-7 | 合格，有一些问题但可以接受 |
| 4-5 | 较差，有明显问题需要修改 |
| 0-3 | 极差，必须重写 |

## 超时和故障处理

### 容错设计理念

钩子脚本采用**乐观容错**的设计理念：

| 场景 | 行为 |
|------|------|
| API 超时（默认 30s） | ✅ **安全策略**：允许提交（因为是服务端问题，不阻止工作） |

### 超时配置

默认超时时间：`30秒`。

如需修改超时时间：

#### Linux/MacOS:
```bash
# 修改 TIMEOUT 变量
TIMEOUT=60  # 改为 60秒
```

#### Windows:
```powershell
$Timeout = 60  # 改为 60秒
```

### 错误处理机制

钩子会在以下场景**允许提交**（exit 0）：

1. API 超时
2. 网络错误（HTTP != 200）
3. 服务不可用
4. 任何服务端异常

### 失败安全策略解释

这样做的好处：

- **不阻塞开发者正常工作**：即使服务暂时不可用不影响提交
- **故障恢复后，可以事后审查**：SVN 提交后，可事后查看审查结果
- **审查可以事后通过 UI 查看**

## 注意事项

1. **依赖要求**：

- 确保 SVN 服务器上需要安装 `jq` (仅 Linux/MacOS)
- 确保能访问代码审查服务的 API
- 如果 API 调用失败时，钩子会允许提交（因为是服务端问题）
- 可以使用 `[skip-review]` 标签跳过审查

2. **故障排查**：

```bash
# 检查脚本是否有执行权限
ls -l pre-commit

# 检查是否有足够的超时时间
# 审查时间可能比较长，建议配置合理的超时
```

3. **测试钩子**：

```bash
# 可以手动运行脚本测试
cd /path/to/svn/repo/hooks
./pre-commit /path/to/svn/repo txn-id
```

4. **VisualSVN Server**：

确保 VisualSVN Server 配置：

1. 打开 VisualSVN Server Manager
2. 选择仓库 -> Properties -> Hooks
3. 添加 Pre-commit hook
4. 指向你的脚本文件

## 安全提示

- 建议同时部署多个仓库配置，AI 会严格，避免被人跳过审查
- 建议，
- 建议

**生产环境部署**
- 生产环境中，建议使用 HTTPS 和
