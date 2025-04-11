import React, { memo } from 'react';
import { Handle, Position } from 'reactflow';
import { QuestionCircleOutlined } from '@ant-design/icons';

interface ConditionNodeProps {
  data: {
    condition: string;
  };
}

const ConditionNode: React.FC<ConditionNodeProps> = memo(({ data }) => {
  return (
    <div style={{
      padding: '10px',
      borderRadius: '5px',
      background: '#fff7e6',
      border: '1px solid #ffd591',
      width: '200px'
    }}>
      <div style={{ display: 'flex', alignItems: 'center' }}>
        <QuestionCircleOutlined style={{ fontSize: '18px', color: '#fa8c16', marginRight: '8px' }} />
        <div style={{ fontSize: '14px', fontWeight: 'bold' }}>条件节点</div>
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
        textOverflow: 'ellipsis',
        fontFamily: 'monospace'
      }}>
        {data.condition || '条件未设置'}
      </div>
      <Handle
        type="target"
        position={Position.Top}
        style={{ background: '#fa8c16' }}
      />
      <Handle
        type="source"
        position={Position.Bottom}
        id="true"
        style={{ background: '#52c41a', left: '30%' }}
      />
      <Handle
        type="source"
        position={Position.Bottom}
        id="false"
        style={{ background: '#f5222d', left: '70%' }}
      />
    </div>
  );
});

export default ConditionNode; 