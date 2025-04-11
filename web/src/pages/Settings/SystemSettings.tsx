import React, { useEffect, useState } from 'react';
import { 
  Form, 
  Input, 
  Button, 
  Select, 
  Switch, 
  InputNumber, 
  Upload, 
  Space,
  Divider,
  message,
  Typography,
  Spin,
  Card,
  Collapse
} from 'antd';
import { 
  SaveOutlined, 
  UploadOutlined, 
  CloudUploadOutlined,
  InfoCircleOutlined
} from '@ant-design/icons';
import { getSystemSettings, updateSystemSettings } from '../../services/dashboardService';
import { ISystemSettings } from '../../types/dashboard';

const { Title, Text } = Typography;
const { Option } = Select;
const { Panel } = Collapse;

const SystemSettings: React.FC = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [settings, setSettings] = useState<ISystemSettings | null>(null);
  const [storageType, setStorageType] = useState<'local' | 's3'>('local');

  // 获取系统设置
  const fetchSettings = async () => {
    try {
      setLoading(true);
      const data = await getSystemSettings();
      setSettings(data);
      setStorageType(data.storageProvider);
      form.setFieldsValue(data);
    } catch (error) {
      console.error('获取系统设置失败:', error);
      message.error('获取系统设置失败');
    } finally {
      setLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    fetchSettings();
  }, []);

  // 保存设置
  const handleSaveSettings = async (values: any) => {
    try {
      setSaving(true);
      await updateSystemSettings(values);
      message.success('设置已保存');
      // 刷新数据
      fetchSettings();
    } catch (error) {
      console.error('保存设置失败:', error);
      message.error('保存设置失败');
    } finally {
      setSaving(false);
    }
  };

  // 上传前验证
  const beforeUpload = (file: File) => {
    const isImage = file.type.startsWith('image/');
    if (!isImage) {
      message.error('只能上传图片文件!');
    }
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      message.error('图片必须小于2MB!');
    }
    return isImage && isLt2M;
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  return (
    <div>
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSaveSettings}
        initialValues={{
          siteName: '智能体构建平台',
          allowRegistration: true,
          defaultLanguage: 'zh-CN',
          apiRateLimit: 100,
          storageProvider: 'local'
        }}
      >
        <Title level={4}>基本设置</Title>
        <Card style={{ marginBottom: 16 }}>
          <Form.Item
            name="siteName"
            label="平台名称"
            rules={[{ required: true, message: '请输入平台名称' }]}
          >
            <Input placeholder="例如: 智能体构建平台" />
          </Form.Item>
          
          <Form.Item
            name="logoUrl"
            label="平台Logo"
            extra="建议尺寸: 200x50px, 格式: PNG, 背景透明"
          >
            <Upload
              name="logo"
              listType="picture"
              maxCount={1}
              beforeUpload={beforeUpload}
              action="/api/v1/settings/upload-logo"
              headers={{
                Authorization: `Bearer ${localStorage.getItem('token')}`,
              }}
            >
              <Button icon={<UploadOutlined />}>上传Logo</Button>
            </Upload>
          </Form.Item>
          
          <Form.Item
            name="defaultLanguage"
            label="默认语言"
            rules={[{ required: true, message: '请选择默认语言' }]}
          >
            <Select>
              <Option value="zh-CN">简体中文</Option>
              <Option value="en-US">English</Option>
            </Select>
          </Form.Item>
          
          <Form.Item
            name="defaultModel"
            label="默认模型"
            rules={[{ required: true, message: '请选择默认模型' }]}
          >
            <Select placeholder="选择默认使用的AI模型">
              <Option value="gpt-4">GPT-4</Option>
              <Option value="claude-3">Claude 3</Option>
              <Option value="llama-3">Llama 3</Option>
            </Select>
          </Form.Item>
        </Card>
        
        <Title level={4}>安全设置</Title>
        <Card style={{ marginBottom: 16 }}>
          <Form.Item
            name="allowRegistration"
            label="允许新用户注册"
            valuePropName="checked"
          >
            <Switch />
          </Form.Item>
          
          <Form.Item
            name="apiRateLimit"
            label="API请求限制(次/分钟)"
            rules={[{ required: true, message: '请输入API请求限制' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>
        </Card>
        
        <Title level={4}>存储配置</Title>
        <Card style={{ marginBottom: 16 }}>
          <Form.Item
            name="storageProvider"
            label="存储提供商"
            rules={[{ required: true }]}
          >
            <Select onChange={(value) => setStorageType(value as 'local' | 's3')}>
              <Option value="local">本地存储</Option>
              <Option value="s3">S3兼容存储</Option>
            </Select>
          </Form.Item>
          
          {storageType === 's3' && (
            <div>
              <Form.Item
                name={['s3Config', 'bucket']}
                label="S3 存储桶名称"
                rules={[{ required: true, message: '请输入S3存储桶名称' }]}
              >
                <Input placeholder="例如: my-agents-bucket" />
              </Form.Item>
              
              <Form.Item
                name={['s3Config', 'region']}
                label="S3 区域"
                rules={[{ required: true, message: '请输入S3区域' }]}
              >
                <Input placeholder="例如: us-west-1" />
              </Form.Item>
              
              <Form.Item
                name={['s3Config', 'accessKey']}
                label="访问密钥ID"
                rules={[{ required: true, message: '请输入访问密钥ID' }]}
              >
                <Input placeholder="S3 Access Key ID" />
              </Form.Item>
              
              <Form.Item
                name={['s3Config', 'secretKey']}
                label="秘密访问密钥"
                rules={[{ required: true, message: '请输入秘密访问密钥' }]}
              >
                <Input.Password placeholder="S3 Secret Access Key" />
              </Form.Item>
            </div>
          )}
        </Card>
        
        <Collapse style={{ marginBottom: 16 }}>
          <Panel header="邮件服务配置" key="email">
            <Form.Item
              name={['emailSettings', 'smtpServer']}
              label="SMTP服务器"
            >
              <Input placeholder="例如: smtp.163.com" />
            </Form.Item>
            
            <Form.Item
              name={['emailSettings', 'smtpPort']}
              label="SMTP端口"
            >
              <InputNumber placeholder="例如: 465" style={{ width: '100%' }} />
            </Form.Item>
            
            <Form.Item
              name={['emailSettings', 'smtpUser']}
              label="SMTP用户名"
            >
              <Input placeholder="邮箱账号" />
            </Form.Item>
            
            <Form.Item
              name={['emailSettings', 'smtpPassword']}
              label="SMTP密码"
            >
              <Input.Password placeholder="邮箱密码或授权码" />
            </Form.Item>
            
            <Form.Item
              name={['emailSettings', 'senderEmail']}
              label="发送者邮箱"
            >
              <Input placeholder="显示的发件人地址" />
            </Form.Item>
          </Panel>
        </Collapse>
        
        <Form.Item>
          <Button 
            type="primary" 
            htmlType="submit" 
            icon={<SaveOutlined />} 
            loading={saving}
          >
            保存设置
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
};

export default SystemSettings;