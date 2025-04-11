import React, { useState, useEffect } from 'react';
import {
  Form,
  Input,
  Card,
  Button,
  Divider,
  message,
  Steps,
  Typography,
  Row,
  Col,
  Slider,
  InputNumber,
  Space,
  Tabs,
  Collapse,
  Radio,
  Tag
} from 'antd';
import {
  SaveOutlined,
  RocketOutlined,
  ArrowLeftOutlined,
  ArrowRightOutlined,
  CodeOutlined
} from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import ModelSelector from '../../components/AgentBuilder/ModelSelector';
import KnowledgeBaseSelector from '../../components/AgentBuilder/KnowledgeBaseSelector';
import ToolSelector from '../../components/AgentBuilder/ToolSelector';
import PromptTemplateSelector from '../../components/AgentBuilder/PromptTemplateSelector';
import { 
  getAgent, 
  createAgent, 
  updateAgent, 
  publishAgent
} from '../../services/agentService';
import { IAgentFormState, IChatAgentConfig, IPromptTemplate } from '../../types/agent';
import { 
  XProvider, 
  Prompts, 
  ThoughtChain 
} from '@ant-design/x';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { Step } = Steps;
const { TabPane } = Tabs;
const { Panel } = Collapse;

const steps = [
  {
    title: '基础设置',
    description: '设置智能体名称和模型'
  },
  {
    title: '知识库和工具',
    description: '配置知识库和工具接入'
  },
  {
    title: '提示词设置',
    description: '设计智能体的提示词'
  },
  {
    title: '参数与测试',
    description: '调优参数并测试效果'
  }
];

