import React, { useState, useEffect, useRef } from 'react';
import {
  Card,
  Input,
  Button,
  Typography,
  Spin,
  List,
  Avatar,
  Divider,
  Space,
  Tabs,
  Badge,
  Table,
  Tag,
  message,
  Collapse,
  Form,
  Select,
  Row,
  Col
} from 'antd';
import {
  SendOutlined,
  UserOutlined,
  RobotOutlined,
  ArrowLeftOutlined,
  CodeOutlined,
  BugOutlined,
  DashboardOutlined,
  ThunderboltOutlined
} from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import ReactJson from 'react-json-view';
import { 
  getAgent,
  testAgent,
  getTestHistory,
  createConversation,
  sendMessage,
  getConversation
} from '../../services/agentService';
import { 
  IAgentFormState, 
  IAgentMessage, 
  ITestResult,
  IAgentConversation,
  IToolCall
} from '../../types/agent';
import ChatInterface from '../../components/AgentBuilder/ChatInterface';
import SuggestionPrompts from '../../components/AgentBuilder/SuggestionPrompts';
import AIWelcome from '../../components/Welcome';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { TabPane } = Tabs;
const { Panel } = Collapse;
const { Option } = Select;

const AgentTester: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();

  // 状态
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const [inputValue, setInputValue] = useState('');
  const [agent, setAgent] = useState<IAgentFormState | null>(null);
  const [conversation, setConversation] = useState<IAgentConversation | null>(null);
  const [testHistory, setTestHistory] = useState<ITestResult[]>([]);
  const [activeTest, setActiveTest] = useState<ITestResult | null>(null);
  const [variables, setVariables] = useState<Record<string, any>>({});
  const [testLoading, setTestLoading] = useState(false);
  const [historyLoading, setHistoryLoading] = useState(false);
  const [showWelcome, setShowWelcome] = useState(true);

  // 获取智能体信息
  useEffect(() => {
    if (id) {
      fetchAgentDetails();
      fetchTestHistory();
    }
  }, [id]);

  // 获取智能体详情
  const fetchAgentDetails = async () => {
    try {
      setLoading(true);
      const data = await getAgent(id!);
      setAgent(data);

      // 提取智能体变量
      const extractedVariables: Record<string, any> = {};
      if (data.type === 'chat' && data.chatConfig?.promptTemplate?.variables) {
        data.chatConfig.promptTemplate.variables.forEach(variable => {
          extractedVariables[variable] = '';
        });
      } else if (data.type === 'flow' && data.flowConfig?.variables) {
        data.flowConfig.variables.forEach(variable => {
          extractedVariables[variable] = '';
        });
      }

      setVariables(extractedVariables);
      form.setFieldsValue(extractedVariables);
      
      // 创建新对话
      await createNewConversation();
    } catch (error) {
      console.error('获取智能体详情失败:', error);
      message.error('获取智能体详情失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取测试历史
  const fetchTestHistory = async () => {
    try {
      setHistoryLoading(true);
      const response = await getTestHistory(id!);
      setTestHistory(response.data);
    } catch (error) {
      console.error('获取测试历史失败:', error);
    } finally {
      setHistoryLoading(false);
    }
  };

  // 创建新对话
  const createNewConversation = async () => {
    try {
      const conversationResponse = await createConversation(id!);
      setConversation(conversationResponse);
    } catch (error) {
      console.error('创建对话失败:', error);
      message.error('创建对话失败');
    }
  };

  // 运行智能体测试
  const handleRunTest = async () => {
    try {
      setTestLoading(true);
      const updatedVariables = form.getFieldsValue();
      setVariables(updatedVariables);
      
      const result = await testAgent(id!, inputValue, updatedVariables);
      setActiveTest(result);
      
      // 刷新测试历史
      await fetchTestHistory();
    } catch (error) {
      console.error('测试失败:', error);
      message.error('测试失败');
    } finally {
      setTestLoading(false);
    }
  };

  // 查看历史测试
  const handleViewTest = (test: ITestResult) => {
    setActiveTest(test);
  };

  // 渲染测试结果
  const renderTestResult = () => {
    if (!activeTest) return null;
    
    return (
      <Card style={{ marginBottom: 20 }}>
        <Row gutter={16}>
          <Col span={12}>
            <Title level={5}>测试详情</Title>
            <Divider />
            <p><strong>输入:</strong> {activeTest.input}</p>
            <p><strong>输出:</strong> {activeTest.output}</p>
            <p>
              <strong>状态:</strong>{' '}
              {activeTest.success ? 
                <Tag color="success">成功</Tag> : 
                <Tag color="error">失败</Tag>
              }
            </p>
            {activeTest.error && (
              <p><strong>错误:</strong> {activeTest.error}</p>
            )}
          </Col>
          <Col span={12}>
            <Title level={5}>性能指标</Title>
            <Divider />
            <p><strong>处理时间:</strong> {activeTest.duration.toFixed(2)}秒</p>
            <p><strong>Token用量:</strong></p>
            <ul>
              <li>输入Token: {activeTest.tokenUsage.prompt}</li>
              <li>输出Token: {activeTest.tokenUsage.completion}</li>
              <li>总Token: {activeTest.tokenUsage.total}</li>
            </ul>
            <p><strong>使用模型:</strong> {activeTest.modelId}</p>
          </Col>
        </Row>
      </Card>
    );
  };

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <Space style={{ marginBottom: 16 }}>
          <Button 
            icon={<ArrowLeftOutlined />} 
            onClick={() => navigate(`/agents/edit/${id}`)}
          >
            返回编辑
          </Button>
          <Title level={4} style={{ margin: 0 }}>
            {agent?.name || '智能体'} - 测试模式
          </Title>
        </Space>
        <Paragraph>
          在此页面您可以测试智能体的表现、查看思维链和运行日志，以及查阅历史测试记录。
        </Paragraph>
      </Card>

      <Row gutter={16}>
        <Col span={16}>
          {conversation ? (
            showWelcome ? (
              <AIWelcome
                agentName={agent?.name}
                agentDescription={agent?.description}
                welcomeMessage={agent?.chatConfig?.welcomeMessage || '有什么我可以帮助你的吗？'}
                onStart={() => setShowWelcome(false)}
              />
            ) : (
              <ChatInterface 
                conversationId={conversation.id}
                initialMessages={conversation.messages}
                agentName={agent?.name}
                variables={variables}
                onMessageSent={() => {
                  // 可以在这里添加消息发送后的回调逻辑
                }}
                loading={loading}
              />
            )
          ) : (
            <Spin tip="加载对话中..." />
          )}
          
          {conversation && !showWelcome && (
            <div style={{ marginTop: '16px', marginBottom: '16px' }}>
              <SuggestionPrompts 
                onSelect={(prompt) => {
                  if (conversation) {
                    sendMessage(conversation.id, prompt, variables)
                      .then(() => getConversation(conversation.id))
                      .then(updatedConv => setConversation(updatedConv))
                      .catch(err => {
                        console.error('发送快捷提示失败:', err);
                        message.error('发送快捷提示失败');
                      });
                  }
                }}
                suggestions={[
                  '你是什么类型的智能体?',
                  '你能帮我解决什么问题?',
                  '你的功能有哪些限制?',
                  '你能接入哪些外部工具?',
                  '如何使用你的高级功能?'
                ]}
              />
            </div>
          )}
          
          {Object.keys(variables).length > 0 && (
            <Card title="变量设置" style={{ marginTop: 16 }}>
              <Form
                form={form}
                layout="vertical"
                initialValues={variables}
              >
                <Row gutter={16}>
                  {Object.keys(variables).map(key => (
                    <Col span={12} key={key}>
                      <Form.Item
                        name={key}
                        label={key}
                      >
                        <Input placeholder={`输入${key}的值`} />
                      </Form.Item>
                    </Col>
                  ))}
                </Row>
              </Form>
            </Card>
          )}
        </Col>
        
        <Col span={8}>
          <Tabs defaultActiveKey="test">
            <TabPane 
              tab={
                <span>
                  <BugOutlined />
                  运行测试
                </span>
              } 
              key="test"
            >
              <Card>
                <Title level={5}>运行性能测试</Title>
                <Input.TextArea
                  value={inputValue}
                  onChange={e => setInputValue(e.target.value)}
                  placeholder="输入测试消息..."
                  autoSize={{ minRows: 3, maxRows: 6 }}
                  style={{ marginBottom: 16 }}
                />
                <Button 
                  type="primary" 
                  icon={<ThunderboltOutlined />} 
                  onClick={handleRunTest}
                  loading={testLoading}
                  block
                >
                  运行测试
                </Button>
              </Card>
              
              {renderTestResult()}
            </TabPane>
            
            <TabPane 
              tab={
                <span>
                  <DashboardOutlined />
                  测试历史
                </span>
              } 
              key="history"
            >
              <List
                loading={historyLoading}
                itemLayout="horizontal"
                dataSource={testHistory}
                renderItem={item => (
                  <List.Item
                    actions={[
                      <Button 
                        key="view" 
                        type="link" 
                        onClick={() => handleViewTest(item)}
                      >
                        查看
                      </Button>
                    ]}
                  >
                    <List.Item.Meta
                      title={
                        <Space>
                          <span>{new Date(item.timestamp).toLocaleString()}</span>
                          {item.success ? 
                            <Tag color="success">成功</Tag> : 
                            <Tag color="error">失败</Tag>
                          }
                        </Space>
                      }
                      description={`输入: ${item.input.substring(0, 50)}${item.input.length > 50 ? '...' : ''}`}
                    />
                    <div>{(item.duration * 1000).toFixed(0)}ms</div>
                  </List.Item>
                )}
              />
            </TabPane>
          </Tabs>
        </Col>
      </Row>
    </div>
  );
};

export default AgentTester; 