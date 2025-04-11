import React, { useState } from 'react';
import { Tabs, Card, Form, Input, Button, message, Typography, Switch } from 'antd';
import {
  UserOutlined,
  KeyOutlined,
  ApiOutlined,
  DatabaseOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import useAuthStore from '../store/useAuthStore';

const { Title, Paragraph } = Typography;
const { TabPane } = Tabs;

const Settings: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const { user } = useAuthStore();
  
  const handleProfileUpdate = (values: any) => {
    setLoading(true);
    // 模拟API调用
    setTimeout(() => {
      setLoading(false);
      message.success('个人资料已更新');
    }, 1000);
  };
  
  const handlePasswordChange = (values: any) => {
    setLoading(true);
    // 模拟API调用
    setTimeout(() => {
      setLoading(false);
      message.success('密码已修改');
    }, 1000);
  };

  return (
    <div>
      <Title level={2}>系统设置</Title>
      <Paragraph>您可以在这里管理您的个人资料和系统配置。</Paragraph>
      
      <Card>
        <Tabs defaultActiveKey="profile">
          <TabPane 
            tab={
              <span>
                <UserOutlined />
                个人资料
              </span>
            }
            key="profile"
          >
            <Form
              layout="vertical"
              initialValues={{
                username: user?.username || '',
                email: user?.email || '',
              }}
              onFinish={handleProfileUpdate}
            >
              <Form.Item
                label="用户名"
                name="username"
                rules={[{ required: true, message: '请输入用户名' }]}
              >
                <Input prefix={<UserOutlined />} />
              </Form.Item>
              
              <Form.Item
                label="电子邮箱"
                name="email"
                rules={[
                  { required: true, message: '请输入电子邮箱' },
                  { type: 'email', message: '请输入有效的电子邮箱' }
                ]}
              >
                <Input />
              </Form.Item>
              
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  保存修改
                </Button>
              </Form.Item>
            </Form>
          </TabPane>
          
          <TabPane
            tab={
              <span>
                <KeyOutlined />
                修改密码
              </span>
            }
            key="password"
          >
            <Form
              layout="vertical"
              onFinish={handlePasswordChange}
            >
              <Form.Item
                label="当前密码"
                name="currentPassword"
                rules={[{ required: true, message: '请输入当前密码' }]}
              >
                <Input.Password />
              </Form.Item>
              
              <Form.Item
                label="新密码"
                name="newPassword"
                rules={[
                  { required: true, message: '请输入新密码' },
                  { min: 8, message: '密码长度不能少于8个字符' }
                ]}
              >
                <Input.Password />
              </Form.Item>
              
              <Form.Item
                label="确认新密码"
                name="confirmPassword"
                dependencies={['newPassword']}
                rules={[
                  { required: true, message: '请确认新密码' },
                  ({ getFieldValue }) => ({
                    validator(_, value) {
                      if (!value || getFieldValue('newPassword') === value) {
                        return Promise.resolve();
                      }
                      return Promise.reject(new Error('两次输入的密码不一致'));
                    },
                  }),
                ]}
              >
                <Input.Password />
              </Form.Item>
              
              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  修改密码
                </Button>
              </Form.Item>
            </Form>
          </TabPane>
          
          <TabPane
            tab={
              <span>
                <ApiOutlined />
                API配置
              </span>
            }
            key="api"
          >
            <Form layout="vertical">
              <Form.Item label="API密钥">
                <Input.Password defaultValue="sk-xxxxxxxxxxxxxxxxxxxxxxxx" />
                <Button type="link" style={{ padding: 0, marginTop: 8 }}>
                  生成新密钥
                </Button>
              </Form.Item>
              
              <Form.Item label="请求速率限制" name="rateLimit">
                <Input addonAfter="请求/分钟" defaultValue="60" />
              </Form.Item>
              
              <Form.Item>
                <Button type="primary">保存配置</Button>
              </Form.Item>
            </Form>
          </TabPane>
          
          <TabPane
            tab={
              <span>
                <SettingOutlined />
                系统配置
              </span>
            }
            key="system"
          >
            <Form layout="vertical">
              <Form.Item label="开启调试模式" name="debugMode" valuePropName="checked">
                <Switch />
              </Form.Item>
              
              <Form.Item label="自动更新" name="autoUpdate" valuePropName="checked">
                <Switch defaultChecked />
              </Form.Item>
              
              <Form.Item label="日志级别" name="logLevel">
                <select style={{ width: 200, height: 32, borderRadius: 2 }}>
                  <option value="info">信息</option>
                  <option value="warning">警告</option>
                  <option value="error">错误</option>
                  <option value="debug">调试</option>
                </select>
              </Form.Item>
              
              <Form.Item>
                <Button type="primary">保存配置</Button>
              </Form.Item>
            </Form>
          </TabPane>
          
          <TabPane
            tab={
              <span>
                <DatabaseOutlined />
                存储设置
              </span>
            }
            key="storage"
          >
            <Paragraph>
              系统当前存储使用情况:
            </Paragraph>
            
            <Card style={{ marginBottom: 16 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                <span>文档存储</span>
                <span>2.3 GB / 10 GB</span>
              </div>
              <div style={{ height: 8, background: '#f0f0f0', borderRadius: 4 }}>
                <div style={{ height: '100%', width: '23%', background: '#1890ff', borderRadius: 4 }}></div>
              </div>
            </Card>
            
            <Card style={{ marginBottom: 16 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                <span>向量存储</span>
                <span>1.7 GB / 5 GB</span>
              </div>
              <div style={{ height: 8, background: '#f0f0f0', borderRadius: 4 }}>
                <div style={{ height: '100%', width: '34%', background: '#52c41a', borderRadius: 4 }}></div>
              </div>
            </Card>
            
            <Button type="primary">升级存储空间</Button>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default Settings; 