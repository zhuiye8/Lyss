import React from 'react';
import { Card, Typography, Button, Space } from 'antd';
import { RobotOutlined, RocketOutlined } from '@ant-design/icons';
import { XProvider, Welcome } from '@ant-design/x';

const { Title, Paragraph } = Typography;

interface AIWelcomeProps {
  agentName?: string;
  agentDescription?: string;
  onStart?: () => void;
  welcomeMessage?: string;
}

const AIWelcome: React.FC<AIWelcomeProps> = ({
  agentName = '智能助手',
  agentDescription = '我是一个功能强大的智能助手，可以帮助你完成各种任务。',
  onStart,
  welcomeMessage = '有什么我可以帮助你的吗？'
}) => {
  return (
    <XProvider>
      <Card 
        bordered={false}
        style={{ borderRadius: '8px', maxWidth: '800px', margin: '0 auto' }}
      >
        <Welcome
          title={agentName}
          description={agentDescription}
          avatar={<RobotOutlined style={{ fontSize: '36px', color: '#1890ff' }} />}
          actions={[
            {
              key: 'start',
              text: '开始对话',
              icon: <RocketOutlined />,
              onClick: onStart
            }
          ]}
          welcomeText={welcomeMessage}
        />
      </Card>
    </XProvider>
  );
};

export default AIWelcome; 