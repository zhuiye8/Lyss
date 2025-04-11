import React from 'react';
import { Row, Col, Card, Statistic, Typography, List, Avatar } from 'antd';
import { RobotOutlined, DatabaseOutlined, ApiOutlined, ClockCircleOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';

const { Title, Paragraph } = Typography;

const Dashboard: React.FC = () => {
  const navigate = useNavigate();

  // 模拟数据
  const statistics = [
    { title: '智能体数量', value: 5, icon: <RobotOutlined />, color: '#1890ff' },
    { title: '知识库数量', value: 3, icon: <DatabaseOutlined />, color: '#52c41a' },
    { title: 'API调用次数', value: 1205, icon: <ApiOutlined />, color: '#722ed1' },
    { title: '平均响应时间', value: '320ms', icon: <ClockCircleOutlined />, color: '#fa8c16' },
  ];

  const recentActivities = [
    { title: '创建了新智能体"客服助手"', time: '10分钟前' },
    { title: '更新了"产品手册"知识库', time: '2小时前' },
    { title: '修改了系统设置', time: '昨天' },
    { title: '添加了新文档到"API文档"知识库', time: '2天前' },
  ];

  return (
    <div>
      <div style={{ marginBottom: 24 }}>
        <Title level={2}>仪表盘</Title>
        <Paragraph>欢迎使用智能体构建平台，您可以在这里查看系统状态和最近活动。</Paragraph>
      </div>

      <Row gutter={[16, 16]}>
        {statistics.map((stat, index) => (
          <Col key={index} xs={24} sm={12} md={6}>
            <Card hoverable>
              <Statistic 
                title={stat.title}
                value={stat.value}
                valueStyle={{ color: stat.color }}
                prefix={stat.icon}
              />
            </Card>
          </Col>
        ))}
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col span={24}>
          <Card title="最近活动" bordered={false}>
            <List
              itemLayout="horizontal"
              dataSource={recentActivities}
              renderItem={(item) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={<Avatar icon={<ClockCircleOutlined />} />}
                    title={item.title}
                    description={item.time}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 24 }}>
        <Col xs={24} md={12}>
          <Card 
            title="快速创建智能体" 
            bordered={false}
            hoverable
            onClick={() => navigate('/agents')}
          >
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <Avatar size={64} icon={<RobotOutlined />} style={{ backgroundColor: '#1890ff' }} />
              <div style={{ marginLeft: 16 }}>
                <Title level={4}>开始创建您的智能体</Title>
                <Paragraph>快速构建能够满足您业务需求的AI助手</Paragraph>
              </div>
            </div>
          </Card>
        </Col>

        <Col xs={24} md={12}>
          <Card 
            title="管理知识库" 
            bordered={false}
            hoverable
            onClick={() => navigate('/knowledge')}
          >
            <div style={{ display: 'flex', alignItems: 'center' }}>
              <Avatar size={64} icon={<DatabaseOutlined />} style={{ backgroundColor: '#52c41a' }} />
              <div style={{ marginLeft: 16 }}>
                <Title level={4}>知识库管理</Title>
                <Paragraph>为您的智能体添加专属知识库，提升回答准确度</Paragraph>
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard; 