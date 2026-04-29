# AI 代码审查机器人

这是一个使用 Go 语言开发的 AI 代码审查机器人，用于在闭网环境中提供智能代码审查服务。

## 功能特性

- ✅ GitLab Merge Request 自动审查
- ✅ SVN Commit 自动审查（通过关键词 [review] 触发）
- ✅ **SVN 手动触发审查**（通过 API 手动触发）
- ✅ **SVN Pre-commit 钩子审查**（带评分，阻止低分提交）
- ✅ **审查历史列表 Web UI**（左侧列表，右侧详情）
- ✅ **分享链接支持**（/?version=xxx 直接查看特定报告）
- ✅ **Ollama 本地 LLM 集成** (默认)
- ✅ **DeepSeek API 支持** (OpenAI 兼容格式)
- ✅ 环境变量配置
- ✅ Docker 部署支持

## 技术栈

- **Web 框架**: Gin
- **配置管理**: Viper
- **HTTP 客户端**: 标准库
- **SVN 操作**: 命令行调用（需要安装 SVN 客户端）

## 项目结构

```
git-svn-reviewbot-go/
├── main.go                  # 主程序入口
├── internal/
│   ├── config/
│   │   └── config.go        # 配置管理
│   ├── controller/
│   │   ├── gitlab.go        # GitLab Webhook 控制器
│   │   ├── svn.go           # SVN Webhook 控制器（含 pre-commit）
│   │   └── health.go        # 健康检查
│   ├── service/
│   │   ├── code_review.go   # 代码审查服务（支持评分）
│   │   ├── gitlab.go        # GitLab API 服务
│   │   ├── svn.go           # SVN 操作服务
│   │   ├── ollama.go        # Ollama API 服务
│   │   └── deepseek.go      # DeepSeek API 服务
│   ├── dto/
│   │   ├── gitlab.go        # GitLab 数据结构
│   │   ├── svn.go           # SVN 数据结构（含 pre-commit 响应）
│   │   ├── ollama.go        # Ollama 数据结构
│   │   └── deepseek.go      # DeepSeek 数据结构
│   └── storage/
│       └── file.go          # 文件存储（含列表功能）
├── pre-commit-hook/
│   ├── README.md            # Pre-commit 钩子文档
│   ├── pre-commit           # Linux/Mac 钩子脚本
│   ├── pre-commit.ps1       # Windows PowerShell 钩子
│   └── pre-commit.bat       # Windows Batch 钩子
├── web/
│   └── dist/                # 前端构建产物
├── docker/
│   ├── Dockerfile           # Docker 镜像构建
│   └── docker-compose.yml   # 完整部署配置
├── go.mod
└── go.sum
```

## 配置项

所有配置通过环境变量设置：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| SERVER_BASE_URL | http://localhost:8080 | 服务器基础 URL |
| SERVER_PORT | 8080 | 服务端口 |
| GITLAB_BASE_URL | http://localhost:8929 | GitLab 服务器地址 |
| GITLAB_PRIVATE_TOKEN | - | GitLab 访问令牌 |
| SVN_BASE_URL | svn://localhost:3690 | SVN 服务器地址 |
| SVN_KEYWORD | [review] | SVN 审查触发关键词 |
| SVN_PATH | ./reviews | SVN 审查结果存储路径 |
| **LLM_BACKEND** | ollama | LLM 后端选择 ("ollama" 或 "deepseek") |
| OLLAMA_BASE_URL | http://localhost:11434 | Ollama 服务器地址 |
| OLLAMA_MODEL | llama3 | 使用的 LLM 模型 |
| **DEEPSEEK_BASE_URL** | https://api.deepseek.com | DeepSeek API 地址 |
| **DEEPSEEK_API_KEY** | - | DeepSeek API 密钥 (需要申请) |
| **DEEPSEEK_MODEL** | deepseek-v4-flash | DeepSeek 模型 |
| REVIEW_PROMPT | (内置) | 自定义审查提示词 |

## 使用 DeepSeek API

要使用 DeepSeek API 进行代码审查，需要进行以下配置：

