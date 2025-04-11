import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { StopOutlined } from '@ant-design/icons';

interface EndNodeProps {
  data: {
    label: string;
  };
}

const EndNode: React.FC<EndNodeProps> = memo(({ data }) => {
  return (
    <div style={{
      padding: '10px',
      borderRadius: '5px',
      background: '#fff1f0',
      border: '1px solid #ffa39e',
      width: '150px',
      textAlign: 'center'
    }}>
      <StopOutlined style={{ fontSize: '24px', color: '#f5222d' }} />
      <div style={{ marginTop: '5px', fontSize: '14px', fontWeight: 'bold' }}>
        {data.label || '结束'}
      </div>
      <Handle
        type="target"
        position={Position.Top}
        style={{ background: '#f5222d' }}
      />
    </div>
  );
});

export default EndNode; 