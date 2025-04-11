import React, { useState } from 'react';
import { Layout, Menu } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  DashboardOutlined,
  RobotOutlined,
  DatabaseOutlined,
  ApiOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
} from '@ant-design/icons';

const { Sider } = Layout;

const AppSidebar: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [collapsed, setCollapsed] = useState(false);

  // 根据当前路径获取默认选中的菜单项
  const getSelectedKey = () => {
    const path = location.pathname;
    if (path === '/') return ['dashboard'];
    if (path.startsWith('/agents')) return ['agents'];
    if (path.startsWith('/knowledge-base')) return ['knowledge'];
    if (path.startsWith('/api')) return ['api'];
    if (path.startsWith('/settings')) return ['settings'];
    return ['dashboard'];
  };

  return (
    <Sider 
      collapsible 
      collapsed={collapsed} 
      onCollapse={value => setCollapsed(value)}
      theme="light"
      width={220}
    >
      <div style={{ height: 64, padding: '16px', textAlign: 'center' }}>
        <h2 style={{ color: '#00b96b', display: collapsed ? 'none' : 'block' }}>智能体平台</h2>
        {collapsed ? <RobotOutlined style={{ fontSize: 24, color: '#00b96b' }} /> : null}
      </div>
      <Menu
        mode="inline"
        selectedKeys={getSelectedKey()}
        items={[
          {
            key: 'dashboard',
            icon: <DashboardOutlined />,
            label: '仪表盘',
            onClick: () => navigate('/'),
          },
          {
            key: 'agents',
            icon: <RobotOutlined />,
            label: '智能体',
            onClick: () => navigate('/agents'),
          },
          {
            key: 'knowledge',
            icon: <DatabaseOutlined />,
            label: '知识库',
            onClick: () => navigate('/knowledge-base'),
          },
          {
            key: 'api',
            icon: <ApiOutlined />,
            label: 'API管理',
            onClick: () => navigate('/api'),
          },
          {
            key: 'settings',
            icon: <SettingOutlined />,
            label: '系统设置',
            onClick: () => navigate('/settings'),
          },
        ]}
      />
      <div 
        style={{ 
          position: 'absolute', 
          bottom: 16, 
          width: '100%', 
          textAlign: 'center',
          cursor: 'pointer',
          color: '#999'
        }}
        onClick={() => setCollapsed(!collapsed)}
      >
        {collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
      </div>
    </Sider>
  );
};

export default AppSidebar; 