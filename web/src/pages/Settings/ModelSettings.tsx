import React, { useEffect, useState } from 'react';
import { 
  Card, 
  Table, 
  Button, 
  Space, 
  Modal, 
  Form, 
  Input, 
  Select, 
  Switch, 
  InputNumber,
  Divider,
  Typography,
  message,
  Popconfirm,
  Badge,
  Tooltip,
  Collapse
} from 'antd';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  CheckCircleOutlined,
  CloseCircleOutlined,
  ApiOutlined,
  InfoCircleOutlined
} from '@ant-design/icons';
import { getModels, addModel, updateModel, deleteModel, testModelConnection } from '../../services/dashboardService';
import { IModelData } from '../../types/dashboard';

const { Title, Text } = Typography;
const { Option } = Select;
const { Panel } = Collapse;

const ModelSettings: React.FC = () => {
  const [models, setModels] = useState<IModelData[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalVisible, setModalVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [testingConnection, setTestingConnection] = useState(false);
  const [editingModel, setEditingModel] = useState<IModelData | null>(null);
  const [form] = Form.useForm();

  // 获取模型列表
  const fetchModels = async () => {
    try {
      setLoading(true);
      const data = await getModels();
      setModels(data);
    } catch (error) {
      console.error('获取模型列表失败:', error);
      message.error('获取模型列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    fetchModels();
  }, []);

  // 新增或编辑模型
  const handleAddOrUpdateModel = async (values: any) => {
    try {
      setConfirmLoading(true);
      if (editingModel) {
        // 更新模型
        await updateModel(editingModel.id, values);
        message.success('模型已更新');
      } else {
        // 添加新模型
        await addModel(values);
        message.success('模型已添加');
      }
      // 重新加载列表
      fetchModels();
      setModalVisible(false);
    } catch (error) {
      console.error('保存模型失败:', error);
      message.error('保存模型失败');
    } finally {
      setConfirmLoading(false);
    }
  };

  // 删除模型
  const handleDeleteModel = async (id: string) => {
    try {
      await deleteModel(id);
      message.success('模型已删除');
      // 重新加载列表
      fetchModels();
    } catch (error) {
      console.error('删除模型失败:', error);
      message.error('删除模型失败');
    }
  };

  // 测试模型连接
  const handleTestConnection = async () => {
    try {
      const values = await form.validateFields();
      setTestingConnection(true);
      const success = await testModelConnection(values);
      if (success) {
        message.success('连接测试成功');
      } else {
        message.error('连接测试失败');
      }
    } catch (error) {
      message.error('表单验证失败或测试连接时出错');
    } finally {
      setTestingConnection(false);
    }
  };

  // 处理模型编辑
  const handleEditModel = (model: IModelData) => {
    setEditingModel(model);
    form.setFieldsValue(model);
    setModalVisible(true);
  };

  // 打开新增模型对话框
  const showAddModal = () => {
    setEditingModel(null);
    form.resetFields();
    setModalVisible(true);
  };

  // 关闭对话框
  const handleCancel = () => {
    setModalVisible(false);
  };

  // 表格列配置
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '提供商',
      dataIndex: 'provider',
      key: 'provider',
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        type === 'local' ? <Tag color="green">本地</Tag> : <Tag color="blue">云端</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'active' 
          ? <Badge status="success" text="已启用" /> 
          : <Badge status="default" text="已禁用" />
      ),
    },
    {
      title: '上下文长度',
      dataIndex: 'contextLength',
      key: 'contextLength',
      render: (length: number) => `${length.toLocaleString()} tokens`,
    },
    {
      title: '函数调用',
      dataIndex: 'supportsFunctionCalling',
      key: 'supportsFunctionCalling',
      render: (supports: boolean) => (
        supports 
          ? <CheckCircleOutlined style={{ color: '#52c41a' }} /> 
          : <CloseCircleOutlined style={{ color: '#f5222d' }} />
      ),
    },
    {
      title: '上次使用',
      dataIndex: 'lastUsed',
      key: 'lastUsed',
      render: (date: string) => date ? new Date(date).toLocaleString() : '未使用',
    },
    {
      title: '操作',
      key: 'action',
      render: (text: string, record: IModelData) => (
        <Space size="small">
          <Button 
            type="text" 
            icon={<EditOutlined />} 
            onClick={() => handleEditModel(record)}
          />
          <Popconfirm
            title="确定要删除这个模型吗？"
            onConfirm={() => handleDeleteModel(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="text" danger icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={4}>模型配置</Title>
        <Button 
          type="primary" 
          icon={<PlusOutlined />} 
          onClick={showAddModal}
        >
          添加模型
        </Button>
      </div>
      <Divider style={{ margin: '16px 0' }} />
      
      <Table
        columns={columns}
        dataSource={models}
        rowKey="id"
        loading={loading}
      />

      {/* 新增/编辑模型对话框 */}
      <Modal
        title={editingModel ? '编辑模型' : '添加模型'}
        open={modalVisible}
        onOk={form.submit}
        onCancel={handleCancel}
        confirmLoading={confirmLoading}
        width={700}
        footer={[
          <Button key="test" onClick={handleTestConnection} loading={testingConnection}>
            测试连接
          </Button>,
          <Button key="cancel" onClick={handleCancel}>
            取消
          </Button>,
          <Button key="submit" type="primary" onClick={form.submit} loading={confirmLoading}>
            保存
          </Button>,
        ]}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleAddOrUpdateModel}
          initialValues={{
            type: 'cloud',
            status: 'active',
            contextLength: 4096,
            supportsFunctionCalling: false,
            parameters: {
              temperature: 0.7,
              topP: 1,
              maxTokens: 1000,
            }
          }}
        >
          <Form.Item
            name="name"
            label="模型名称"
            rules={[{ required: true, message: '请输入模型名称' }]}
          >
            <Input placeholder="例如: GPT-4, Claude 3" />
          </Form.Item>
          
          <Form.Item
            name="provider"
            label="提供商"
            rules={[{ required: true, message: '请输入提供商' }]}
          >
            <Input placeholder="例如: OpenAI, Anthropic, 百度, 阿里" />
          </Form.Item>
          
          <Form.Item
            name="type"
            label="类型"
            rules={[{ required: true }]}
          >
            <Select>
              <Option value="local">本地</Option>
              <Option value="cloud">云端</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="status"
            label="状态"
            valuePropName="checked"
          >
            <Switch checkedChildren="启用" unCheckedChildren="禁用" />
          </Form.Item>
          
          <Form.Item
            name="contextLength"
            label="上下文长度 (tokens)"
            rules={[{ required: true, message: '请输入上下文长度' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
          
          <Form.Item
            name="supportsFunctionCalling"
            label="支持函数调用"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
          
          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) => prevValues.type !== currentValues.type}
          >
            {({ getFieldValue }) => 
              getFieldValue('type') === 'cloud' ? (
                <>
                  <Form.Item
                    name="apiKey"
                    label={
                      <span>
                        API密钥 
                        <Tooltip title="请提供模型API访问所需的密钥">
                          <InfoCircleOutlined style={{ marginLeft: 4 }} />
                        </Tooltip>
                      </span>
                    }
                    rules={[{ required: true, message: '请输入API密钥' }]}
                  >
                    <Input.Password placeholder="您的API密钥" />
                  </Form.Item>
                  
                  <Form.Item
                    name="baseUrl"
                    label="基础URL"
                  >
                    <Input placeholder="API基础URL（可选）" />
                  </Form.Item>
                </>
              ) : null
            }
          </Form.Item>
          
          <Collapse>
            <Panel header="高级参数" key="1">
              <Form.Item label="温度" name={['parameters', 'temperature']}>
                <InputNumber min={0} max={2} step={0.1} style={{ width: '100%' }} />
              </Form.Item>
              
              <Form.Item label="Top P" name={['parameters', 'topP']}>
                <InputNumber min={0} max={1} step={0.05} style={{ width: '100%' }} />
              </Form.Item>
              
              <Form.Item label="最大输出tokens" name={['parameters', 'maxTokens']}>
                <InputNumber min={1} style={{ width: '100%' }} />
              </Form.Item>
            </Panel>
          </Collapse>
        </Form>
      </Modal>
    </div>
  );
};

export default ModelSettings; 