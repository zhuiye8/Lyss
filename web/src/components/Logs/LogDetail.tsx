import React from 'react';
import { Descriptions, Tag, Typography, Collapse, Badge } from 'antd';
import moment from 'moment';
import ReactJson from 'react-json-view';
import { Log, APILog, ErrorLog, ModelCallLog, LogLevel } from '../../types/log';

const { Panel } = Collapse;
const { Text } = Typography;

interface LogDetailProps {
  log: Log | APILog | ErrorLog | ModelCallLog;
}

const LogDetail: React.FC<LogDetailProps> = ({ log }) => {
  // 日志级别对应的标签颜色
  const getLevelColor = (level: LogLevel) => {
    switch (level) {
      case LogLevel.DEBUG:
        return 'blue';
      case LogLevel.INFO:
        return 'green';
      case LogLevel.WARN:
        return 'orange';
      case LogLevel.ERROR:
        return 'red';
      case LogLevel.FATAL:
        return 'purple';
      default:
        return 'default';
    }
  };

  // 判断日志类型
  const isAPILog = (log: any): log is APILog => 'method' in log && 'path' in log;
  const isErrorLog = (log: any): log is ErrorLog => 'stack_trace' in log;
  const isModelCallLog = (log: any): log is ModelCallLog => 'model_name' in log;

  return (
    <div>
      <Descriptions bordered column={2} size="small">
        <Descriptions.Item label="ID" span={2}>{log.id}</Descriptions.Item>
        <Descriptions.Item label="级别">
          <Tag color={getLevelColor(log.level)}>{log.level.toUpperCase()}</Tag>
        </Descriptions.Item>
        <Descriptions.Item label="类别">
          <Tag>{log.category}</Tag>
        </Descriptions.Item>
        <Descriptions.Item label="时间" span={2}>
          {moment(log.created_at).format('YYYY-MM-DD HH:mm:ss')}
        </Descriptions.Item>
        <Descriptions.Item label="消息" span={2}>
          <Text>{log.message}</Text>
        </Descriptions.Item>
        
        {log.user_id && (
          <Descriptions.Item label="用户ID">{log.user_id}</Descriptions.Item>
        )}
        
        {/* API日志特定信息 */}
        {isAPILog(log) && (
          <>
            <Descriptions.Item label="HTTP方法">
              <Tag color="blue">{log.method}</Tag>
            </Descriptions.Item>
            <Descriptions.Item label="API路径" span={2}>
              {log.path}
            </Descriptions.Item>
            <Descriptions.Item label="状态码">
              <Badge 
                status={log.status_code < 400 ? 'success' : 'error'} 
                text={log.status_code} 
              />
            </Descriptions.Item>
            <Descriptions.Item label="响应时间">
              {log.duration}ms
            </Descriptions.Item>
            {log.ip && (
              <Descriptions.Item label="IP地址">{log.ip}</Descriptions.Item>
            )}
            {log.request_id && (
              <Descriptions.Item label="请求ID">{log.request_id}</Descriptions.Item>
            )}
          </>
        )}
        
        {/* 错误日志特定信息 */}
        {isErrorLog(log) && (
          <>
            {log.error_code && (
              <Descriptions.Item label="错误代码">{log.error_code}</Descriptions.Item>
            )}
            {log.source && (
              <Descriptions.Item label="错误来源">{log.source}</Descriptions.Item>
            )}
            {log.resolved_at && (
              <Descriptions.Item label="解决时间">
                {moment(log.resolved_at).format('YYYY-MM-DD HH:mm:ss')}
              </Descriptions.Item>
            )}
            {log.resolved_by && (
              <Descriptions.Item label="解决人">{log.resolved_by}</Descriptions.Item>
            )}
          </>
        )}
        
        {/* 模型调用日志特定信息 */}
        {isModelCallLog(log) && (
          <>
            <Descriptions.Item label="模型名称">{log.model_name}</Descriptions.Item>
            <Descriptions.Item label="调用结果">
              <Tag color={log.success ? 'green' : 'red'}>
                {log.success ? '成功' : '失败'}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="耗时">{log.duration}ms</Descriptions.Item>
            {log.total_tokens && (
              <Descriptions.Item label="总Token数">{log.total_tokens}</Descriptions.Item>
            )}
            {log.prompt_tokens && (
              <Descriptions.Item label="提示Token">{log.prompt_tokens}</Descriptions.Item>
            )}
            {log.comp_tokens && (
              <Descriptions.Item label="完成Token">{log.comp_tokens}</Descriptions.Item>
            )}
            {log.application_id && (
              <Descriptions.Item label="应用ID">{log.application_id}</Descriptions.Item>
            )}
            {log.project_id && (
              <Descriptions.Item label="项目ID">{log.project_id}</Descriptions.Item>
            )}
          </>
        )}
      </Descriptions>
      
      {/* 展示堆栈跟踪（如果有） */}
      {isErrorLog(log) && log.stack_trace && (
        <Collapse style={{ marginTop: 16 }}>
          <Panel header="堆栈跟踪" key="1">
            <pre style={{ 
              maxHeight: '300px', 
              overflow: 'auto', 
              backgroundColor: '#f5f5f5', 
              padding: '10px',
              borderRadius: '4px'
            }}>
              {log.stack_trace}
            </pre>
          </Panel>
        </Collapse>
      )}
      
      {/* 展示元数据（如果有） */}
      {log.metadata && (
        <Collapse style={{ marginTop: 16 }}>
          <Panel header="元数据" key="1">
            <ReactJson 
              src={log.metadata} 
              name={false} 
              collapsed={2} 
              displayDataTypes={false}
              displayObjectSize={false}
            />
          </Panel>
        </Collapse>
      )}
    </div>
  );
};

export default LogDetail; 