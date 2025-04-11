import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { DatabaseOutlined } from '@ant-design/icons';

interface KnowledgeNodeProps {
  data: {
    knowledgeBaseId: string;
  };
}

const KnowledgeNode: React.FC<KnowledgeNodeProps> = memo(({ data }) => {
  return (
    <div style={{
      padding: '10px',
      borderRadius: '5px',
      background: '#e6f7ff',
      border: '1px solid #91d5ff',
      width: '200px'
    }}>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <DatabaseOutlined style={{ fontSize: '18px', color: '#1890ff', marginRight: '8px' }} />
        <div style={{ fontSize: '14px', fontWeight: 'bold' }}>知识库节点</div>
      </div>
      <div style={{ 
        marginTop: '8px', 
        fontSize: '12px',
        padding: '6px',
        background: '#ffffff',
        borderRadius: '4px',
        border: '1px dashed #d9d9d9'
      }}>
        <div>知识库ID: {data.knowledgeBaseId || '未选择'}</div>
      </div>
      <Handle
        type="target"
        position={Position.Top}
        style={{ background: '#1890ff' }}
      />
      <Handle
        type="source"
        position={Position.Bottom}
        style={{ background: '#1890ff' }}
      />
    </div>
  );
});

export default KnowledgeNode; 