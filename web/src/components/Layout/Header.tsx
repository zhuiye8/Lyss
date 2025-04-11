import React from 'react';
import { Layout, Menu, Avatar, Dropdown, Space, Button } from 'antd';
import { UserOutlined, BellOutlined, SettingOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Header } = Layout;

const AppHeader: React.FC = () => {
  const navigate = useNavigate();

  const userMenu = (
    <Menu
      items={[
        {
          key: 'profile',
          label: '个人资料',
          onClick: () => navigate('/profile'),
        },
        {
          key: 'settings',
          label: '系统设置',
          onClick: () => navigate('/settings'),
        },
        {
          type: 'divider',
        },
        {
          key: 'logout',
          label: '退出登录',
          danger: true,
          onClick: () => {
            // 处理登出逻辑
            console.log('用户登出');
          },
        },
      ]}
    />
  );

  const notificationMenu = (
    <Menu
      items={[
        {
          key: 'notification1',
          label: '新消息通知 1',
        },
        {
          key: 'notification2',
          label: '新消息通知 2',
        },
        {
          key: 'notification3',
          label: '查看全部通知',
          onClick: () => navigate('/notifications'),
        },
      ]}
    />
  );

  return (
    <Header style={{ background: '#fff', padding: '0 20px', display: 'flex', justifyContent: 'flex-end', alignItems: 'center' }}>
      <div style={{ flex: 1 }}>
        <Button type="primary" onClick={() => navigate('/agents/create')}>
          创建智能体
        </Button>
      </div>
      <Space size="large">
        <Dropdown overlay={notificationMenu} placement="bottomRight" arrow>
          <Button type="text" icon={<BellOutlined />} />
        </Dropdown>
        <Dropdown overlay={userMenu} placement="bottomRight" arrow>
          <Space>
            <Avatar icon={<UserOutlined />} />
            <span>管理员</span>
          </Space>
        </Dropdown>
      </Space>
    </Header>
  );
};

export default AppHeader; 