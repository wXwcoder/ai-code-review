import React, { useState, useCallback, useRef, useEffect } from 'react';
import { JSONTree } from 'react-json-tree';

/**
 * 主题配置
 */
const theme = {
  scheme: 'monokai',
  base00: '#f8f9fa',
  base01: '#e9ecef',
  base02: '#dee2e6',
  base03: '#adb5bd',
  base04: '#6c757d',
  base05: '#495057',
  base06: '#343a40',
  base07: '#212529',
  base08: '#dc3545',
  base09: '#fd7e14',
  base0A: '#ffc107',
  base0B: '#198754',
  base0C: '#0dcaf0',
  base0D: '#0d6efd',
  base0E: '#6610f2',
  base0F: '#d63384',
  arrow: {
    color: '#6c757d',
    fontSize: '10px',
  },
  nestedNode: {
    style: {
      padding: '0 0 0 18px',
      margin: '0',
      listStyle: 'none',
    },
  },
  nestedNodeItem: {
    style: {
      margin: '0',
      padding: '0',
    },
  },
  rootNode: {
    style: {
      padding: '0',
      margin: '0',
      listStyle: 'none',
    },
  },
};

/**
 * JSON查看器组件
 * @param {Object} data - JSON数据
 */
function JsonViewer({ data }) {
  const [key, setKey] = useState(0);
  const [expandAllState, setExpandAllState] = useState(false);

  /**
   * 展开所有
   */
  const expandAll = useCallback(() => {
    setExpandAllState(true);
    setKey(prev => prev + 1);
  }, []);

  /**
   * 收起所有
   */
  const collapseAll = useCallback(() => {
    setExpandAllState(false);
    setKey(prev => prev + 1);
  }, []);

  /**
   * 自定义标签渲染
   */
  const getLabelRenderer = () => {
    return function labelRenderer([key]) {
      return <span style={{ fontWeight: 600, color: '#0d6efd' }}>"{key}"</span>;
    };
  };

  /**
   * 自定义值渲染
   */
  const getValueRenderer = () => {
    return function valueRenderer(value) {
      if (typeof value === 'string') {
        return <span style={{ color: '#dc3545' }}>"{value}"</span>;
      }
      if (typeof value === 'number') {
        return <span style={{ color: '#198754', fontWeight: 500 }}>{value}</span>;
      }
      if (typeof value === 'boolean') {
        return <span style={{ color: '#fd7e14', fontWeight: 600 }}>{value.toString()}</span>;
      }
      if (value === null) {
        return <span style={{ color: '#6c757d', fontStyle: 'italic' }}>null</span>;
      }
      return value;
    };
  };

  /**
   * 自定义项字符串渲染
   */
  const getItemString = () => {
    return function itemString(type, data, itemType, itemString) {
      let count = '';
      if (type === 'Object') {
        count = ` ${Object.keys(data).length} keys`;
      } else if (type === 'Array') {
        count = ` ${data.length} items`;
      }
      return <span style={{ color: '#6c757d', fontSize: '11px', fontWeight: 500 }}>{count}</span>;
    };
  };

  /**
   * 检查是否应该展开
   */
  const shouldExpandNodeInitially = useCallback((keyPath, data, level) => {
    return expandAllState;
  }, [expandAllState]);

  return (
    <div className="json-viewer-wrapper">
      <div className="json-viewer-toolbar">
        <button className="toolbar-btn" onClick={expandAll} title="展开全部">
          <span className="btn-icon">📂</span> 展开全部
        </button>
        <button className="toolbar-btn" onClick={collapseAll} title="收起全部">
          <span className="btn-icon">📁</span> 收起全部
        </button>
      </div>
      <div className="json-viewer">
        <JSONTree
          key={key}
          data={data}
          theme={theme}
          invertTheme={false}
          hideRoot={false}
          getItemString={getItemString()}
          labelRenderer={getLabelRenderer()}
          valueRenderer={getValueRenderer()}
          shouldExpandNodeInitially={shouldExpandNodeInitially}
        />
      </div>
    </div>
  );
}

export default JsonViewer;
