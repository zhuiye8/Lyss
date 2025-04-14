import React from 'react';
import { Card, List, Avatar, Tag, Spin } from 'antd';
import { 
  RobotOutlined, 
  CommentOutlined, 
  DatabaseOutlined,
  UserOutlined 
} from '@ant-design/icons';
import { IRecentActivity } from '../../types/dashboard';

interface RecentActivitiesProps {
  activities: IRecentActivity[];
  loading: boolean;
}

const RecentActivities: React.FC<RecentActivitiesProps> = ({ activities, loading }) => {
  // 根据活动类型返回图标
  const getIcon = (type: string) => {
    switch (type) {
      case 'agent_created':
      case 'agent_updated':
        return <RobotOutlined style={{ color: '#1890ff' }} />;
      case 'conversation':
        return <CommentOutlined style={{ color: '#52c41a' }} />;
      case 'knowledge_base':
        return <DatabaseOutlined style={{ color: '#722ed1' }} />;
      default:
        return <UserOutlined style={{ color: '#fa8c16' }} />;
    }
  };

  // 根据活动类型返回颜色
  const getTagColor = (type: string) => {
    switch (type) {
      case 'agent_created':
        return 'blue';
      case 'agent_updated':
        return 'cyan';
      case 'conversation':
        return 'green';
      case 'knowledge_base':
        return 'purple';
      default:
        return 'orange';
    }
  };

  // 根据活动类型返回标签文字
  const getTagText = (type: string) => {
    switch (type) {
      case 'agent_created':
        return '创建智能体';
      case 'agent_updated':
        return '更新智能体';
      case 'conversation':
        return '对话';
      case 'knowledge_base':
        return '知识库';
      default:
        return '活动';
    }
  };

  return (
    <Card title="最近活动">
      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px 0' }}>
          <Spin />
        </div>
      ) : (
      <List
        itemLayout="horizontal"
        dataSource={activities}
          renderItem={item => (
          <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon={getIcon(item.type)} />}
                title={
                  <div style={{ display: 'flex', alignItems: 'center' }}>
                    <span>{item.content}</span>
                    <Tag color={getTagColor(item.type)} style={{ marginLeft: 8 }}>
                      {getTagText(item.type)}
                    </Tag>
                  </div>
                }
                description={item.time}
              />
          </List.Item>
        )}
      />
      )}
    </Card>
  );
};

export default RecentActivities; 