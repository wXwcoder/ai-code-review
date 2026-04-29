import React, { useState, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import axios from 'axios'

function App() {
  const [reviews, setReviews] = useState([])
  const [selectedReview, setSelectedReview] = useState(null)
  const [loading, setLoading] = useState(false)

  const fetchReviews = async () => {
    setLoading(true)
    try {
      const response = await axios.get('/api/reviews/list')
      setReviews(response.data.reviews || [])
    } catch (error) {
      console.error('获取审查列表失败:', error)
    } finally {
      setLoading(false)
    }
  }

  const fetchSingleReview = async (repo, version) => {
    try {
      const response = await axios.get('/api/review/detail', {
        params: { repo, version }
      })
      const data = response.data
      setSelectedReview({
        repoPath: data.repoPath,
        revision: data.revision,
        content: data.report,
        timestamp: data.timestamp,
        fileName: data.fileName
      })
    } catch (error) {
      console.error('获取单个审查报告失败:', error)
    }
  }

  useEffect(() => {
    // 检查 URL 参数
    const params = new URLSearchParams(window.location.search)
    const version = params.get('version')
    const repo = params.get('repo') || null
    
    if (version) {
      fetchSingleReview(repo, version)
    }
    // 总是获取列表
    fetchReviews()
  }, [])

  const formatDate = (dateString) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  }

  const getDisplayRepo = (repoPath) => {
    if (!repoPath || repoPath === 'default') {
      return '代码审查'
    }
    return repoPath
  }

  const handleSelectReview = (review) => {
    setSelectedReview(review)
    // 更新 URL 参数，方便分享
    const params = new URLSearchParams()
    if (review.repoPath && review.repoPath !== 'default') {
      params.set('repo', review.repoPath)
    }
    params.set('version', review.revision.toString())
    const url = `${window.location.pathname}?${params.toString()}`
    window.history.replaceState(null, '', url)
  }

  return (
    <div className="app">
      {/* ===== Sidebar ===== */}
      <div className="sidebar">
        <div className="sidebar-header">
          <h1>代码审查</h1>
          <p>AI 代码审查历史</p>
          <button 
            className="refresh-button" 
            onClick={fetchReviews}
            disabled={loading}
          >
            <span className="refresh-icon">{loading ? '⏳' : '🔄'}</span>
            {loading ? '刷新中...' : '刷新列表'}
          </button>
        </div>
        <div className="sidebar-content">
          {reviews.length === 0 ? (
            <div className="empty-state">
              <div className="empty-icon">📝</div>
              <p>暂无审查记录</p>
            </div>
          ) : (
            <div className="review-list">
              {reviews.map((review, index) => (
                <div
                  key={index}
                  className={`review-item ${selectedReview && selectedReview.revision === review.revision ? 'active' : ''}`}
                  onClick={() => handleSelectReview(review)}
                >
                  <div className="review-item-header">
                    <span className="review-repo">{getDisplayRepo(review.repoPath)}</span>
                    <span className="review-revision">r{review.revision}</span>
                  </div>
                  <div className="review-meta">
                    <span className="review-date">
                      <span>📅</span>
                      {formatDate(review.timestamp)}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* ===== Main Content ===== */}
      <div className="main-content">
        {selectedReview ? (
          <>
            <div className="main-header">
              <h2>审查详情</h2>
              <div className="main-meta">
                <span className="main-revision">r{selectedReview.revision}</span>
                <span>{formatDate(selectedReview.timestamp)}</span>
              </div>
            </div>
            <div className="content-container">
              <div className="markdown-content">
                <ReactMarkdown remarkPlugins={[remarkGfm]}>
                  {selectedReview.content}
                </ReactMarkdown>
              </div>
            </div>
          </>
        ) : (
          <div className="no-selection">
            <div className="no-selection-icon">👈</div>
            <p>请从左侧选择一条审查记录</p>
          </div>
        )}
      </div>
    </div>
  )
}

export default App
