import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { ToolOutlined } from '@ant-design/icons';

interface ToolNodeProps {
  data: {
    toolId: string;
    toolParams?: Record<string, any>;
  };
}

const ToolNode: React.FC<ToolNodeProps> = memo(({ data }) => {
  return (
    <div style={{
      padding: '10px',
      borderRadius: '5px',
      background: '#f9f0ff',
      border: '1px solid #d3adf7',
      width: '200px'
    }}>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <ToolOutlined style={{ fontSize: '18px', color: '#722ed1', marginRight: '8px' }} />
        <div style={{ fontSize: '14px', fontWeight: 'bold' }}>工具节点</div>
      </div>
      <div style={{ 
        marginTop: '8px', 
        fontSize: '12px',
        padding: '6px',
        background: '#ffffff',
        borderRadius: '4px',
        border: '1px dashed #d9d9d9'
      }}>
        <div style={{ fontWeight: 'bold' }}>工具ID: {data.toolId || '未选择'}</div>
        {data.toolParams && Object.keys(data.toolParams).length > 0 && (
          <div style={{ 
            marginTop: '4px',
            fontSize: '11px',
            maxHeight: '60px',
            overflow: 'hidden'
          }}>
            参数: {JSON.stringify(data.toolParams).substring(0, 50)}
            {JSON.stringify(data.toolParams).length > 50 ? '...' : ''}
          </div>
        )}
      </div>
      <Handle
        type="target"
        position={Position.Top}
        style={{ background: '#722ed1' }}
      />
      <Handle
        type="source"
        position={Position.Bottom}
        style={{ background: '#722ed1' }}
      />
    </div>
  );
});

export default ToolNode; 