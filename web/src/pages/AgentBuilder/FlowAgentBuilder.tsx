import React, { useState, useEffect, useCallback } from 'react';
import {
  Form,
  Input,
  Card,
  Button,
  message,
  Typography,
  Space,
  Tabs,
  Divider,
  Spin,
  Modal,
  Select,
  Tag
} from 'antd';
import {
  SaveOutlined,
  RocketOutlined,
  ArrowLeftOutlined,
  DeleteOutlined,
  PlusOutlined,
  EditOutlined
} from '@ant-design/icons';
import { useNavigate, useParams } from 'react-router-dom';
import ReactFlow, { 
  Background, 
  Controls, 
  MiniMap, 
  NodeTypes, 
  Edge, 
  Node,
  addEdge,
  useNodesState,
  useEdgesState,
  MarkerType,
  NodeChange,
  EdgeChange,
  Connection
} from 'reactflow';
import 'reactflow/dist/style.css';

import ModelSelector from '../../components/AgentBuilder/ModelSelector';
import KnowledgeBaseSelector from '../../components/AgentBuilder/KnowledgeBaseSelector';
import ToolSelector from '../../components/AgentBuilder/ToolSelector';
import { 
  getAgent, 
  createAgent, 
  updateAgent, 
  publishAgent,
  getTools,
  getKnowledgeBases 
} from '../../services/agentService';
import { 
  IAgentFormState, 
  IFlowAgentConfig, 
  IAgentNode, 
  IAgentEdge,
  ITool,
  IKnowledgeBase
} from '../../types/agent';

// 自定义节点组件
import StartNode from '../../components/AgentBuilder/FlowNodes/StartNode';
import MessageNode from '../../components/AgentBuilder/FlowNodes/MessageNode';
import ConditionNode from '../../components/AgentBuilder/FlowNodes/ConditionNode';
import ToolNode from '../../components/AgentBuilder/FlowNodes/ToolNode';
import KnowledgeNode from '../../components/AgentBuilder/FlowNodes/KnowledgeNode';
import EndNode from '../../components/AgentBuilder/FlowNodes/EndNode';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;
const { Option } = Select;
const { TabPane } = Tabs;

// 注册自定义节点类型
const nodeTypes: NodeTypes = {
  start: StartNode,
  message: MessageNode,
  condition: ConditionNode,
  tool: ToolNode,
  knowledge: KnowledgeNode,
  end: EndNode
};