const ChatAgentBuilder: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const isEditing = !!id;

  const [currentStep, setCurrentStep] = useState(0);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [publishing, setPublishing] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState<IPromptTemplate | null>(null);
  const [formValues, setFormValues] = useState<IAgentFormState>({
    name: '',
    description: '',
    type: 'chat',
    visibility: 'private',
    modelId: '',
    status: 'draft',
    chatConfig: {
      welcomeMessage: '有什么我可以帮您的吗？',
      systemPrompt: '你是一个智能助手，回答用户的问题。',
      tools: [],
      knowledgeBases: [],
      parameters: {
        temperature: 0.7,
        topP: 0.9,
        maxTokens: 1000,
        presencePenalty: 0,
        frequencyPenalty: 0
      }
    }
  });

  // 获取智能体信息
  useEffect(() => {
    if (isEditing) {
      fetchAgentDetails();
    }
  }, [id]);

  const fetchAgentDetails = async () => {
    try {
      setLoading(true);
      const data = await getAgent(id!);
      
      if (data.type !== 'chat') {
        message.error('当前智能体不是对话式应用');
        navigate('/agents');
        return;
      }
      
      setFormValues(data);
      form.setFieldsValue({
        ...data,
        promptTemplateId: data.chatConfig?.promptTemplate?.id,
        temperature: data.chatConfig?.parameters?.temperature || 0.7,
        topP: data.chatConfig?.parameters?.topP || 0.9,
        maxTokens: data.chatConfig?.parameters?.maxTokens || 1000,
        presencePenalty: data.chatConfig?.parameters?.presencePenalty || 0,
        frequencyPenalty: data.chatConfig?.parameters?.frequencyPenalty || 0
      });
      
      if (data.chatConfig?.promptTemplate) {
        setSelectedTemplate(data.chatConfig.promptTemplate);
      }

    } catch (error) {
      console.error('获取智能体详情失败:', error);
      message.error('获取智能体详情失败');
    } finally {
      setLoading(false);
    }
  };

  // 保存表单数据
  const handleSave = async () => {
    try {
      const values = await form.validateFields();
      setSaving(true);

      // 构建保存的数据结构
      const chatConfig: Partial<IChatAgentConfig> = {
        ...formValues.chatConfig,
        welcomeMessage: values.welcomeMessage,
        systemPrompt: values.systemPrompt,
        knowledgeBases: values.knowledgeBases || [],
        tools: values.tools || [],
        promptTemplate: selectedTemplate || undefined,
        parameters: {
          temperature: values.temperature,
          topP: values.topP,
          maxTokens: values.maxTokens,
          presencePenalty: values.presencePenalty,
          frequencyPenalty: values.frequencyPenalty
        }
      };

      const agentData: IAgentFormState = {
        name: values.name,
        description: values.description,
        type: 'chat',
        visibility: values.visibility,
        modelId: values.modelId,
        status: formValues.status,
        chatConfig
      };

      let result;
      if (isEditing) {
        result = await updateAgent(id!, agentData);
        message.success('智能体更新成功');
      } else {
        result = await createAgent(agentData);
        message.success('智能体创建成功');
        // 创建成功后跳转到编辑页
        navigate(`/agents/edit/${result.id}`);
      }

      setFormValues(result);
    } catch (error) {
      console.error('保存智能体失败:', error);
      message.error('保存智能体失败，请检查表单填写');
    } finally {
      setSaving(false);
    }
  };

  // 发布智能体
  const handlePublish = async () => {
    try {
      await handleSave();
      setPublishing(true);
      await publishAgent(id!);
      message.success('智能体已发布');
      // 刷新数据
      await fetchAgentDetails();
    } catch (error) {
      console.error('发布智能体失败:', error);
      message.error('发布智能体失败');
    } finally {
      setPublishing(false);
    }
  };

  // 处理步骤变化
  const handleStepChange = (step: number) => {
    form.validateFields().then(() => {
      setCurrentStep(step);
    }).catch(error => {
      // 如果当前表单验证失败，不跳转
      console.error('表单验证失败:', error);
    });
  };

  // 下一步
  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      form.validateFields().then(() => {
        setCurrentStep(currentStep + 1);
      }).catch(error => {
        console.error('表单验证失败:', error);
      });
    }
  };

  // 上一步
  const handlePrev = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  // 处理提示词模板选择
  const handleTemplateSelect = (template: IPromptTemplate) => {
    setSelectedTemplate(template);
    // 更新表单中相关的字段
    if (template) {
      form.setFieldsValue({
        systemPrompt: template.content
      });
    }
  };

  // 使用Ant Design X组件渲染提示词编辑器
  const renderPromptEditor = () => {
    const templates = selectedTemplate ? [selectedTemplate] : [];
    
    return (
      <XProvider>
        <Prompts
          templates={templates.map(tpl => ({
            id: tpl.id,
            title: tpl.name || '默认模板',
            content: tpl.content,
            variables: tpl.variables.map(v => ({ name: v })),
            description: tpl.description || ''
          }))}
          activeId={selectedTemplate?.id}
          onSelect={(templateId) => {
            const template = templates.find(t => t.id === templateId);
            if (template) {
              setSelectedTemplate(template);
              form.setFieldsValue({
                systemPrompt: template.content
              });
            }
          }}
          editable={true}
          onChange={(updatedTemplate) => {
            if (selectedTemplate && updatedTemplate.id === selectedTemplate.id) {
              const newTemplate = { 
                ...selectedTemplate,
                content: updatedTemplate.content,
                variables: updatedTemplate.variables.map(v => v.name)
              };
              setSelectedTemplate(newTemplate);
              form.setFieldsValue({
                systemPrompt: updatedTemplate.content
              });
            }
          }}
        />
      </XProvider>
    );
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
        <Space>
          <Button 
            icon={<ArrowLeftOutlined />} 
            onClick={() => navigate('/agents')}
          >
            返回
          </Button>
          <Title level={3} style={{ margin: 0 }}>
            {isEditing ? '编辑对话式智能体' : '创建对话式智能体'}
          </Title>
        </Space>
        <Space>
          <Button 
            type="primary" 
            icon={<SaveOutlined />} 
            onClick={handleSave} 
            loading={saving}
          >
            保存
          </Button>
          {isEditing && (
            <Button 
              type="primary" 
              danger 
              icon={<RocketOutlined />} 
              onClick={handlePublish} 
              loading={publishing}
            >
              发布
            </Button>
          )}
        </Space>
      </div>

      <Steps 
        current={currentStep} 
        onChange={handleStepChange} 
        items={steps}
        style={{ marginBottom: 24 }}
      />

      <Form
        form={form}
        layout="vertical"
        initialValues={{
          ...formValues,
          welcomeMessage: formValues.chatConfig?.welcomeMessage,
          systemPrompt: formValues.chatConfig?.systemPrompt,
          knowledgeBases: formValues.chatConfig?.knowledgeBases || [],
          tools: formValues.chatConfig?.tools || [],
          promptTemplateId: formValues.chatConfig?.promptTemplate?.id,
          temperature: formValues.chatConfig?.parameters?.temperature || 0.7,
          topP: formValues.chatConfig?.parameters?.topP || 0.9,
          maxTokens: formValues.chatConfig?.parameters?.maxTokens || 1000,
          presencePenalty: formValues.chatConfig?.parameters?.presencePenalty || 0,
          frequencyPenalty: formValues.chatConfig?.parameters?.frequencyPenalty || 0
        }}
      >
        <div style={{ display: currentStep === 0 ? 'block' : 'none' }}>
          <Card title="基础信息" style={{ marginBottom: 16 }}>
            <Form.Item
              name="name"
              label="智能体名称"
              rules={[{ required: true, message: '请输入智能体名称' }]}
            >
              <Input placeholder="例如: 客服助手、销售顾问" />
            </Form.Item>
            
            <Form.Item
              name="description"
              label="描述"
              rules={[{ required: true, message: '请输入智能体描述' }]}
            >
              <TextArea 
                placeholder="描述智能体的功能和用途" 
                autoSize={{ minRows: 3, maxRows: 6 }}
              />
            </Form.Item>
            
            <Form.Item
              name="visibility"
              label="可见性"
              rules={[{ required: true }]}
            >
              <Radio.Group>
                <Radio value="public">公开</Radio>
                <Radio value="private">私有</Radio>
              </Radio.Group>
            </Form.Item>
          </Card>
          
          <Card title="对话设置">
            <Form.Item
              name="welcomeMessage"
              label="欢迎语"
              rules={[{ required: true, message: '请输入欢迎语' }]}
            >
              <TextArea 
                placeholder="用户首次访问时，智能体的欢迎语" 
                autoSize={{ minRows: 2, maxRows: 4 }}
              />
            </Form.Item>
            
            <Form.Item
              name="modelId"
              label="语言模型"
              rules={[{ required: true, message: '请选择语言模型' }]}
            >
              <ModelSelector />
            </Form.Item>
          </Card>
        </div>

        <div style={{ display: currentStep === 1 ? 'block' : 'none' }}>
          <Card title="知识库配置" style={{ marginBottom: 16 }}>
            <Paragraph>
              选择智能体可以访问的知识库，允许智能体基于知识库的内容回答问题。
            </Paragraph>
            <Form.Item
              name="knowledgeBases"
              label="知识库"
            >
              <KnowledgeBaseSelector />
            </Form.Item>
          </Card>
          
          <Card title="工具配置">
            <Paragraph>
              选择智能体可以使用的工具，使智能体能够执行特定操作。
            </Paragraph>
            <Form.Item
              name="tools"
              label="工具"
            >
              <ToolSelector />
            </Form.Item>
          </Card>
        </div>

        <div style={{ display: currentStep === 2 ? 'block' : 'none' }}>
          <Card title="提示词设置">
            <Tabs defaultActiveKey="template">
              <TabPane tab="使用模板" key="template">
                <Form.Item
                  name="promptTemplateId"
                  label="提示词模板"
                >
                  <PromptTemplateSelector 
                    onSelect={handleTemplateSelect}
                    value={selectedTemplate?.id}
                  />
                </Form.Item>
                
                {selectedTemplate && (
                  <Card 
                    type="inner" 
                    title="模板预览" 
                    style={{ marginBottom: 16, backgroundColor: '#f9f9f9' }}
                  >
                    <Paragraph>
                      <pre style={{ whiteSpace: 'pre-wrap', fontFamily: 'monospace' }}>
                        {selectedTemplate.content}
                      </pre>
                    </Paragraph>
                    <Divider />
                    <div>
                      <Text strong>变量列表:</Text>
                      <div style={{ marginTop: 8 }}>
                        {selectedTemplate.variables.map((variable, index) => (
                          <Tag key={index} color="blue">{`{{${variable}}}`}</Tag>
                        ))}
                      </div>
                    </div>
                  </Card>
                )}
              </TabPane>
              <TabPane tab="自定义提示词" key="custom">
                <Form.Item
                  name="systemPrompt"
                  label="系统提示词"
                  rules={[{ required: true, message: '请输入系统提示词' }]}
                  extra="设定智能体的行为、性格、限制和专业领域等"
                >
                  {renderPromptEditor()}
                </Form.Item>
              </TabPane>
            </Tabs>
          </Card>
        </div>

        <div style={{ display: currentStep === 3 ? 'block' : 'none' }}>
          <Card title="模型参数" style={{ marginBottom: 16 }}>
            <Collapse defaultActiveKey={['1']} ghost>
              <Panel header="温度 (Temperature)" key="1">
                <Paragraph>
                  控制生成文本的随机性和创造性。较高的值使输出更随机，较低的值使输出更确定和集中。
                </Paragraph>
                <Form.Item name="temperature">
                  <Row>
                    <Col span={19}>
                      <Slider
                        min={0}
                        max={2}
                        step={0.1}
                        onChange={value => form.setFieldsValue({ temperature: value })}
                      />
                    </Col>
                    <Col span={4} offset={1}>
                      <InputNumber
                        min={0}
                        max={2}
                        step={0.1}
                        style={{ width: '100%' }}
                        onChange={value => form.setFieldsValue({ temperature: value })}
                      />
                    </Col>
                  </Row>
                </Form.Item>
              </Panel>
              
              <Panel header="Top P (采样概率)" key="2">
                <Paragraph>
                  核采样，控制模型考虑的词汇范围。较低的值使输出更保守，较高的值使模型考虑更多可能性。
                </Paragraph>
                <Form.Item name="topP">
                  <Row>
                    <Col span={19}>
                      <Slider
                        min={0}
                        max={1}
                        step={0.05}
                        onChange={value => form.setFieldsValue({ topP: value })}
                      />
                    </Col>
                    <Col span={4} offset={1}>
                      <InputNumber
                        min={0}
                        max={1}
                        step={0.05}
                        style={{ width: '100%' }}
                        onChange={value => form.setFieldsValue({ topP: value })}
                      />
                    </Col>
                  </Row>
                </Form.Item>
              </Panel>
              
              <Panel header="最大输出长度 (Max Tokens)" key="3">
                <Paragraph>
                  控制生成文本的最大长度，防止模型输出过长。
                </Paragraph>
                <Form.Item name="maxTokens">
                  <Row>
                    <Col span={19}>
                      <Slider
                        min={100}
                        max={4000}
                        step={100}
                        onChange={value => form.setFieldsValue({ maxTokens: value })}
                      />
                    </Col>
                    <Col span={4} offset={1}>
                      <InputNumber
                        min={100}
                        max={4000}
                        step={100}
                        style={{ width: '100%' }}
                        onChange={value => form.setFieldsValue({ maxTokens: value })}
                      />
                    </Col>
                  </Row>
                </Form.Item>
              </Panel>
              
              <Panel header="高级参数" key="4">
                <Paragraph>
                  <Text strong>存在惩罚 (Presence Penalty):</Text> 降低模型重复已经出现过的内容的可能性。
                </Paragraph>
                <Form.Item name="presencePenalty">
                  <Row>
                    <Col span={19}>
                      <Slider
                        min={-2}
                        max={2}
                        step={0.1}
                        onChange={value => form.setFieldsValue({ presencePenalty: value })}
                      />
                    </Col>
                    <Col span={4} offset={1}>
                      <InputNumber
                        min={-2}
                        max={2}
                        step={0.1}
                        style={{ width: '100%' }}
                        onChange={value => form.setFieldsValue({ presencePenalty: value })}
                      />
                    </Col>
                  </Row>
                </Form.Item>
                
                <Paragraph style={{ marginTop: 16 }}>
                  <Text strong>频率惩罚 (Frequency Penalty):</Text> 降低模型重复使用相同词语的可能性。
                </Paragraph>
                <Form.Item name="frequencyPenalty">
                  <Row>
                    <Col span={19}>
                      <Slider
                        min={-2}
                        max={2}
                        step={0.1}
                        onChange={value => form.setFieldsValue({ frequencyPenalty: value })}
                      />
                    </Col>
                    <Col span={4} offset={1}>
                      <InputNumber
                        min={-2}
                        max={2}
                        step={0.1}
                        style={{ width: '100%' }}
                        onChange={value => form.setFieldsValue({ frequencyPenalty: value })}
                      />
                    </Col>
                  </Row>
                </Form.Item>
              </Panel>
            </Collapse>
          </Card>
          
          <Card title="测试智能体">
            <Paragraph>
              保存智能体后，您可以在这里测试智能体的表现。调整上面的参数以优化智能体效果。
            </Paragraph>
            <Button 
              type="primary" 
              icon={<CodeOutlined />} 
              disabled={!isEditing}
              onClick={() => navigate(`/agents/test/${id}`)}
            >
              测试智能体
            </Button>
          </Card>
        </div>
      </Form>
      
      <div style={{ marginTop: 24, display: 'flex', justifyContent: 'space-between' }}>
        <Button 
          disabled={currentStep === 0} 
          onClick={handlePrev}
          icon={<ArrowLeftOutlined />}
        >
          上一步
        </Button>
        <Button 
          type="primary" 
          disabled={currentStep === steps.length - 1} 
          onClick={handleNext}
          icon={<ArrowRightOutlined />}
        >
          下一步
        </Button>
      </div>
    </div>
  );
};

export default ChatAgentBuilder; 