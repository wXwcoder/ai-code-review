package dto

// GitLabWebhookMergeRequestEvent GitLab Webhook 合并请求事件
type GitLabWebhookMergeRequestEvent struct {
	ObjectKind        string                    `json:"object_kind"`
	User              GitLabUser                `json:"user"`
	Project           GitLabProject             `json:"project"`
	ObjectAttributes  GitLabMergeRequestAttrs   `json:"object_attributes"`
	Changes           *GitLabChanges            `json:"changes,omitempty"`
}

// GitLabUser GitLab 用户信息
type GitLabUser struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
}

// GitLabProject GitLab 项目信息
type GitLabProject struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
}

// GitLabMergeRequestAttrs GitLab 合并请求属性
type GitLabMergeRequestAttrs struct {
	ID           int64         `json:"id"`
	IID          int64         `json:"iid"`
	Title        string        `json:"title"`
	State        string        `json:"state"`
	SourceBranch string        `json:"source_branch"`
	TargetBranch string        `json:"target_branch"`
	URL          string        `json:"url"`
	LastCommit   *GitLabCommit `json:"last_commit,omitempty"`
	Action       string        `json:"action"`
}

// GitLabCommit GitLab 提交信息
type GitLabCommit struct {
	ID        string     `json:"id"`
	Message   string     `json:"message"`
	Timestamp string     `json:"timestamp"`
	URL       string     `json:"url"`
	Author    GitLabUser `json:"author"`
}

// GitLabChanges GitLab 变更信息
type GitLabChanges struct {
	UpdatedByID *GitLabChangeDetail `json:"updated_by_id,omitempty"`
}

// GitLabChangeDetail GitLab 变更详情
type GitLabChangeDetail struct {
	Previous *int64 `json:"previous"`
	Current  *int64 `json:"current"`
}

// GitLabDiffResponse GitLab 差异响应
type GitLabDiffResponse struct {
	Diff        string `json:"diff"`
	NewPath     string `json:"new_path"`
	OldPath     string `json:"old_path"`
	DeletedFile bool   `json:"deleted_file"`
	NewFile     bool   `json:"new_file"`
	RenamedFile bool   `json:"renamed_file"`
}

// GitLabCommentRequest GitLab 评论请求
type GitLabCommentRequest struct {
	Body string `json:"body"`
}

// GitLabCommentResponse GitLab 评论响应
type GitLabCommentResponse struct {
	ID        int64      `json:"id"`
	Body      string     `json:"body"`
	Author    GitLabUser `json:"author"`
	CreatedAt string     `json:"created_at"`
}