const FlowAgentBuilder: React.FC = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const isEditing = !!id;

  // 状态
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [publishing, setPublishing] = useState(false);
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [selectedNode, setSelectedNode] = useState<Node<any> | null>(null);
  const [selectedEdge, setSelectedEdge] = useState<Edge<any> | null>(null);
  const [nodeEditModalVisible, setNodeEditModalVisible] = useState(false);
  const [edgeEditModalVisible, setEdgeEditModalVisible] = useState(false);
  const [tools, setTools] = useState<ITool[]>([]);
  const [knowledgeBases, setKnowledgeBases] = useState<IKnowledgeBase[]>([]);
  const [formValues, setFormValues] = useState<IAgentFormState>({
    name: '',
    description: '',
    type: 'flow',
    visibility: 'private',
    modelId: '',
    status: 'draft',
    flowConfig: {
      nodes: [],
      edges: [],
      tools: [],
      knowledgeBases: [],
      variables: []
    }
  });

  // 获取工具和知识库列表
  useEffect(() => {
    const fetchData = async () => {
      try {
        const [toolsData, kbData] = await Promise.all([
          getTools(),
          getKnowledgeBases()
        ]);
        setTools(toolsData);
        setKnowledgeBases(kbData);
      } catch (error) {
        console.error('获取数据失败:', error);
        message.error('获取工具或知识库数据失败');
      }
    };
    
    fetchData();
  }, []);

  // 获取智能体信息
  useEffect(() => {
    if (isEditing) {
      fetchAgentDetails();
    } else {
      // 如果是新建，初始化一个默认流程
      initializeDefaultFlow();
    }
  }, [id]);

  // 初始化默认流程
  const initializeDefaultFlow = () => {
    const initialNodes: Node[] = [
      {
        id: 'start-1',
        type: 'start',
        position: { x: 250, y: 25 },
        data: { label: '开始' }
      },
      {
        id: 'message-1',
        type: 'message',
        position: { x: 250, y: 150 },
        data: { message: '你好，我是智能助手，有什么可以帮你的吗？' }
      },
      {
        id: 'end-1',
        type: 'end',
        position: { x: 250, y: 275 },
        data: { label: '结束' }
      }
    ];

    const initialEdges: Edge[] = [
      {
        id: 'edge-start-message',
        source: 'start-1',
        target: 'message-1',
        type: 'smoothstep',
        markerEnd: { type: MarkerType.ArrowClosed }
      },
      {
        id: 'edge-message-end',
        source: 'message-1',
        target: 'end-1',
        type: 'smoothstep',
        markerEnd: { type: MarkerType.ArrowClosed }
      }
    ];

    setNodes(initialNodes);
    setEdges(initialEdges);
  };

  // 获取智能体详情
  const fetchAgentDetails = async () => {
    try {
      setLoading(true);
      const data = await getAgent(id!);
      
      if (data.type !== 'flow') {
        message.error('当前智能体不是流程式应用');
        navigate('/agents');
        return;
      }
      
      setFormValues(data);
      form.setFieldsValue({
        ...data,
        tools: data.flowConfig?.tools || [],
        knowledgeBases: data.flowConfig?.knowledgeBases || []
      });
      
      // 转换节点和边数据
      if (data.flowConfig?.nodes && data.flowConfig?.edges) {
        const flowNodes: Node[] = data.flowConfig.nodes.map((node: IAgentNode) => ({
          id: node.id,
          type: node.type,
          position: node.position,
          data: node.data
        }));
        
        const flowEdges: Edge[] = data.flowConfig.edges.map((edge: IAgentEdge) => ({
          id: edge.id,
          source: edge.source,
          target: edge.target,
          label: edge.label,
          data: { condition: edge.condition },
          type: 'smoothstep',
          markerEnd: { type: MarkerType.ArrowClosed }
        }));
        
        setNodes(flowNodes);
        setEdges(flowEdges);
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

      // 将React Flow的节点和边转换为保存格式
      const agentNodes: IAgentNode[] = nodes.map(node => ({
        id: node.id,
        type: node.type as any,
        position: node.position,
        data: node.data
      }));

      const agentEdges: IAgentEdge[] = edges.map(edge => ({
        id: edge.id,
        source: edge.source,
        target: edge.target,
        label: edge.label,
        condition: edge.data?.condition
      }));

      // 构建保存的数据结构
      const flowConfig: Partial<IFlowAgentConfig> = {
        ...formValues.flowConfig,
        nodes: agentNodes,
        edges: agentEdges,
        knowledgeBases: values.knowledgeBases || [],
        tools: values.tools || [],
        // 从节点内容中提取变量
        variables: extractVariablesFromNodes(nodes)
      };

      const agentData: IAgentFormState = {
        name: values.name,
        description: values.description,
        type: 'flow',
        visibility: values.visibility || 'private',
        modelId: values.modelId,
        status: formValues.status,
        flowConfig
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

  // 从节点内容中提取变量
  const extractVariablesFromNodes = (nodes: Node[]) => {
    const variables: string[] = [];
    
    // 遍历所有节点，查找可能包含变量的内容
    nodes.forEach(node => {
      if (node.data) {
        // 检查消息内容
        if (node.data.message) {
          const matches = node.data.message.match(/\{\{([^}]+)\}\}/g) || [];
          matches.forEach(match => {
            const variable = match.replace(/\{\{|\}\}/g, '').trim();
            if (variable && !variables.includes(variable)) {
              variables.push(variable);
            }
          });
        }
        
        // 检查条件内容
        if (node.data.condition) {
          const matches = node.data.condition.match(/\{\{([^}]+)\}\}/g) || [];
          matches.forEach(match => {
            const variable = match.replace(/\{\{|\}\}/g, '').trim();
            if (variable && !variables.includes(variable)) {
              variables.push(variable);
            }
          });
        }
      }
    });
    
    return variables;
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

  // 连接处理
  const onConnect = useCallback((params: Connection) => {
    const newEdge = {
      ...params,
      id: `edge-${params.source}-${params.target}-${Date.now()}`,
      type: 'smoothstep',
      markerEnd: { type: MarkerType.ArrowClosed }
    };
    setEdges(edges => addEdge(newEdge, edges));
  }, []);

  // 节点点击处理
  const onNodeClick = (event: React.MouseEvent, node: Node) => {
    setSelectedNode(node);
    setNodeEditModalVisible(true);
  };

  // 边点击处理
  const onEdgeClick = (event: React.MouseEvent, edge: Edge) => {
    setSelectedEdge(edge);
    setEdgeEditModalVisible(true);
  };

  // 更新节点数据
  const updateNodeData = (nodeId: string, newData: any) => {
    setNodes(prevNodes => 
      prevNodes.map(node => 
        node.id === nodeId 
          ? { ...node, data: { ...node.data, ...newData } }
          : node
      )
    );
  };

  // 更新边数据
  const updateEdgeData = (edgeId: string, newData: any) => {
    setEdges(prevEdges => 
      prevEdges.map(edge => 
        edge.id === edgeId 
          ? { ...edge, ...newData }
          : edge
      )
    );
  };

  // 添加新节点
  const addNewNode = (type: string) => {
    const newNodeId = `${type}-${Date.now()}`;
    const newNode: Node = {
      id: newNodeId,
      type,
      position: { x: 250, y: 200 },
      data: {}
    };

    // 根据节点类型设置初始数据
    switch (type) {
      case 'message':
        newNode.data.message = '';
        break;
      case 'condition':
        newNode.data.condition = '';
        break;
      case 'tool':
        newNode.data.toolId = '';
        newNode.data.toolParams = {};
        break;
      case 'knowledge':
        newNode.data.knowledgeBaseId = '';
        break;
      case 'start':
        newNode.data.label = '开始';
        break;
      case 'end':
        newNode.data.label = '结束';
        break;
    }

    setNodes(prevNodes => [...prevNodes, newNode]);
    
    // 选中新节点进行编辑
    setSelectedNode(newNode);
    setNodeEditModalVisible(true);
  };

  // 删除节点
  const deleteNode = (nodeId: string) => {
    // 删除与节点相关的边
    setEdges(prevEdges => 
      prevEdges.filter(edge => edge.source !== nodeId && edge.target !== nodeId)
    );
    
    // 删除节点
    setNodes(prevNodes => prevNodes.filter(node => node.id !== nodeId));
    
    setNodeEditModalVisible(false);
  };

  // 删除边
  const deleteEdge = (edgeId: string) => {
    setEdges(prevEdges => prevEdges.filter(edge => edge.id !== edgeId));
    setEdgeEditModalVisible(false);
  };

  // 生成工具选项
  const generateToolOptions = () => {
    return tools.map(tool => (
      <Option key={tool.id} value={tool.id}>
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <span>{tool.name}</span>
          <Tag color="blue" style={{ marginLeft: 8 }}>{tool.type}</Tag>
        </div>
        <div style={{ fontSize: 12, color: '#666' }}>{tool.description}</div>
      </Option>
    ));
  };

  // 生成知识库选项
  const generateKnowledgeBaseOptions = () => {
    return knowledgeBases.map(kb => (
      <Option key={kb.id} value={kb.id}>
        <div>{kb.name}</div>
        <div style={{ fontSize: 12, color: '#666' }}>{kb.description}</div>
      </Option>
    ));
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

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
            {isEditing ? '编辑流程式智能体' : '创建流程式智能体'}
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

      <Form
        form={form}
        layout="vertical"
        initialValues={{
          ...formValues,
          tools: formValues.flowConfig?.tools || [],
          knowledgeBases: formValues.flowConfig?.knowledgeBases || []
        }}
      >
        <div style={{ display: 'flex', gap: '16px' }}>
          <div style={{ flex: 1 }}>
            <Card title="基础信息" style={{ marginBottom: 16 }}>
              <Form.Item
                name="name"
                label="智能体名称"
                rules={[{ required: true, message: '请输入智能体名称' }]}
              >
                <Input placeholder="例如: 订单处理流程、客户服务流程" />
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
                name="modelId"
                label="语言模型"
                rules={[{ required: true, message: '请选择语言模型' }]}
              >
                <ModelSelector />
              </Form.Item>
            </Card>
            
            <Card title="资源配置" style={{ marginBottom: 16 }}>
              <Form.Item
                name="knowledgeBases"
                label="知识库"
              >
                <KnowledgeBaseSelector />
              </Form.Item>
              
              <Form.Item
                name="tools"
                label="工具"
              >
                <ToolSelector />
              </Form.Item>
            </Card>
            
            <Card title="节点面板">
              <Paragraph>点击添加节点到流程图</Paragraph>
              <Space wrap>
                <Button 
                  type="dashed" 
                  icon={<PlusOutlined />} 
                  onClick={() => addNewNode('message')}
                >
                  消息节点
                </Button>
                <Button 
                  type="dashed" 
                  icon={<PlusOutlined />} 
                  onClick={() => addNewNode('condition')}
                >
                  条件节点
                </Button>
                <Button 
                  type="dashed" 
                  icon={<PlusOutlined />} 
                  onClick={() => addNewNode('tool')}
                >
                  工具节点
                </Button>
                <Button 
                  type="dashed" 
                  icon={<PlusOutlined />} 
                  onClick={() => addNewNode('knowledge')}
                >
                  知识库节点
                </Button>
                <Button 
                  type="dashed" 
                  icon={<PlusOutlined />} 
                  onClick={() => addNewNode('end')}
                >
                  结束节点
                </Button>
              </Space>
              
              <Divider />
              
              <Paragraph>
                <Text strong>使用说明：</Text>
              </Paragraph>
              <ul>
                <li>添加节点到画布</li>
                <li>点击节点进行编辑</li>
                <li>拖拽连线建立节点关系</li>
                <li>点击连线可编辑连线属性</li>
              </ul>
            </Card>
          </div>
          
          <div style={{ flex: 1.5 }}>
            <Card title="流程设计" style={{ height: '700px' }}>
              <div style={{ height: '650px' }}>
                <ReactFlow
                  nodes={nodes}
                  edges={edges}
                  onNodesChange={onNodesChange}
                  onEdgesChange={onEdgesChange}
                  onConnect={onConnect}
                  onNodeClick={onNodeClick}
                  onEdgeClick={onEdgeClick}
                  nodeTypes={nodeTypes}
                  fitView
                >
                  <Controls />
                  <MiniMap />
                  <Background gap={12} size={1} />
                </ReactFlow>
              </div>
            </Card>
          </div>
        </div>
      </Form>

      {/* 节点编辑对话框 */}
      <Modal
        title={`编辑${selectedNode?.type === 'message' ? '消息' : 
                    selectedNode?.type === 'condition' ? '条件' : 
                    selectedNode?.type === 'tool' ? '工具' : 
                    selectedNode?.type === 'knowledge' ? '知识库' : ''}节点`}
        open={nodeEditModalVisible}
        onCancel={() => setNodeEditModalVisible(false)}
        footer={[
          <Button 
            key="delete" 
            danger 
            icon={<DeleteOutlined />} 
            onClick={() => selectedNode && deleteNode(selectedNode.id)}
          >
            删除节点
          </Button>,
          <Button key="cancel" onClick={() => setNodeEditModalVisible(false)}>
            取消
          </Button>,
          <Button 
            key="submit" 
            type="primary" 
            onClick={() => {
              setNodeEditModalVisible(false);
            }}
          >
            确定
          </Button>
        ]}
        width={600}
      >
        {selectedNode?.type === 'message' && (
          <Form layout="vertical">
            <Form.Item
              label="消息内容"
              rules={[{ required: true, message: '请输入消息内容' }]}
            >
              <TextArea 
                rows={5} 
                placeholder="输入消息内容，可以使用{{变量}}引用变量"
                value={selectedNode.data?.message}
                onChange={(e) => updateNodeData(selectedNode.id, { message: e.target.value })}
              />
            </Form.Item>
          </Form>
        )}
        
        {selectedNode?.type === 'condition' && (
          <Form layout="vertical">
            <Form.Item
              label="条件表达式"
              rules={[{ required: true, message: '请输入条件表达式' }]}
              extra="条件表达式例如: {{user_age}} > 18，支持基本的比较和逻辑运算"
            >
              <TextArea 
                rows={3} 
                placeholder="输入条件表达式，可以使用{{变量}}引用变量"
                value={selectedNode.data?.condition}
                onChange={(e) => updateNodeData(selectedNode.id, { condition: e.target.value })}
              />
            </Form.Item>
          </Form>
        )}
        
        {selectedNode?.type === 'tool' && (
          <Form layout="vertical">
            <Form.Item
              label="选择工具"
              rules={[{ required: true, message: '请选择工具' }]}
            >
              <Select
                placeholder="选择要调用的工具"
                value={selectedNode.data?.toolId}
                onChange={(value) => updateNodeData(selectedNode.id, { toolId: value })}
                style={{ width: '100%' }}
              >
                {generateToolOptions()}
              </Select>
            </Form.Item>
            
            {selectedNode.data?.toolId && (
              <Form.Item
                label="工具参数"
                extra="参数格式为JSON，可以使用{{变量}}引用变量"
              >
                <TextArea 
                  rows={5} 
                  placeholder='{
  "param1": "value1",
  "param2": {{variable}}
}'
                  value={JSON.stringify(selectedNode.data?.toolParams || {}, null, 2)}
                  onChange={(e) => {
                    try {
                      const params = JSON.parse(e.target.value);
                      updateNodeData(selectedNode.id, { toolParams: params });
                    } catch (error) {
                      // 如果JSON解析失败，仍然更新值，但在保存时会验证
                      console.error('JSON解析失败:', error);
                    }
                  }}
                />
              </Form.Item>
            )}
          </Form>
        )}
        
        {selectedNode?.type === 'knowledge' && (
          <Form layout="vertical">
            <Form.Item
              label="选择知识库"
              rules={[{ required: true, message: '请选择知识库' }]}
            >
              <Select
                placeholder="选择要访问的知识库"
                value={selectedNode.data?.knowledgeBaseId}
                onChange={(value) => updateNodeData(selectedNode.id, { knowledgeBaseId: value })}
                style={{ width: '100%' }}
              >
                {generateKnowledgeBaseOptions()}
              </Select>
            </Form.Item>
          </Form>
        )}
      </Modal>

      {/* 边编辑对话框 */}
      <Modal
        title="编辑连线"
        open={edgeEditModalVisible}
        onCancel={() => setEdgeEditModalVisible(false)}
        footer={[
          <Button 
            key="delete" 
            danger 
            icon={<DeleteOutlined />} 
            onClick={() => selectedEdge && deleteEdge(selectedEdge.id)}
          >
            删除连线
          </Button>,
          <Button key="cancel" onClick={() => setEdgeEditModalVisible(false)}>
            取消
          </Button>,
          <Button 
            key="submit" 
            type="primary" 
            onClick={() => {
              setEdgeEditModalVisible(false);
            }}
          >
            确定
          </Button>
        ]}
      >
        <Form layout="vertical">
          <Form.Item
            label="连线标签"
          >
            <Input 
              placeholder="可选，为连线添加标签"
              value={selectedEdge?.label}
              onChange={(e) => selectedEdge && updateEdgeData(selectedEdge.id, { label: e.target.value })}
            />
          </Form.Item>
          
          <Form.Item
            label="条件表达式"
            extra="可选，来源为条件节点时的路径条件"
          >
            <TextArea 
              rows={3} 
              placeholder="例如: 'true'、'false' 或自定义表达式"
              value={selectedEdge?.data?.condition}
              onChange={(e) => selectedEdge && updateEdgeData(selectedEdge.id, { data: { ...selectedEdge.data, condition: e.target.value } })}
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default FlowAgentBuilder; 