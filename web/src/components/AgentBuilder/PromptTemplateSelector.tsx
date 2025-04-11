import React, { useEffect, useState } from 'react';
import { Select, Spin, Tag, Typography, Empty, Button, Modal, Form, Input, Space, message } from 'antd';
import { FileTextOutlined, PlusOutlined } from '@ant-design/icons';
import { getPromptTemplates, createPromptTemplate } from '../../services/agentService';
import { IPromptTemplate } from '../../types/agent';

const { TextArea } = Input;

interface PromptTemplateSelectorProps {
  value?: string;
  onChange?: (value: string, template?: IPromptTemplate) => void;
  disabled?: boolean;
  onSelect?: (template: IPromptTemplate) => void;
}

const PromptTemplateSelector: React.FC<PromptTemplateSelectorProps> = ({
  value,
  onChange,
  disabled = false,
  onSelect,
}) => {
  const [templates, setTemplates] = useState<IPromptTemplate[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [submitting, setSubmitting] = useState(false);

  // 获取提示词模板列表
  const fetchTemplates = async () => {
    try {
      setLoading(true);
      const data = await getPromptTemplates();
      setTemplates(data);
    } catch (error) {
      console.error('获取提示词模板列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTemplates();
  }, []);

  // 创建新模板
  const handleCreateTemplate = async (values: any) => {
    try {
      setSubmitting(true);
      // 解析变量
      const variableMatches = values.content.match(/\{\{([^}]+)\}\}/g) || [];
      const variables = variableMatches.map((match: string) => 
        match.replace(/\{\{|\}\}/g, '').trim()
      );
      
      const newTemplate = await createPromptTemplate({
        ...values,
        variables: [...new Set(variables)], // 去重
      });
      
      setTemplates([...templates, newTemplate]);
      message.success('模板创建成功');
      setModalVisible(false);
      form.resetFields();
      
      // 如果有onChange，自动选择新创建的模板
      if (onChange) {
        onChange(newTemplate.id, newTemplate);
      }
      
    } catch (error) {
      console.error('创建提示词模板失败:', error);
      message.error('创建提示词模板失败');
    } finally {
      setSubmitting(false);
    }
  };

  // 选择模板
  const handleSelectTemplate = (templateId: string) => {
    if (onChange) {
      const template = templates.find(t => t.id === templateId);
      onChange(templateId, template);
      
      if (onSelect && template) {
        onSelect(template);
      }
    }
  };

  // 渲染模板选项
  const renderTemplateOption = (template: IPromptTemplate) => {
    return (
      <Select.Option key={template.id} value={template.id}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <FileTextOutlined style={{ marginRight: 8, color: '#1890ff' }} />
            <Typography.Text strong>{template.name || '未命名模板'}</Typography.Text>
          </div>
          <div style={{ marginTop: 4, marginLeft: 22 }}>
            <Typography.Text type="secondary" style={{ fontSize: 12 }}>
              {template.description || '没有描述'}
            </Typography.Text>
          </div>
          <div style={{ marginTop: 4, marginLeft: 22 }}>
            {template.variables.map((variable, index) => (
              <Tag key={index} color="blue">{`{{${variable}}}`}</Tag>
            ))}
          </div>
        </div>
      </Select.Option>
    );
  };

  // 自定义下拉菜单页脚
  const dropdownRender = (menu: React.ReactNode) => (
    <div>
      {menu}
      <div style={{ padding: '8px', textAlign: 'center' }}>
        <Button 
          type="dashed" 
          icon={<PlusOutlined />} 
          onClick={() => setModalVisible(true)}
          style={{ width: '100%' }}
        >
          创建新模板
        </Button>
      </div>
    </div>
  );

  return (
    <>
      <Select
        placeholder="选择提示词模板"
        value={value}
        onChange={handleSelectTemplate}
        style={{ width: '100%' }}
        loading={loading}
        disabled={disabled}
        allowClear
        showSearch
        optionFilterProp="children"
        dropdownRender={dropdownRender}
        notFoundContent={
          loading ? (
            <Spin size="small" />
          ) : (
            <Empty description="没有可用的提示词模板" image={Empty.PRESENTED_IMAGE_SIMPLE} />
          )
        }
      >
        {templates.map(renderTemplateOption)}
      </Select>

      {/* 创建模板对话框 */}
      <Modal
        title="创建提示词模板"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateTemplate}
        >
          <Form.Item
            name="name"
            label="模板名称"
            rules={[{ required: true, message: '请输入模板名称' }]}
          >
            <Input placeholder="给模板起个名字" />
          </Form.Item>
          
          <Form.Item
            name="description"
            label="描述"
          >
            <Input placeholder="模板用途描述（可选）" />
          </Form.Item>
          
          <Form.Item
            name="content"
            label="模板内容"
            rules={[{ required: true, message: '请输入模板内容' }]}
            extra="使用 {{变量名}} 定义变量，例如：{{用户名}}、{{日期}}"
          >
            <TextArea
              placeholder="请输入提示词模板内容..."
              autoSize={{ minRows: 6, maxRows: 12 }}
            />
          </Form.Item>
          
          <Form.Item>
            <Space style={{ display: 'flex', justifyContent: 'flex-end' }}>
              <Button onClick={() => setModalVisible(false)}>取消</Button>
              <Button type="primary" htmlType="submit" loading={submitting}>
                创建
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
};

export default PromptTemplateSelector; 