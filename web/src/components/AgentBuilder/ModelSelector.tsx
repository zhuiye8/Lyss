import React, { useEffect, useState } from 'react';
import { Select, Spin, Typography, Badge, Tooltip } from 'antd';
import { InfoCircleOutlined } from '@ant-design/icons';
import { getModels } from '../../services/dashboardService';
import { IModelData } from '../../types/dashboard';

interface ModelSelectorProps {
  value?: string;
  onChange?: (value: string) => void;
  disabled?: boolean;
}

const ModelSelector: React.FC<ModelSelectorProps> = ({
  value,
  onChange,
  disabled = false,
}) => {
  const [models, setModels] = useState<IModelData[]>([]);
  const [loading, setLoading] = useState(false);

  // 获取模型列表
  const fetchModels = async () => {
    try {
      setLoading(true);
      const data = await getModels();
      const activeModels = data.filter(model => model.status === 'active');
      setModels(activeModels);
    } catch (error) {
      console.error('获取模型列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchModels();
  }, []);

  // 渲染模型选项
  const renderModelOption = (model: IModelData) => {
    return (
      <Select.Option key={model.id} value={model.id}>
        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <div>
            <Typography.Text strong>{model.name}</Typography.Text>
            <Typography.Text type="secondary" style={{ marginLeft: 8 }}>
              ({model.provider})
            </Typography.Text>
          </div>
          <div>
            {model.supportsFunctionCalling && (
              <Tooltip title="支持函数调用">
                <Badge status="success" text="函数调用" />
              </Tooltip>
            )}
          </div>
        </div>
      </Select.Option>
    );
  };

  return (
    <Select
      placeholder="选择模型"
      value={value}
      onChange={onChange}
      style={{ width: '100%' }}
      loading={loading}
      disabled={disabled}
      notFoundContent={loading ? <Spin size="small" /> : '无可用模型'}
      optionLabelProp="label"
    >
      {models.map(renderModelOption)}
    </Select>
  );
};

export default ModelSelector; 