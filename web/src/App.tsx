import React from 'react';
import { ConfigProvider, Layout, theme } from 'antd';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import zhCN from 'antd/lib/locale/zh_CN';

// 导入自定义组件
import AppHeader from './components/Layout/Header';
import AppSidebar from './components/Layout/Sidebar';
import Dashboard from './pages/Dashboard';
import AgentBuilder from './pages/AgentBuilder';
import AgentList from './pages/AgentList';
import KnowledgeBase from './pages/KnowledgeBase';
import Settings from './pages/Settings';

// 布局组件
const { Content } = Layout;

const App: React.FC = () => {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        algorithm: theme.defaultAlgorithm,
        token: {
          colorPrimary: '#00b96b',
          borderRadius: 4,
        },
      }}
    >
      <Router>
        <Layout style={{ minHeight: '100vh' }}>
          <AppSidebar />
          <Layout>
            <AppHeader />
            <Content style={{ margin: '16px' }}>
              <div style={{ padding: 24, background: '#fff', minHeight: 360 }}>
                <Routes>
                  <Route path="/" element={<Dashboard />} />
                  <Route path="/agents" element={<AgentList />} />
                  <Route path="/agents/create" element={<AgentBuilder />} />
                  <Route path="/agents/edit/:id" element={<AgentBuilder />} />
                  <Route path="/knowledge-base" element={<KnowledgeBase />} />
                  <Route path="/settings" element={<Settings />} />
                </Routes>
              </div>
            </Content>
          </Layout>
        </Layout>
      </Router>
    </ConfigProvider>
  );
};

export default App; 