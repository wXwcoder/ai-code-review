import React, { useState } from 'react'

function DataInputModal({ isOpen, onClose, onSubmit }) {
  const [inputData, setInputData] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()
    if (inputData.trim()) {
      onSubmit(inputData)
      setInputData('')
    }
  }

  const handleClose = () => {
    setInputData('')
    onClose()
  }

  if (!isOpen) return null

  return (
    <div className="modal-overlay">
      <div className="modal">
        <div className="modal-header">
          <h3>输入JSON数据</h3>
          <button className="close-btn" onClick={handleClose}>×</button>
        </div>
        <form onSubmit={handleSubmit}>
          <div className="modal-body">
            <textarea
              value={inputData}
              onChange={(e) => setInputData(e.target.value)}
              placeholder="请输入JSON数据..."
              rows={10}
              className="json-input"
            />
          </div>
          <div className="modal-footer">
            <button type="button" className="btn btn-secondary" onClick={handleClose}>
              取消
            </button>
            <button type="submit" className="btn btn-primary">
              确认
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

export default DataInputModal