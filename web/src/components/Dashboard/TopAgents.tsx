import React from 'react';
import { Card, Table, Progress, Button } from 'antd';
import { RobotOutlined, ArrowRightOutlined } from '@ant-design/icons';
import { ITopAgent } from '../../types/dashboard';
import { useNavigate } from 'react-router-dom';

interface TopAgentsProps {
  agents: ITopAgent[];
  loading?: boolean;
}

const TopAgents: React.FC<TopAgentsProps> = ({
  agents,
  loading = false,
}) => {
  const navigate = useNavigate();
  
  // 找出使用量最大值，用于计算进度条百分比
  const maxUsage = Math.max(...agents.map(a => a.usageCount), 1);
  
  const columns = [
    {
      title: '智能体',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <div style={{ display: 'flex', alignItems: 'center' }}>
          <RobotOutlined style={{ marginRight: 8, color: '#1890ff' }} />
          <span>{text}</span>
        </div>
      ),
    },
    {
      title: '使用次数',
      dataIndex: 'usageCount',
      key: 'usageCount',
      render: (count: number) => (
        <div style={{ width: 150 }}>
          <Progress 
            percent={Math.round((count / maxUsage) * 100)} 
            format={() => count}
            strokeColor={{ from: '#108ee9', to: '#87d068' }}
          />
        </div>
      ),
    },
    {
      title: 'Token 使用量',
      dataIndex: 'tokenUsage',
      key: 'tokenUsage',
      render: (tokens: number) => (
        <span>{tokens.toLocaleString()}</span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (text: string, record: ITopAgent) => (
        <Button 
          type="link" 
          icon={<ArrowRightOutlined />}
          onClick={() => navigate(`/agents/edit/${record.id}`)}
        >
          查看
        </Button>
      ),
    },
  ];

  return (
    <Card 
      title="热门智能体" 
      extra={<Button type="link" onClick={() => navigate('/agents')}>查看全部</Button>}
    >
      <Table 
        dataSource={agents} 
        columns={columns} 
        pagination={false}
        rowKey="id"
        loading={loading}
        size="small"
      />
    </Card>
  );
};

export default TopAgents; 