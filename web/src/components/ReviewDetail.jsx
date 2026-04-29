import React, { useState, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import axios from 'axios'

const ReviewDetail = ({ onBack }) => {
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [review, setReview] = useState(null)

  useEffect(() => {
    // 从 URL 参数获取
    const params = new URLSearchParams(window.location.search)
    const version = params.get('version')
    const repo = params.get('repo') || 'default'

    if (!version) {
      setError('缺少版本号参数')
      setLoading(false)
      return
    }

    fetchReview(repo, version)
  }, [])

  const fetchReview = async (repo, version) => {
    try {
      setLoading(true)
      const response = await axios.get('/api/review/detail', {
        params: { repo, version }
      })
      setReview(response.data)
    } catch (err) {
      console.error('获取审查报告失败:', err)
      if (err.response && err.response.status === 404) {
        setError('未找到审查报告')
      } else {
        setError('获取审查报告失败: ' + (err.response?.data?.error || err.message))
      }
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="review-detail-page">
        <div className="loading-container">
          <div className="spinner"></div>
          <p>加载中...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="review-detail-page">
        <div className="error-container">
          <div className="error-icon">⚠️</div>
          <h2>出错了</h2>
          <p>{error}</p>
          {onBack && (
            <button className="back-button" onClick={onBack}>
              返回列表
            </button>
          )}
        </div>
      </div>
    )
  }

  return (
    <div className="review-detail-page">
      <div className="review-detail-header">
        {onBack && (
          <button className="back-button" onClick={onBack}>
            ← 返回列表
          </button>
        )}
        <div className="review-info">
          <h1>代码审查报告</h1>
          {review && (
            <div className="review-meta">
              {review.repoPath && <span className="meta-item">仓库: {review.repoPath}</span>}
              <span className="meta-item">版本: {review.revision}</span>
              <span className="meta-item">时间: {review.timestamp}</span>
            </div>
          )}
        </div>
      </div>
      <div className="review-content">
        <div className="markdown-wrapper">
          {review && <ReactMarkdown remarkPlugins={[remarkGfm]}>{review.report}</ReactMarkdown>}
        </div>
      </div>
    </div>
  )
}

export default ReviewDetail