```bash
# 设置环境变量
export LLM_BACKEND=deepseek
export DEEPSEEK_API_KEY=your_deepseek_api_key_here
export DEEPSEEK_MODEL=deepseek-v4-flash  # 或 deepseek-v4-pro

# 运行程序
./review-bot.exe
```

DeepSeek API 使用 OpenAI 兼容的接口格式，支持以下模型：
- `deepseek-v4-flash` (默认，快速响应)
- `deepseek-v4-pro` (更强的推理能力)

## 快速开始

### 本地运行

```bash
# 克隆项目
cd git-svn-reviewbot-go

# 安装依赖
go mod tidy

# 构建
go build -o review-bot.exe main.go

# 运行 (默认使用 Ollama)
./review-bot.exe

# 或使用 DeepSeek
set LLM_BACKEND=deepseek
set DEEPSEEK_API_KEY=your_api_key
./review-bot.exe
```

### Docker 运行

```bash
cd docker
docker-compose up -d
```

## Web UI 界面

启动服务器后，访问以下地址：

- **首页**: `http://localhost:8080/`
  - 左侧：审查历史列表（按时间倒序）
  - 右侧：选中审查报告详情

- **直接查看特定审查**: `http://localhost:8080/?version=123`
- **带仓库参数**: `http://localhost:8080/?repo=xxx&version=123`

## API 端点

### 通用
- `GET /health` - 健康检查

### GitLab
- `POST /webhook/gitlab` - GitLab Webhook (自动触发审查)

### SVN
- `POST /webhook/svn` - SVN Webhook (自动触发审查)
- `GET /api/reviews/svn?repo=xxx&revision=xxx` - 获取 SVN 审查结果 (Markdown)
- `GET /api/reviews/svn/trigger?repo=xxx&revision=xxx` - **手动触发 SVN 审查**
- `POST /api/reviews/svn/trigger?repo=xxx&revision=xxx` - **手动触发 SVN 审查**
- `GET /api/reviews/list` - **获取所有审查历史列表**
- `GET /api/review/detail?repo=xxx&version=xxx` - **获取审查报告详情 (JSON)**
- `POST /api/reviews/svn/pre-commit` - **Pre-commit 审查 API（带评分）**

### 手动触发 SVN 审查接口说明

#### 接口地址
```
GET|POST /api/reviews/svn/trigger?repo=xxx&revision=xxx
```

#### 参数说明
- **查询参数**：
  - `repo`: SVN 仓库路径，支持多层目录，如 `rod_s/dinosaur_island/dev_trunk`（**可选**，见下面配置说明）
  - `revision`: 版本号（必填）
  - `author`: 作者名（可选），默认 "manual"
  - `message`: 提交信息（可选），默认 "[review] Manual review trigger"

#### 配置说明

支持两种配置方式：

##### 方式 1: 完整路径配置（推荐，适合单仓库）
直接把完整路径放在 `SVN_BASE_URL` 里：
```bash
SVN_BASE_URL=svn://10.55.19.205/rod_s/dinosaur_island/dev_trunk
```
调用时不需要 `repo` 参数：
```bash
curl http://localhost:8080/api/reviews/svn/trigger?revision=12345
```

##### 方式 2: 服务器路径配置（适合多仓库）
把服务器地址放在 `SVN_BASE_URL`，仓库路径通过 `repo` 参数传递：
```bash
SVN_BASE_URL=svn://10.55.19.205
```
调用时加上 `repo` 参数：
```bash
curl http://localhost:8080/api/reviews/svn/trigger?repo=rod_s/dinosaur_island/dev_trunk&revision=12345
```

#### 使用示例

**示例 1: 完整路径配置（推荐）**
```bash
# 配置
SVN_BASE_URL=svn://10.55.19.205/rod_s/dinosaur_island/dev_trunk

# 触发审查
curl http://localhost:8080/api/reviews/svn/trigger?revision=13262

# 查询结果
http://localhost:8080/?revision=13262
```

