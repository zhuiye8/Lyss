import React from 'react';
import { Card, Radio, Spin } from 'antd';
import { IUsageData } from '../../types/dashboard';

interface UsageTrendChartProps {
  data: IUsageData[];
  loading: boolean;
  timeRange: '7d' | '30d' | '90d';
  onTimeRangeChange: (range: '7d' | '30d' | '90d') => void;
}

const UsageTrendChart: React.FC<UsageTrendChartProps> = ({
  data,
  loading,
  timeRange,
  onTimeRangeChange
}) => {
  return (
    <Card
      title="使用趋势"
      extra={
        <Radio.Group 
          value={timeRange}
          onChange={e => onTimeRangeChange(e.target.value)}
          optionType="button"
          buttonStyle="solid"
          size="small"
        >
          <Radio.Button value="7d">7天</Radio.Button>
          <Radio.Button value="30d">30天</Radio.Button>
          <Radio.Button value="90d">90天</Radio.Button>
        </Radio.Group>
      }
    >
      {loading ? (
        <div style={{ height: 300, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
          <Spin />
        </div>
      ) : (
        <div style={{ height: 300, padding: '20px 0' }}>
          {/* 简化版，只显示数据而不绘制图表 */}
          <div style={{ maxHeight: 260, overflowY: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
              <thead>
                <tr style={{ borderBottom: '1px solid #f0f0f0' }}>
                  <th style={{ padding: '8px', textAlign: 'left' }}>日期</th>
                  <th style={{ padding: '8px', textAlign: 'right' }}>对话数</th>
                  <th style={{ padding: '8px', textAlign: 'right' }}>Token 使用量</th>
                </tr>
              </thead>
              <tbody>
                {data.map((item, index) => (
                  <tr key={index} style={{ borderBottom: '1px solid #f0f0f0' }}>
                    <td style={{ padding: '8px', textAlign: 'left' }}>{item.date}</td>
                    <td style={{ padding: '8px', textAlign: 'right' }}>{item.conversations}</td>
                    <td style={{ padding: '8px', textAlign: 'right' }}>{item.tokens.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </Card>
  );
};

export default UsageTrendChart; 