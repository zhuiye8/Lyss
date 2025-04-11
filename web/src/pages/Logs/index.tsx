import React from 'react';
import { Tabs } from 'antd';
import LogList from '../../components/Logs/LogList';
import LogStats from '../../components/Logs/LogStats';
import { LogType } from '../../types/log';

const { TabPane } = Tabs;

const LogsPage: React.FC = () => {
  return (
    <div className="logs-page">
      <h2>日志管理</h2>
      <Tabs defaultActiveKey="all" size="large">
        <TabPane tab="全部日志" key="all">
          <LogList type={LogType.ALL} />
        </TabPane>
        
        <TabPane tab="API日志" key="api">
          <LogList type={LogType.API} />
        </TabPane>
        
        <TabPane tab="错误日志" key="error">
          <LogList type={LogType.ERROR} />
        </TabPane>
        
        <TabPane tab="模型调用日志" key="model_call">
          <LogList type={LogType.MODEL_CALL} />
        </TabPane>
        
        <TabPane tab="统计与监控" key="stats">
          <LogStats />
        </TabPane>
      </Tabs>
    </div>
  );
};

export default LogsPage; 