**示例 2: 服务器路径配置**
```bash
# 配置
SVN_BASE_URL=svn://10.55.19.205

# 触发审查
curl http://localhost:8080/api/reviews/svn/trigger?repo=rod_s/dinosaur_island/dev_trunk&revision=13262

# 查询结果
http://localhost:8080/?repo=rod_s/dinosaur_island/dev_trunk&revision=13262
```

**示例 3: 带自定义参数**
```bash
curl "http://localhost:8080/api/reviews/svn/trigger?revision=12345&author=zhang_san&message=[review]+修复登录模块"
```

**示例 4: 使用 POST**
```bash
curl -X POST "http://localhost:8080/api/reviews/svn/trigger?revision=12345"
```

**响应示例**：
```json
{
  "reviewId": 12345,
  "reviewUrl": "http://localhost:8080/?revision=12345",
  "message": "Manual review triggered successfully."
}
```

**查询审查结果**（在触发后访问）：
```
http://localhost:8080/?revision=12345
```

#### 特性
- ✅ 支持 GET 和 POST 两种方法
- ✅ 支持完整多层路径，如 `rod_s/dinosaur_island/dev_trunk`
- ✅ 自动添加 [review] 关键字（如果消息中没有）
- ✅ 异步处理，立即返回结果 URL
- ✅ 可自定义作者和提交信息
- ✅ 文件命名安全处理，路径中的 `/` 自动转为 `_`

## SVN Pre-commit 钩子

### 功能说明

Pre-commit 钩子在代码提交前进行 AI 审查，根据评分决定是否阻止提交：
- **评分 > 6 分**: 允许提交
- **评分 ≤ 6 分**: 阻止提交

### 配置方式

在你的 SVN 仓库 `hooks` 目录中安装钩子：

#### Linux/MacOS
```bash
cd /path/to/svn/repo/hooks
cp /path/to/git-svn-reviewbot-go/pre-commit-hook/pre-commit ./pre-commit
chmod +x pre-commit
# 编辑 pre-commit 中的配置
```

#### Windows
```batch
cd C:\path\to\svn\repo\hooks
copy C:\path\to\git-svn-reviewbot-go\pre-commit-hook\pre-commit.bat pre-commit.bat
copy C:\path\to\git-svn-reviewbot-go\pre-commit-hook\pre-commit.ps1 pre-commit.ps1
# 编辑 pre-commit.ps1 中的配置
```

### 钩子配置

修改脚本中的以下参数：
- `REVIEW_API_URL`: 审查服务地址
- `REPO_NAME`: 仓库名称
- `TIMEOUT`: API 超时（默认 30 秒，超时时允许提交）

### Pre-commit API 说明

#### 接口地址
```
POST /api/reviews/svn/pre-commit
```

#### 请求体
```json
{
  "repo": "my-repo",
  "diff": "svn diff 内容",
  "author": "username",
  "message": "commit message"
}
```

#### 响应
```json
{
  "allowed": true,
  "score": 8,
  "report": "完整 Markdown 报告...",
  "url": "http://localhost:8080/?repo=xxx&version=xxx",
  "message": "代码审查通过！评分：8/10 分"
}
```

### 评分标准

| 评分 | 说明 |
|-----|------|
| 10 | 完美，没有任何问题 |
| 8-9 | 优秀，只有很小的改进空间 |
| 6-7 | 合格，有一些问题但可以接受 |
| 4-5 | 较差，有明显问题需要修改 |
| 0-3 | 极差，必须重写 |

### 跳过审查

在提交信息中添加 `[skip-review]` 标记可以跳过审查：
```bash
svn commit -m "[skip-review] 紧急修复！"
```

详细说明见 `pre-commit-hook/README.md`。

## 注意事项

1. **SVN 客户端**: 如果使用命令行方案，需要确保环境中已安装 SVN 客户端
2. **Ollama 模型**: 需要在 Ollama 服务器上下载并运行指定的模型
3. **DeepSeek API**: 使用 DeepSeek 时需要申请 API 密钥并配置环境变量
4. **GitLab Webhook**: 需要在 GitLab 项目设置中配置 Webhook 指向 `/webhook/gitlab`
5. **Pre-commit 钩子**: 确保脚本有执行权限且能访问审查服务
