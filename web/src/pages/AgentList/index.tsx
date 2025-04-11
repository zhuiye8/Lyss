import React, { useEffect, useState } from 'react';
import { 
  Table, 
  Button, 
  Input, 
  Space, 
  Typography, 
  Tag, 
  Popconfirm, 
  message, 
  Divider,
  Badge
} from 'antd';
import { 
  SearchOutlined, 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  PoweroffOutlined 
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getAgents, deleteAgent, updateAgentStatus } from '../../services/dashboardService';
import { IAgentData } from '../../types/dashboard';

const { Title } = Typography;

const AgentList: React.FC = () => {
  const navigate = useNavigate();
  const [agents, setAgents] = useState<IAgentData[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0
  });

  // 获取智能体列表
  const fetchAgents = async (page = 1, pageSize = 10, search = '') => {
    try {
      setLoading(true);
      const response = await getAgents(page, pageSize, search);
      setAgents(response.data);
      setPagination({
        ...pagination,
        current: page,
        total: response.total
      });
    } catch (error) {
      console.error('获取智能体列表失败:', error);
      message.error('获取智能体列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    fetchAgents();
  }, []);

  // 处理搜索
  const handleSearch = () => {
    fetchAgents(1, pagination.pageSize, searchQuery);
  };

  // 处理表格分页、排序、筛选变化
  const handleTableChange = (pagination: any) => {
    fetchAgents(pagination.current, pagination.pageSize, searchQuery);
  };

  // 删除智能体
  const handleDelete = async (id: string) => {
    try {
      await deleteAgent(id);
      message.success('智能体已删除');
      // 重新加载列表
      fetchAgents(pagination.current, pagination.pageSize, searchQuery);
    } catch (error) {
      console.error('删除智能体失败:', error);
      message.error('删除智能体失败');
    }
  };

  // 更新智能体状态
  const handleUpdateStatus = async (id: string, currentStatus: 'active' | 'inactive') => {
    const newStatus = currentStatus === 'active' ? 'inactive' : 'active';
    try {
      await updateAgentStatus(id, newStatus);
      message.success(`智能体状态已更新为${newStatus === 'active' ? '启用' : '禁用'}`);
      // 更新本地状态
      setAgents(agents.map(agent => 
        agent.id === id ? { ...agent, status: newStatus } : agent
      ));
    } catch (error) {
      console.error('更新智能体状态失败:', error);
      message.error('更新智能体状态失败');
    }
  };

  // 表格列配置
  const columns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
      render: (text: string, record: IAgentData) => (
        <a onClick={() => navigate(`/agents/edit/${record.id}`)}>{text}</a>
      ),
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        type === 'chat' ? <Tag color="blue">对话型</Tag> : <Tag color="purple">流程型</Tag>
      ),
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        if (status === 'active') {
          return <Badge status="success" text="已启用" />;
        } else if (status === 'inactive') {
          return <Badge status="default" text="已禁用" />;
        } else {
          return <Badge status="warning" text="草稿" />;
        }
      },
    },
    {
      title: '使用次数',
      dataIndex: 'usageCount',
      key: 'usageCount',
      sorter: true,
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (date: string) => new Date(date).toLocaleString(),
      sorter: true,
    },
    {
      title: '上次访问',
      dataIndex: 'lastAccessed',
      key: 'lastAccessed',
      render: (date: string) => new Date(date).toLocaleString(),
      sorter: true,
    },
    {
      title: '操作',
      key: 'action',
      render: (text: string, record: IAgentData) => (
        <Space size="small">
          <Button 
            type="text" 
            icon={<EditOutlined />} 
            onClick={() => navigate(`/agents/edit/${record.id}`)}
          />
          <Button
            type="text"
            icon={<PoweroffOutlined />}
            onClick={() => handleUpdateStatus(record.id, record.status as 'active' | 'inactive')}
            danger={record.status === 'active'}
          />
          <Popconfirm
            title="确定要删除这个智能体吗？"
            onConfirm={() => handleDelete(record.id)}
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
        <Title level={3}>智能体管理</Title>
        <Button 
          type="primary" 
          icon={<PlusOutlined />} 
          onClick={() => navigate('/agents/create')}
        >
          创建智能体
        </Button>
      </div>
      <Divider style={{ margin: '16px 0' }} />
      
      {/* 搜索栏 */}
      <div style={{ marginBottom: 16 }}>
        <Input
          placeholder="搜索智能体名称或描述"
          value={searchQuery}
          onChange={e => setSearchQuery(e.target.value)}
          onPressEnter={handleSearch}
          suffix={
            <Button type="text" icon={<SearchOutlined />} onClick={handleSearch} />
          }
          style={{ width: 300 }}
        />
      </div>
      
      {/* 表格 */}
      <Table
        columns={columns}
        dataSource={agents}
        rowKey="id"
        loading={loading}
        pagination={pagination}
        onChange={handleTableChange}
      />
    </div>
  );
};

export default AgentList; 