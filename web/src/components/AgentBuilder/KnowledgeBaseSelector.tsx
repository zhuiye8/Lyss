import React, { useEffect, useState } from 'react';
import { Select, Spin, Tag, Typography, Empty } from 'antd';
import { DatabaseOutlined } from '@ant-design/icons';
import { getKnowledgeBases } from '../../services/agentService';
import { IKnowledgeBase } from '../../types/agent';

interface KnowledgeBaseSelectorProps {
  value?: string[];
  onChange?: (value: string[]) => void;
  disabled?: boolean;
  mode?: 'multiple' | 'single';
}

const KnowledgeBaseSelector: React.FC<KnowledgeBaseSelectorProps> = ({
  value,
  onChange,
  disabled = false,
  mode = 'multiple',
}) => {
  const [knowledgeBases, setKnowledgeBases] = useState<IKnowledgeBase[]>([]);
  const [loading, setLoading] = useState(false);

  // 获取知识库列表
  const fetchKnowledgeBases = async () => {
    try {
      setLoading(true);
      const data = await getKnowledgeBases();
      setKnowledgeBases(data);
    } catch (error) {
      console.error('获取知识库列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchKnowledgeBases();
  }, []);

  // 渲染知识库选项
  const renderKnowledgeBaseOption = (kb: IKnowledgeBase) => {
    return (
      <Select.Option key={kb.id} value={kb.id}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <DatabaseOutlined style={{ marginRight: 8, color: '#1890ff' }} />
          <div>
            <Typography.Text strong>{kb.name}</Typography.Text>
            <div>
              <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                {kb.description || '没有描述'}
              </Typography.Text>
            </div>
            <div style={{ marginTop: 4 }}>
              <Tag color="blue">{kb.documentCount} 文档</Tag>
              <Tag color="green">{kb.vectorCount} 向量</Tag>
            </div>
          </div>
        </div>
      </Select.Option>
    );
  };

  return (
    <Select
      placeholder="选择知识库"
      value={value}
      onChange={onChange}
      style={{ width: '100%' }}
      loading={loading}
      disabled={disabled}
      mode={mode === 'multiple' ? 'multiple' : undefined}
      allowClear
      showSearch
      optionFilterProp="children"
      notFoundContent={
        loading ? (
          <Spin size="small" />
        ) : (
          <Empty description="没有可用的知识库" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        )
      }
    >
      {knowledgeBases.map(renderKnowledgeBaseOption)}
    </Select>
  );
};

export default KnowledgeBaseSelector; 