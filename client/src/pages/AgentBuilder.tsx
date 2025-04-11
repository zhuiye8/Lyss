import React, { useState } from 'react';
import { 
  Typography, 
  Card, 
  Tabs, 
  Form, 
  Input, 
  Select, 
  Switch, 
  Button, 
  Row, 
  Col,
  Divider,
  Space,
  Tag,
  List,
  Avatar,
  Modal,
  message
} from 'antd';
import {
  RobotOutlined,
  SettingOutlined,
  DatabaseOutlined,
  ApiOutlined,
  PlusOutlined,
  DeleteOutlined,
  EditOutlined,
  MessageOutlined,
  CodeOutlined,
  SaveOutlined,
  PlayCircleOutlined
} from '@ant-design/icons';

const { Title, Paragraph, Text } = Typography;
const { TabPane } = Tabs;
const { TextArea } = Input;
const { Option } = Select;

// 模拟数据
const mockModels = [
  { label: 'GPT-4', value: 'gpt-4' },
  { label: 'GPT-3.5-Turbo', value: 'gpt-3.5-turbo' },
  { label: 'Claude 3 Opus', value: 'claude-3-opus' },
  { label: 'Claude 3 Sonnet', value: 'claude-3-sonnet' },
  { label: 'Llama 3', value: 'llama-3' },
  { label: 'Mistral', value: 'mistral' },
];

const mockKnowledgeBases = [
  { label: '产品手册', value: 'product-manual' },
  { label: 'API文档', value: 'api-docs' },
  { label: '常见问题', value: 'faq' },
];

const mockTools = [
  { name: '网络搜索', id: 'web-search', description: '允许智能体在互联网上搜索信息' },
  { name: '代码解释器', id: 'code-interpreter', description: '允许智能体执行Python代码' },
  { name: '天气查询', id: 'weather', description: '允许智能体查询天气信息' },
  { name: '日程管理', id: 'calendar', description: '允许智能体管理日程和提醒' },
  { name: '数据分析', id: 'data-analysis', description: '允许智能体分析数据集' },
  { name: '图像生成', id: 'image-gen', description: '允许智能体生成图像' },
];

const mockAgents = [
  { 
    id: '1', 
    name: '客服助手', 
    description: '为客户提供支持和帮助的智能体', 
    model: 'gpt-4',
    createdAt: '2025-03-10 14:30',
    updatedAt: '2025-04-05 09:20',
  },
  { 
    id: '2', 
    name: '代码助手', 
    description: '帮助开发人员编写和调试代码的智能体', 
    model: 'claude-3-opus',
    createdAt: '2025-03-15 10:45',
    updatedAt: '2025-04-02 16:30',
  },
  { 
    id: '3', 
    name: '数据分析师', 
    description: '分析数据并提供见解的智能体', 
    model: 'gpt-4',
    createdAt: '2025-03-22 09:15',
    updatedAt: '2025-04-01 11:20',
  },
];

