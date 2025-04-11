import React from 'react';
import { 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend, 
  ResponsiveContainer 
} from 'recharts';
import { Card, Radio, Spin } from 'antd';
import { IUsageData } from '../../types/dashboard';

interface UsageTrendChartProps {
  data: IUsageData[];
  loading?: boolean;
  timeRange: '7d' | '30d' | '90d';
  onTimeRangeChange: (range: '7d' | '30d' | '90d') => void;
}

const UsageTrendChart: React.FC<UsageTrendChartProps> = ({
  data,
  loading = false,
  timeRange,
  onTimeRangeChange,
}) => {
  return (
    <Card
      title="使用趋势"
      extra={
        <Radio.Group 
          value={timeRange}
          onChange={(e) => onTimeRangeChange(e.target.value)}
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
        <div style={{ height: 300, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <Spin tip="加载中..." />
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={300}>
          <LineChart
            data={data}
            margin={{
              top: 5,
              right: 30,
              left: 20,
              bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis yAxisId="left" />
            <YAxis yAxisId="right" orientation="right" />
            <Tooltip />
            <Legend />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="tokens"
              name="Token 使用量"
              stroke="#8884d8"
              activeDot={{ r: 8 }}
            />
            <Line
              yAxisId="right"
              type="monotone"
              dataKey="conversations"
              name="对话数"
              stroke="#82ca9d"
            />
          </LineChart>
        </ResponsiveContainer>
      )}
    </Card>
  );
};

export default UsageTrendChart; 