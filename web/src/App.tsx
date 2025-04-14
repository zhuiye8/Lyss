import React, { useEffect, useState } from 'react';
import { Layout } from 'antd';
import { BrowserRouter as Router, Routes, Route, Navigate, useNavigate, useLocation, Link } from 'react-router-dom';
import axios from 'axios';

// 导入自定义组件
import AppHeader from './components/Layout/Header';
import AppSidebar from './components/Layout/Sidebar';
import Dashboard from './pages/Dashboard';
import AgentBuilder from './pages/AgentBuilder';
import AgentList from './pages/AgentList';
import Settings from './pages/Settings';
// 暂时注释掉缺少的组件
// import AgentTester from './pages/AgentBuilder/AgentTester';
// import KnowledgeBase from './pages/KnowledgeBase';
// import LogsPage from './pages/Logs';
// import Settings from './pages/Settings';

// 导入XConfigProvider替代原有ConfigProvider
import XConfigProvider from './components/XConfig';

// 布局组件
const { Content } = Layout;

// 注册页面组件
const RegisterPage = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleRegister = async () => {
    setLoading(true);
    setError("");
    console.log('开始注册流程...');
    
    try {
      // 调用注册API
      const response = await axios.post('/api/v1/auth/register', {
        username,
        email,
        password,
        full_name: fullName
      });
      
      console.log('注册API响应:', response.data);
      
      // 从响应中提取token
      const { token } = response.data;
      
      if (token) {
        // 存储token到localStorage
        localStorage.setItem('token', token);
        console.log('注册成功并获取Token:', token);
        
        // 注册成功后跳转到首页
        navigate('/');
      } else {
        setError("注册成功但未收到有效令牌");
      }
    } catch (error) {
      console.error('注册请求失败:', error);
      
      if (axios.isAxiosError(error) && error.response) {
        console.error('服务器响应:', error.response.data);
        setError(error.response.data.error || "注册失败，请检查输入信息");
      } else {
        setError("网络错误，请稍后重试");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <div style={{ width: 300, padding: 20, border: '1px solid #ddd', borderRadius: 8, background: 'white' }}>
        <h2>创建新账户</h2>
        <p>请填写注册信息</p>
        
        {error && (
          <div style={{ color: 'red', marginBottom: 16, padding: 8, background: '#ffeeee', borderRadius: 4 }}>
            {error}
          </div>
        )}
        
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 8 }}>用户名</label>
          <input 
            type="text" 
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }} 
          />
        </div>
        
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 8 }}>邮箱</label>
          <input 
            type="email" 
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }} 
          />
        </div>
        
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 8 }}>姓名</label>
          <input 
            type="text" 
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
            style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }} 
          />
        </div>
        
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 8 }}>密码</label>
          <input 
            type="password" 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }} 
          />
        </div>
        
        <button 
          onClick={handleRegister} 
          disabled={loading}
          style={{ 
            width: '100%', 
            padding: '10px', 
            marginTop: 16, 
            background: loading ? '#cccccc' : '#1890ff', 
            color: 'white', 
            border: 'none', 
            borderRadius: 4,
            cursor: loading ? 'not-allowed' : 'pointer'
          }}
        >
          {loading ? '注册中...' : '创建账户'}
        </button>
        
        <div style={{ marginTop: 16, textAlign: 'center' }}>
          <Link to="/login" style={{ color: '#1890ff' }}>已有账户？点击登录</Link>
        </div>
      </div>
    </div>
  );
};

// 登录页面组件
const LoginPage = () => {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleLogin = async () => {
    setLoading(true);
    setError("");
    console.log('开始登录流程...');
    
    try {
      // 调用真实的登录API获取JWT令牌
      const response = await axios.post('/api/v1/auth/login', {
        email,
        password
      });
      
      console.log('登录API响应:', response.data);
      
      // 从响应中提取token
      const { token, user } = response.data;
      
      if (token) {
        // 存储token到localStorage
        localStorage.setItem('token', token);
        console.log('Token设置成功:', token);
        
        // 存储用户信息(可选)
        if (user) {
          localStorage.setItem('user', JSON.stringify(user));
        }
        
        // 登录成功后跳转到首页
        navigate('/');
      } else {
        setError("登录成功但未收到有效令牌");
      }
    } catch (error) {
      console.error('登录请求失败:', error);
      
      if (axios.isAxiosError(error) && error.response) {
        console.error('服务器响应:', error.response.data);
        setError(error.response.data.error || "登录失败，请检查用户名和密码");
      } else {
        setError("网络错误，请稍后重试");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <div style={{ width: 300, padding: 20, border: '1px solid #ddd', borderRadius: 8, background: 'white' }}>
        <h2>智能体平台登录</h2>
        <p>请输入您的账号信息</p>
        
        {error && (
          <div style={{ color: 'red', marginBottom: 16, padding: 8, background: '#ffeeee', borderRadius: 4 }}>
            {error}
          </div>
        )}
        
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 8 }}>邮箱</label>
          <input 
            type="email" 
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }} 
          />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label style={{ display: 'block', marginBottom: 8 }}>密码</label>
          <input 
            type="password" 
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{ width: '100%', padding: '8px', borderRadius: 4, border: '1px solid #d9d9d9' }} 
          />
        </div>
        <button 
          onClick={handleLogin} 
          disabled={loading}
          style={{ 
            width: '100%', 
            padding: '10px', 
            marginTop: 16, 
            background: loading ? '#cccccc' : '#1890ff', 
            color: 'white', 
            border: 'none', 
            borderRadius: 4,
            cursor: loading ? 'not-allowed' : 'pointer'
          }}
        >
          {loading ? '登录中...' : '登录系统'}
        </button>
        
        <div style={{ marginTop: 16, textAlign: 'center' }}>
          <Link to="/register" style={{ color: '#1890ff' }}>没有账户？点击注册</Link>
        </div>
      </div>
    </div>
  );
};

// 主应用布局
const MainLayout = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    // 检查是否已认证
    const token = localStorage.getItem('token');
    const sessionAuth = sessionStorage.getItem('isLoggedIn');
    
    console.log('当前认证状态:', { token, sessionAuth });
    
    if (!token && !sessionAuth) {
      console.log('未检测到认证信息，重定向到登录页面');
      navigate('/login', { replace: true });
    } else {
      setIsAuthenticated(true);
    }
  }, [navigate, location.pathname]);

  // 如果未认证，返回null（不渲染内容）
  if (!isAuthenticated) {
    return <div style={{ padding: 20 }}>检查认证状态...</div>;
  }

  return (
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
              <Route path="/settings" element={<Settings />} />
              {/* 暂时注释掉缺少的路由 */}
              {/* <Route path="/agents/test/:id" element={<AgentTester />} /> */}
              {/* <Route path="/knowledge-base" element={<KnowledgeBase />} /> */}
              {/* <Route path="/logs" element={<LogsPage />} /> */}
              {/* <Route path="/settings" element={<Settings />} /> */}
            </Routes>
          </div>
        </Content>
      </Layout>
    </Layout>
  );
};

const App: React.FC = () => {
  return (
    <XConfigProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/*" element={<MainLayout />} />
        </Routes>
      </Router>
    </XConfigProvider>
  );
};

export default App; 