const AgentBuilder: React.FC = () => {
  const [activeKey, setActiveKey] = useState<string>('my-agents');
  const [isBuilding, setIsBuilding] = useState<boolean>(false);
  const [currentAgent, setCurrentAgent] = useState<any>(null);
  const [selectedTools, setSelectedTools] = useState<string[]>([]);
  const [isTestModalVisible, setIsTestModalVisible] = useState<boolean>(false);
  const [testMessages, setTestMessages] = useState<any[]>([]);
  const [messageInput, setMessageInput] = useState<string>('');
  const [form] = Form.useForm();

  // 开始创建新智能体
  const handleCreateNewAgent = () => {
    setCurrentAgent(null);
    form.resetFields();
    setSelectedTools([]);
    setIsBuilding(true);
    setActiveKey('builder');
  };

  // 编辑现有智能体
  const handleEditAgent = (agent: any) => {
    setCurrentAgent(agent);
    form.setFieldsValue({
      name: agent.name,
      description: agent.description,
      model: agent.model,
      // 其他字段...
    });
    setSelectedTools(['web-search', 'code-interpreter']); // 模拟已选工具
    setIsBuilding(true);
    setActiveKey('builder');
  };

  // 保存智能体
  const handleSaveAgent = () => {
    form.validateFields().then(values => {
      console.log('Form values:', values);
      console.log('Selected tools:', selectedTools);
      
      message.success(currentAgent ? '智能体更新成功' : '智能体创建成功');
      setIsBuilding(false);
      setActiveKey('my-agents');
    });
  };

  // 取消构建
  const handleCancelBuild = () => {
    Modal.confirm({
      title: '确定要取消构建？',
      content: '所有未保存的更改将会丢失。',
      onOk: () => {
        setIsBuilding(false);
        setActiveKey('my-agents');
      },
    });
  };

  // 处理工具选择
  const handleToolSelect = (toolId: string) => {
    if (selectedTools.includes(toolId)) {
      setSelectedTools(selectedTools.filter(id => id !== toolId));
    } else {
      setSelectedTools([...selectedTools, toolId]);
    }
  };

  // 打开测试模态框
  const handleTestAgent = () => {
    setTestMessages([
      { type: 'system', content: '这是一个测试对话，您可以测试智能体的功能。' },
    ]);
    setIsTestModalVisible(true);
  };

  // 发送测试消息
  const handleSendTestMessage = () => {
    if (!messageInput.trim()) return;
    
    // 添加用户消息
    setTestMessages([
      ...testMessages,
      { type: 'user', content: messageInput },
    ]);
    
    // 模拟智能体回复
    setTimeout(() => {
      setTestMessages(prev => [
        ...prev,
        { 
          type: 'agent', 
          content: `这是对"${messageInput}"的模拟回复。在实际应用中，这里将显示来自AI模型的响应。`
        },
      ]);
    }, 1000);
    
    setMessageInput('');
  };

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Title level={2}>智能体构建</Title>
        <Paragraph>构建和管理您的AI智能体，可以为特定任务定制功能。</Paragraph>
      </div>

      <Tabs activeKey={activeKey} onChange={setActiveKey}>
        <TabPane 
          tab={
            <span>
              <RobotOutlined />
              我的智能体
            </span>
          } 
          key="my-agents"
        >
          <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'flex-end' }}>
            <Button 
              type="primary" 
              icon={<PlusOutlined />} 
              onClick={handleCreateNewAgent}
            >
              创建新智能体
            </Button>
          </div>
          
          <Row gutter={[16, 16]}>
            {mockAgents.map(agent => (
              <Col xs={24} sm={12} md={8} key={agent.id}>
                <Card
                  actions={[
                    <EditOutlined key="edit" onClick={() => handleEditAgent(agent)} />,
                    <MessageOutlined key="test" onClick={handleTestAgent} />,
                    <DeleteOutlined key="delete" onClick={() => message.success('删除成功')} />,
                  ]}
                >
                  <div style={{ display: 'flex', alignItems: 'center', marginBottom: 12 }}>
                    <Avatar 
                      icon={<RobotOutlined />} 
                      style={{ backgroundColor: '#1890ff', marginRight: 12 }} 
                    />
                    <div>
                      <div style={{ fontWeight: 'bold', fontSize: 16 }}>{agent.name}</div>
                      <div style={{ fontSize: 12, color: '#666' }}>模型: {agent.model}</div>
                    </div>
                  </div>
                  <Paragraph ellipsis={{ rows: 2 }}>{agent.description}</Paragraph>
                  <div style={{ fontSize: 12, color: '#999', marginTop: 8 }}>
                    更新于: {agent.updatedAt}
                  </div>
                </Card>
              </Col>
            ))}
          </Row>
        </TabPane>

        <TabPane 
          tab={
            <span>
              <SettingOutlined />
              构建器
            </span>
          } 
          key="builder"
          disabled={!isBuilding}
        >
          <Form
            form={form}
            layout="vertical"
            initialValues={{
              model: 'gpt-4',
              temperature: 0.7,
              maxTokens: 2048,
              enableKnowledgeBase: false,
            }}
          >
            <Card title="基本信息" style={{ marginBottom: 16 }}>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="name"
                    label="智能体名称"
                    rules={[{ required: true, message: '请输入智能体名称' }]}
                  >
                    <Input placeholder="例如：客服助手、代码助手" />
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="model"
                    label="选择模型"
                    rules={[{ required: true, message: '请选择模型' }]}
                  >
                    <Select>
                      {mockModels.map(model => (
                        <Option key={model.value} value={model.value}>{model.label}</Option>
                      ))}
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
              
              <Form.Item
                name="description"
                label="描述"
                rules={[{ required: true, message: '请输入智能体描述' }]}
              >
                <TextArea 
                  placeholder="请描述智能体的功能和用途..." 
                  rows={3}
                />
              </Form.Item>
            </Card>
            
            <Card title="系统提示词" style={{ marginBottom: 16 }}>
              <Form.Item
                name="systemPrompt"
                label="系统提示词"
                rules={[{ required: true, message: '请输入系统提示词' }]}
                initialValue="你是一个有帮助的AI助手，你会尽可能地回答用户的问题。"
              >
                <TextArea 
                  placeholder="请输入系统提示词，定义智能体的角色和行为..." 
                  rows={6}
                />
              </Form.Item>
              
              <div style={{ marginTop: 8 }}>
                <Text type="secondary">系统提示词用于指导AI的行为和角色定位，对智能体的表现有重要影响。</Text>
              </div>
            </Card>
            
            <Card title="模型参数" style={{ marginBottom: 16 }}>
              <Row gutter={16}>
                <Col span={12}>
                  <Form.Item
                    name="temperature"
                    label="温度 (Temperature)"
                    help="控制输出的随机性，值越低回答越确定，值越高回答越多样"
                  >
                    <Select>
                      <Option value={0}>0 - 精确</Option>
                      <Option value={0.3}>0.3 - 平衡</Option>
                      <Option value={0.7}>0.7 - 创造性</Option>
                      <Option value={1}>1 - 多样性</Option>
                    </Select>
                  </Form.Item>
                </Col>
                <Col span={12}>
                  <Form.Item
                    name="maxTokens"
                    label="最大输出长度 (Max Tokens)"
                    help="每次回复最多生成的token数量"
                  >
                    <Select>
                      <Option value={1024}>1024 - 短回复</Option>
                      <Option value={2048}>2048 - 中等长度</Option>
                      <Option value={4096}>4096 - 长回复</Option>
                      <Option value={8192}>8192 - 非常长</Option>
                    </Select>
                  </Form.Item>
                </Col>
              </Row>
            </Card>
            
            <Card title="知识库" style={{ marginBottom: 16 }}>
              <Form.Item
                name="enableKnowledgeBase"
                valuePropName="checked"
              >
                <Switch /> <span style={{ marginLeft: 8 }}>启用知识库增强</span>
              </Form.Item>
              
              <Form.Item
                name="knowledgeBase"
                label="选择知识库"
                dependencies={['enableKnowledgeBase']}
                rules={[
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (getFieldValue('enableKnowledgeBase') && !value) {
                        return Promise.reject(new Error('请选择知识库'));
                      }
                      return Promise.resolve();
                    },
                  }),
                ]}
              >
                <Select 
                  placeholder="请选择知识库" 
                  disabled={!form.getFieldValue('enableKnowledgeBase')}
                >
                  {mockKnowledgeBases.map(kb => (
                    <Option key={kb.value} value={kb.value}>{kb.label}</Option>
                  ))}
                </Select>
              </Form.Item>
              
              <div style={{ marginTop: 8 }}>
                <Text type="secondary">
                  <DatabaseOutlined /> 知识库允许智能体访问特定领域的信息，提升回答质量和准确性。
                </Text>
              </div>
            </Card>
            
            <Card title="工具" style={{ marginBottom: 16 }}>
              <Paragraph>选择您希望智能体能够使用的工具：</Paragraph>
              
              <List
                dataSource={mockTools}
                renderItem={tool => (
                  <List.Item 
                    key={tool.id}
                    extra={
                      <Button 
                        type={selectedTools.includes(tool.id) ? 'primary' : 'default'}
                        onClick={() => handleToolSelect(tool.id)}
                      >
                        {selectedTools.includes(tool.id) ? '已选' : '选择'}
                      </Button>
                    }
                  >
                    <List.Item.Meta
                      avatar={<Avatar icon={<ApiOutlined />} style={{ backgroundColor: '#722ed1' }} />}
                      title={tool.name}
                      description={tool.description}
                    />
                  </List.Item>
                )}
              />
              
              <div style={{ marginTop: 16 }}>
                <Text type="secondary">
                  <CodeOutlined /> 工具使智能体能够执行特定任务和访问外部资源，显著扩展其功能。
                </Text>
              </div>
            </Card>
            
            <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 24, marginBottom: 24 }}>
              <Button onClick={handleCancelBuild}>取消</Button>
              <Space>
                <Button icon={<PlayCircleOutlined />} onClick={handleTestAgent}>测试</Button>
                <Button type="primary" icon={<SaveOutlined />} onClick={handleSaveAgent}>保存</Button>
              </Space>
            </div>
          </Form>
        </TabPane>
      </Tabs>

      {/* 测试对话模态框 */}
      <Modal
        title="测试智能体"
        open={isTestModalVisible}
        onCancel={() => setIsTestModalVisible(false)}
        footer={null}
        width={700}
      >
        <div style={{ height: 400, overflowY: 'scroll', marginBottom: 16, padding: 8, border: '1px solid #f0f0f0' }}>
          {testMessages.map((msg, index) => (
            <div 
              key={index} 
              style={{ 
                marginBottom: 12,
                textAlign: msg.type === 'user' ? 'right' : 'left',
              }}
            >
              {msg.type === 'system' && (
                <div style={{ 
                  padding: '8px 12px', 
                  background: '#f0f0f0', 
                  borderRadius: 8,
                  color: '#666',
                  display: 'inline-block',
                  maxWidth: '80%',
                }}>
                  {msg.content}
                </div>
              )}
              
              {msg.type === 'user' && (
                <div style={{ 
                  padding: '8px 12px', 
                  background: '#1890ff', 
                  color: 'white',
                  borderRadius: 8,
                  display: 'inline-block',
                  maxWidth: '80%',
                }}>
                  {msg.content}
                </div>
              )}
              
              {msg.type === 'agent' && (
                <div style={{ 
                  padding: '8px 12px', 
                  background: '#f9f9f9', 
                  border: '1px solid #d9d9d9',
                  borderRadius: 8,
                  display: 'inline-block',
                  maxWidth: '80%',
                }}>
                  {msg.content}
                </div>
              )}
            </div>
          ))}
        </div>
        
        <div style={{ display: 'flex' }}>
          <Input 
            value={messageInput}
            onChange={e => setMessageInput(e.target.value)}
            onPressEnter={handleSendTestMessage}
            placeholder="输入测试消息..."
          />
          <Button type="primary" onClick={handleSendTestMessage} style={{ marginLeft: 8 }}>
            发送
          </Button>
        </div>
      </Modal>
    </div>
  );
};

export default AgentBuilder; 