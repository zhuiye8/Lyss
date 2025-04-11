import React from 'react';
import { List, Card, Tag, Avatar, Skeleton } from 'antd';
import { 
  RobotOutlined, 
  EditOutlined, 
  CommentOutlined, 
  DatabaseOutlined,
  UserOutlined 
} from '@ant-design/icons';
import { IRecentActivity } from '../../types/dashboard';

interface RecentActivitiesProps {
  activities: IRecentActivity[];
  loading?: boolean;
}

const RecentActivities: React.FC<RecentActivitiesProps> = ({
  activities,
  loading = false,
}) => {
  // 根据活动类型获取图标
  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'agent_created':
        return <RobotOutlined style={{ color: '#1890ff' }} />;
      case 'agent_updated':
        return <EditOutlined style={{ color: '#faad14' }} />;
      case 'conversation':
        return <CommentOutlined style={{ color: '#52c41a' }} />;
      case 'knowledge_updated':
        return <DatabaseOutlined style={{ color: '#722ed1' }} />;
      default:
        return <RobotOutlined />;
    }
  };

  // 根据活动类型获取标签
  const getActivityTag = (type: string) => {
    switch (type) {
      case 'agent_created':
        return <Tag color="blue">创建智能体</Tag>;
      case 'agent_updated':
        return <Tag color="orange">更新智能体</Tag>;
      case 'conversation':
        return <Tag color="green">新对话</Tag>;
      case 'knowledge_updated':
        return <Tag color="purple">更新知识库</Tag>;
      default:
        return <Tag>活动</Tag>;
    }
  };

  return (
    <Card title="最近活动">
      <List
        itemLayout="horizontal"
        dataSource={activities}
        loading={loading}
        renderItem={(item) => (
          <List.Item>
            <Skeleton avatar title={false} loading={loading} active>
              <List.Item.Meta
                avatar={
                  <Avatar icon={<UserOutlined />} />
                }
                title={
                  <div style={{ display: 'flex', alignItems: 'center' }}>
                    <span style={{ marginRight: 8 }}>{item.user}</span>
                    {getActivityTag(item.type)}
                  </div>
                }
                description={
                  <div>
                    <div style={{ display: 'flex', alignItems: 'center', marginBottom: 4 }}>
                      {getActivityIcon(item.type)}
                      <span style={{ marginLeft: 8 }}>
                        {item.agentName && `${item.agentName}: `}{item.details}
                      </span>
                    </div>
                    <div style={{ color: '#8c8c8c', fontSize: 12 }}>
                      {new Date(item.timestamp).toLocaleString()}
                    </div>
                  </div>
                }
              />
            </Skeleton>
          </List.Item>
        )}
      />
    </Card>
  );
};

export default RecentActivities; 