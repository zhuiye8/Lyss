import React from 'react';
import { Layout } from 'antd';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

// 导入自定义组件
import AppHeader from './components/Layout/Header';
import AppSidebar from './components/Layout/Sidebar';
import Dashboard from './pages/Dashboard';
import AgentBuilder from './pages/AgentBuilder';
import AgentList from './pages/AgentList';
import AgentTester from './pages/AgentBuilder/AgentTester';
import KnowledgeBase from './pages/KnowledgeBase';
import LogsPage from './pages/Logs';
import Settings from './pages/Settings';

// 导入XConfigProvider替代原有ConfigProvider
import XConfigProvider from './components/XConfig';

// 布局组件
const { Content } = Layout;

const App: React.FC = () => {
  return (
    <XConfigProvider>
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
                  <Route path="/agents/test/:id" element={<AgentTester />} />
                  <Route path="/knowledge-base" element={<KnowledgeBase />} />
                  <Route path="/logs" element={<LogsPage />} />
                  <Route path="/settings" element={<Settings />} />
                </Routes>
              </div>
            </Content>
          </Layout>
        </Layout>
      </Router>
    </XConfigProvider>
  );
};

export default App; 