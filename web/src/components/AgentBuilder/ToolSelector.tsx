import React, { useEffect, useState } from 'react';
import { Select, Spin, Tag, Typography, Empty, Tooltip } from 'antd';
import { ApiOutlined, CodeOutlined, AppstoreOutlined } from '@ant-design/icons';
import { getTools } from '../../services/agentService';
import { ITool } from '../../types/agent';

interface ToolSelectorProps {
  value?: string[];
  onChange?: (value: string[]) => void;
  disabled?: boolean;
  mode?: 'multiple' | 'single';
}

const ToolSelector: React.FC<ToolSelectorProps> = ({
  value,
  onChange,
  disabled = false,
  mode = 'multiple',
}) => {
  const [tools, setTools] = useState<ITool[]>([]);
  const [loading, setLoading] = useState(false);

  // 获取工具列表
  const fetchTools = async () => {
    try {
      setLoading(true);
      const data = await getTools();
      setTools(data);
    } catch (error) {
      console.error('获取工具列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTools();
  }, []);

  // 获取工具图标
  const getToolIcon = (type: string) => {
    switch (type) {
      case 'api':
        return <ApiOutlined style={{ color: '#1890ff' }} />;
      case 'function':
        return <CodeOutlined style={{ color: '#52c41a' }} />;
      case 'builtin':
        return <AppstoreOutlined style={{ color: '#faad14' }} />;
      default:
        return <ApiOutlined />;
    }
  };

  // 获取工具类型标签
  const getToolTypeTag = (type: string) => {
    switch (type) {
      case 'api':
        return <Tag color="blue">API</Tag>;
      case 'function':
        return <Tag color="green">函数</Tag>;
      case 'builtin':
        return <Tag color="orange">内置</Tag>;
      default:
        return <Tag>工具</Tag>;
    }
  };

  // 渲染工具选项
  const renderToolOption = (tool: ITool) => {
    return (
      <Select.Option key={tool.id} value={tool.id}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          {getToolIcon(tool.type)}
          <div style={{ marginLeft: 8 }}>
            <Typography.Text strong>{tool.name}</Typography.Text>
            <div>
              <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                {tool.description || '没有描述'}
              </Typography.Text>
            </div>
            <div style={{ marginTop: 4 }}>
              {getToolTypeTag(tool.type)}
              <Tag color="purple">{tool.parameters.length} 参数</Tag>
            </div>
          </div>
        </div>
      </Select.Option>
    );
  };

  return (
    <Select
      placeholder="选择工具"
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
          <Empty description="没有可用的工具" image={Empty.PRESENTED_IMAGE_SIMPLE} />
        )
      }
    >
      {tools.map(renderToolOption)}
    </Select>
  );
};

export default ToolSelector; 