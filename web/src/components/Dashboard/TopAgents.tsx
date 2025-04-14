import React from 'react';
import { Card, List, Avatar, Spin, Progress } from 'antd';
import { RobotOutlined } from '@ant-design/icons';
import { ITopAgent } from '../../types/dashboard';

interface TopAgentsProps {
  agents: ITopAgent[];
  loading: boolean;
}

const TopAgents: React.FC<TopAgentsProps> = ({ agents, loading }) => {
  return (
    <Card title="热门智能体">
      {loading ? (
        <div style={{ display: 'flex', justifyContent: 'center', padding: '20px 0' }}>
          <Spin />
        </div>
      ) : (
        <List
          dataSource={agents}
          renderItem={agent => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon={<RobotOutlined />} style={{ backgroundColor: '#1890ff' }} />}
                title={agent.name}
                description={
          <Progress 
                    percent={Math.round(agent.successRate * 100)} 
                    size="small" 
                    status="active"
                    format={percent => `${percent}% 成功率`}
                  />
                }
              />
              <div>
                <span style={{ color: '#8c8c8c' }}>使用次数:</span> {agent.usage}
        </div>
            </List.Item>
          )}
        />
      )}
    </Card>
  );
};

export default TopAgents; 