import React from 'react';
import { Card, Statistic } from 'antd';
import { 
  RobotOutlined,
  CommentOutlined,
  UserOutlined,
  ThunderboltOutlined
} from '@ant-design/icons';

interface StatisticCardProps {
  title: string;
  value: number | string;
  prefix?: string;
  suffix?: string;
  type: 'agents' | 'conversations' | 'users' | 'tokens';
  loading?: boolean;
}

const StatisticCard: React.FC<StatisticCardProps> = ({
  title,
  value,
  prefix,
  suffix,
  type,
  loading = false
}) => {
  // 根据类型选择图标
  const getIcon = () => {
    switch (type) {
      case 'agents':
        return <RobotOutlined style={{ fontSize: 20, color: '#1890ff' }} />;
      case 'conversations':
        return <CommentOutlined style={{ fontSize: 20, color: '#52c41a' }} />;
      case 'users':
        return <UserOutlined style={{ fontSize: 20, color: '#fa8c16' }} />;
      case 'tokens':
        return <ThunderboltOutlined style={{ fontSize: 20, color: '#722ed1' }} />;
      default:
        return null;
    }
  };

  // 根据类型选择卡片颜色
  const getCardStyle = () => {
    switch (type) {
      case 'agents':
        return { borderTop: '3px solid #1890ff' };
      case 'conversations':
        return { borderTop: '3px solid #52c41a' };
      case 'users':
        return { borderTop: '3px solid #fa8c16' };
      case 'tokens':
        return { borderTop: '3px solid #722ed1' };
      default:
        return {};
    }
  };

  return (
    <Card 
      style={{ ...getCardStyle(), height: '100%' }} 
      loading={loading}
    >
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: 8 }}>
        {getIcon()}
        <span style={{ marginLeft: 8, color: '#8c8c8c' }}>{title}</span>
      </div>
      <Statistic 
        value={value} 
        prefix={prefix} 
        suffix={suffix}
        valueStyle={{ color: '#000000', fontSize: 24, fontWeight: 'bold' }}
      />
    </Card>
  );
};

export default StatisticCard; 