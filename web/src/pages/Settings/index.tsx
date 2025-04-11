import React, { useEffect, useState } from 'react';
import { Tabs, Typography, Card, Divider } from 'antd';
import ModelSettings from './ModelSettings';
import SystemSettings from './SystemSettings';

const { Title } = Typography;

// 设置页面的主组件
const Settings: React.FC = () => {
  return (
    <div>
      <Title level={3}>系统设置</Title>
      <Divider style={{ margin: '16px 0' }} />
      
      <Card>
        <Tabs
          defaultActiveKey="system"
          items={[
            {
              key: 'system',
              label: '基本设置',
              children: <SystemSettings />,
            },
            {
              key: 'models',
              label: '模型配置',
              children: <ModelSettings />,
            }
          ]}
        />
      </Card>
    </div>
  );
};

export default Settings; 