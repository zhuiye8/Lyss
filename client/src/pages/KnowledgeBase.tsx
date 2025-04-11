import React, { useState } from 'react';
import { 
  Typography, 
  Button, 
  Table, 
  Card, 
  Space, 
  Modal, 
  Form, 
  Input, 
  Upload, 
  message,
  Tabs,
  Progress,
  Tag
} from 'antd';
import { 
  PlusOutlined, 
  UploadOutlined, 
  EditOutlined, 
  DeleteOutlined,
  FolderOpenOutlined,
  FileTextOutlined,
  InfoCircleOutlined
} from '@ant-design/icons';
import type { UploadProps } from 'antd';

const { Title, Paragraph } = Typography;
const { TabPane } = Tabs;

// 模拟数据
const mockKnowledgeBases = [
  {
    id: '1',
    name: '产品手册',
    description: '包含所有产品说明书和使用指南',
    documentCount: 15,
    chunkCount: 230,
    createdAt: '2025-03-15 14:30',
  },
  {
    id: '2',
    name: 'API文档',
    description: '系统API接口文档',
    documentCount: 8,
    chunkCount: 120,
    createdAt: '2025-03-20 09:15',
  },
  {
    id: '3',
    name: '常见问题',
    description: '用户常见问题与解答',
    documentCount: 25,
    chunkCount: 310,
    createdAt: '2025-04-01 16:45',
  }
];

const mockDocuments = [
  {
    id: '1',
    name: '用户手册.pdf',
    type: 'pdf',
    size: '2.5MB',
    chunkCount: 45,
    createdAt: '2025-03-16 10:30',
  },
  {
    id: '2',
    name: '安装指南.docx',
    type: 'docx',
    size: '1.2MB',
    chunkCount: 28,
    createdAt: '2025-03-18 14:20',
  },
  {
    id: '3',
    name: '常见问题.md',
    type: 'markdown',
    size: '0.5MB',
    chunkCount: 35,
    createdAt: '2025-03-22 09:15',
  }
];

