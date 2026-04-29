package dto

// SvnWebhookRequest SVN Webhook 请求
type SvnWebhookRequest struct {
	Repo     string `json:"repo"`
	Revision int64  `json:"revision"`
	Author   string `json:"author"`
	Message  string `json:"message"`
}

// SvnWebhookResponse SVN Webhook 响应
type SvnWebhookResponse struct {
	ReviewID  int64  `json:"reviewId"`
	ReviewURL string `json:"reviewUrl"`
	Message   string `json:"message"`
}

// PreCommitRequest pre-commit 审查请求
type PreCommitRequest struct {
	Repo    string `json:"repo"`
	Diff    string `json:"diff"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

// PreCommitResponse pre-commit 审查响应
type PreCommitResponse struct {
	Allowed bool   `json:"allowed"` // 是否允许提交
	Score   int    `json:"score"`   // 评分
	Report  string `json:"report"`  // 审查报告内容
	URL     string `json:"url"`     // 审查报告链接（如果已保存）
	Message string `json:"message"` // 额外消息
}

// ReviewDetailResponse 审查报告详情响应
type ReviewDetailResponse struct {
	RepoPath  string `json:"repoPath"`
	Revision  int64  `json:"revision"`
	Report    string `json:"report"`
	Timestamp string `json:"timestamp"`
	FileName  string `json:"fileName"`
}
