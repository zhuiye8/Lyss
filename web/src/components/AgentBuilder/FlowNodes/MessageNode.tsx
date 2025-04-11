import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { MessageOutlined } from '@ant-design/icons';

interface MessageNodeProps {
  data: {
    message: string;
  };
}

const MessageNode: React.FC<MessageNodeProps> = memo(({ data }) => {
  return (
    <div style={{
      padding: '10px',
      borderRadius: '5px',
      background: '#f6ffed',
      border: '1px solid #b7eb8f',
      width: '200px'
    }}>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <MessageOutlined style={{ fontSize: '18px', color: '#52c41a', marginRight: '8px' }} />
        <div style={{ fontSize: '14px', fontWeight: 'bold' }}>消息节点</div>
      </div>
      <div style={{ 
        marginTop: '8px', 
        fontSize: '12px',
        padding: '6px',
        background: '#ffffff',
        borderRadius: '4px',
        border: '1px dashed #d9d9d9',
        maxHeight: '80px',
        overflow: 'hidden',
        textOverflow: 'ellipsis'
      }}>
        {data.message || '无消息内容'}
      </div>
      <Handle
        type="target"
        position={Position.Top}
        style={{ background: '#52c41a' }}
      />
      <Handle
        type="source"
        position={Position.Bottom}
        style={{ background: '#52c41a' }}
      />
    </div>
  );
});

export default MessageNode; 