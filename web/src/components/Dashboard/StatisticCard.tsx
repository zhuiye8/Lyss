import React from 'react';
import { Card, Statistic } from 'antd';
import { 
  RobotOutlined, 
  CommentOutlined, 
  UserOutlined, 
  ApiOutlined 
} from '@ant-design/icons';

interface StatisticCardProps {
  title: string;
  value: number | string;
  type: 'agents' | 'conversations' | 'users' | 'tokens';
  loading?: boolean;
}

const StatisticCard: React.FC<StatisticCardProps> = ({
  title,
  value,
  type,
  loading = false
}) => {
  // 根据类型选择图标
  const getIcon = () => {
    switch (type) {
      case 'agents':
        return <RobotOutlined />;
      case 'conversations':
        return <CommentOutlined />;
      case 'users':
        return <UserOutlined />;
      case 'tokens':
        return <ApiOutlined />;
      default:
        return <RobotOutlined />;
    }
  };

  // 根据类型设置颜色
  const getColor = () => {
    switch (type) {
      case 'agents':
        return '#1890ff';
      case 'conversations':
        return '#52c41a';
      case 'users':
        return '#722ed1';
      case 'tokens':
        return '#fa8c16';
      default:
        return '#1890ff';
    }
  };

  return (
    <Card loading={loading} bodyStyle={{ padding: '20px' }}>
      <Statistic
        title={title}
        value={value}
        valueStyle={{ color: getColor() }}
        prefix={getIcon()}
      />
    </Card>
  );
};

export default StatisticCard; 