const KnowledgeBase: React.FC = () => {
  const [activeKnowledgeBase, setActiveKnowledgeBase] = useState<string | null>(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalType, setModalType] = useState<'create' | 'edit'>('create');
  const [isUploadModalVisible, setIsUploadModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [uploadForm] = Form.useForm();

  // 打开创建知识库模态框
  const showCreateModal = () => {
    setModalType('create');
    form.resetFields();
    setIsModalVisible(true);
  };

  // 打开编辑知识库模态框
  const showEditModal = (record: any) => {
    setModalType('edit');
    form.setFieldsValue({
      name: record.name,
      description: record.description,
    });
    setIsModalVisible(true);
  };

  // 处理模态框确认
  const handleModalOk = () => {
    form.validateFields().then(values => {
      console.log('Form values:', values);
      setIsModalVisible(false);
      message.success(modalType === 'create' ? '知识库创建成功' : '知识库更新成功');
    });
  };

  // 打开上传文档模态框
  const showUploadModal = () => {
    uploadForm.resetFields();
    setIsUploadModalVisible(true);
  };

  // 处理上传文档模态框确认
  const handleUploadOk = () => {
    uploadForm.validateFields().then(values => {
      console.log('Upload form values:', values);
      setIsUploadModalVisible(false);
      message.success('文档上传成功');
    });
  };

  // 上传组件配置
  const uploadProps: UploadProps = {
    name: 'file',
    action: '/api/upload',
    headers: {
      authorization: 'authorization-text',
    },
    onChange(info) {
      if (info.file.status !== 'uploading') {
        console.log(info.file, info.fileList);
      }
      if (info.file.status === 'done') {
        message.success(`${info.file.name} 上传成功`);
      } else if (info.file.status === 'error') {
        message.error(`${info.file.name} 上传失败`);
      }
    },
  };

  // 知识库表格列定义
  const kbColumns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: any) => (
        <a onClick={() => setActiveKnowledgeBase(record.id)}>{text}</a>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '文档数量',
      dataIndex: 'documentCount',
      key: 'documentCount',
    },
    {
      title: '文本块数量',
      dataIndex: 'chunkCount',
      key: 'chunkCount',
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<FolderOpenOutlined />} 
            onClick={() => setActiveKnowledgeBase(record.id)}
          />
          <Button 
            type="text" 
            icon={<EditOutlined />} 
            onClick={() => showEditModal(record)}
          />
          <Button 
            type="text" 
            icon={<DeleteOutlined />} 
            danger
            onClick={() => message.success('删除成功')}
          />
        </Space>
      ),
    },
  ];

  // 文档表格列定义
  const docColumns = [
    {
      title: '文档名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: any) => (
        <Space>
          <FileTextOutlined />
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={
          type === 'pdf' ? 'red' : 
          type === 'docx' ? 'blue' : 
          type === 'markdown' ? 'green' : 'default'
        }>
          {type.toUpperCase()}
        </Tag>
      ),
    },
    {
      title: '大小',
      dataIndex: 'size',
      key: 'size',
    },
    {
      title: '文本块数量',
      dataIndex: 'chunkCount',
      key: 'chunkCount',
    },
    {
      title: '上传时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<InfoCircleOutlined />} 
            onClick={() => message.info('查看文档详情')}
          />
          <Button 
            type="text" 
            icon={<DeleteOutlined />} 
            danger
            onClick={() => message.success('删除成功')}
          />
        </Space>
      ),
    },
  ];

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Title level={2}>知识库管理</Title>
        <Paragraph>创建和管理知识库，上传文档以增强智能体的回答能力。</Paragraph>
      </div>

      {!activeKnowledgeBase ? (
        <Card>
          <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
            <div>
              <Title level={4}>知识库列表</Title>
            </div>
            <Button type="primary" icon={<PlusOutlined />} onClick={showCreateModal}>
              创建知识库
            </Button>
          </div>
          <Table 
            columns={kbColumns} 
            dataSource={mockKnowledgeBases} 
            rowKey="id"
          />
        </Card>
      ) : (
        <div>
          <div style={{ marginBottom: 16 }}>
            <Button type="link" onClick={() => setActiveKnowledgeBase(null)}>
              &lt; 返回知识库列表
            </Button>
          </div>
          
          <Card>
            <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
              <div>
                <Title level={4}>
                  {mockKnowledgeBases.find(kb => kb.id === activeKnowledgeBase)?.name}
                </Title>
                <Paragraph>
                  {mockKnowledgeBases.find(kb => kb.id === activeKnowledgeBase)?.description}
                </Paragraph>
              </div>
              <Button type="primary" icon={<UploadOutlined />} onClick={showUploadModal}>
                上传文档
              </Button>
            </div>
            
            <Tabs defaultActiveKey="documents">
              <TabPane 
                tab={
                  <span>
                    <FileTextOutlined />
                    文档列表
                  </span>
                }
                key="documents"
              >
                <Table 
                  columns={docColumns} 
                  dataSource={mockDocuments} 
                  rowKey="id"
                />
              </TabPane>
              
              <TabPane 
                tab={
                  <span>
                    <InfoCircleOutlined />
                    知识库信息
                  </span>
                }
                key="info"
              >
                <div style={{ maxWidth: 600 }}>
                  <div style={{ marginBottom: 24 }}>
                    <Paragraph>文档数量: 15</Paragraph>
                    <Paragraph>文本块数量: 230</Paragraph>
                    <Paragraph>创建时间: 2025-03-15 14:30</Paragraph>
                    <Paragraph>最后更新: 2025-04-05 09:20</Paragraph>
                  </div>
                  
                  <div style={{ marginBottom: 16 }}>
                    <Title level={5}>存储使用情况</Title>
                    <div style={{ marginBottom: 8 }}>
                      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                        <span>存储空间</span>
                        <span>4.2 MB / 1 GB</span>
                      </div>
                      <Progress percent={1} size="small" />
                    </div>
                  </div>
                </div>
              </TabPane>
            </Tabs>
          </Card>
        </div>
      )}

      {/* 创建/编辑知识库模态框 */}
      <Modal
        title={modalType === 'create' ? '创建知识库' : '编辑知识库'}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={() => setIsModalVisible(false)}
      >
        <Form
          form={form}
          layout="vertical"
        >
          <Form.Item
            name="name"
            label="知识库名称"
            rules={[{ required: true, message: '请输入知识库名称' }]}
          >
            <Input placeholder="请输入知识库名称" />
          </Form.Item>
          
          <Form.Item
            name="description"
            label="描述"
            rules={[{ required: true, message: '请输入知识库描述' }]}
          >
            <Input.TextArea placeholder="请输入知识库描述" rows={4} />
          </Form.Item>
          
          <Form.Item
            name="embedding_model"
            label="向量模型"
            rules={[{ required: true, message: '请选择向量模型' }]}
            initialValue="default"
          >
            <select style={{ width: '100%', height: 32, borderRadius: 2 }}>
              <option value="default">默认模型</option>
              <option value="openai">OpenAI Embeddings</option>
              <option value="bge">BGE Embeddings</option>
            </select>
          </Form.Item>
        </Form>
      </Modal>

      {/* 上传文档模态框 */}
      <Modal
        title="上传文档"
        open={isUploadModalVisible}
        onOk={handleUploadOk}
        onCancel={() => setIsUploadModalVisible(false)}
      >
        <Form
          form={uploadForm}
          layout="vertical"
        >
          <Form.Item
            name="upload"
            label="选择文件"
            rules={[{ required: true, message: '请上传文件' }]}
          >
            <Upload {...uploadProps}>
              <Button icon={<UploadOutlined />}>选择文件</Button>
            </Upload>
          </Form.Item>
          
          <Paragraph>
            支持的文件格式: PDF, DOCX, TXT, MD, HTML
          </Paragraph>
          <Paragraph>
            单个文件大小限制: 10MB
          </Paragraph>
        </Form>
      </Modal>
    </div>
  );
};

export default KnowledgeBase; 