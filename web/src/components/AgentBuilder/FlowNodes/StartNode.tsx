import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { PlayCircleOutlined } from '@ant-design/icons';

interface StartNodeProps {
  data: {
    label: string;
  };
}

const StartNode: React.FC<StartNodeProps> = memo(({ data }) => {
  return (
    <div style={{
      padding: '10px',
      borderRadius: '5px',
      background: '#f0f5ff',
      border: '1px solid #adc6ff',
      width: '150px',
      textAlign: 'center'
    }}>
      <PlayCircleOutlined style={{ fontSize: '24px', color: '#1890ff' }} />
      <div style={{ marginTop: '5px', fontSize: '14px', fontWeight: 'bold' }}>
        {data.label || '开始'}
      </div>
      <Handle
        type="source"
        position={Position.Bottom}
        style={{ background: '#1890ff' }}
      />
    </div>
  );
});

export default StartNode; 