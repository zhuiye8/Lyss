import React, { useState, useEffect } from 'react';
import { Row, Col, Card, Statistic, Table, DatePicker, Button, Select, Spin } from 'antd';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { getLogStats, getMetrics } from '../../services/logService';
import { LogStatsResponse, MetricsResponse } from '../../types/log';

const { RangePicker } = DatePicker;
const { Option } = Select;

const LogStats: React.FC = () => {
  const [loading, setLoading] = useState<boolean>(false);
  const [metricsLoading, setMetricsLoading] = useState<boolean>(false);
  const [stats, setStats] = useState<LogStatsResponse['data'] | null>(null);
  const [metrics, setMetrics] = useState<MetricsResponse['data'] | null>(null);
  const [timeRange, setTimeRange] = useState<[Date, Date] | null>(null);
  const [metricsRange, setMetricsRange] = useState<string>('1h');

  // 获取日志统计数据
  const fetchStats = async () => {
    setLoading(true);
    try {
      let startTime, endTime;
      if (timeRange && timeRange.length === 2) {
        startTime = timeRange[0].toISOString();
        endTime = timeRange[1].toISOString();
      }

      const response = await getLogStats(startTime, endTime);
      const data = response as LogStatsResponse;
      setStats(data.data);
    } catch (error) {
      console.error('获取日志统计信息失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 获取系统监控指标
  const fetchMetrics = async () => {
    setMetricsLoading(true);
    try {
      const response = await getMetrics(metricsRange);
      const data = response as MetricsResponse;
      setMetrics(data.data);
    } catch (error) {
      console.error('获取系统监控指标失败:', error);
    } finally {
      setMetricsLoading(false);
    }
  };

  // 初始加载
  useEffect(() => {
    fetchStats();
    fetchMetrics();
  }, []);

  // 系统指标时间范围变化时重新获取数据
  useEffect(() => {
    fetchMetrics();
  }, [metricsRange]);

  // 格式化指标数据为图表可用格式
  const formatMetricsForChart = (metricData: any[]) => {
    if (!metricData) return [];

    return metricData.map(item => ({
      time: new Date(item.created_at).toLocaleTimeString(),
      value: item.metric_value,
    }));
  };

  // 错误分类统计表格列
  const errorColumns = [
    {
      title: '分类',
      dataIndex: 'category',
      key: 'category',
    },
    {
      title: '错误数量',
      dataIndex: 'count',
      key: 'count',
      sorter: (a: any, b: any) => a.count - b.count,
    },
    {
      title: '占比',
      dataIndex: 'percentage',
      key: 'percentage',
    },
  ];

  // 将对象格式的统计数据转换为表格数据
  const formatStatsForTable = (data: Record<string, number>, total: number) => {
    if (!data) return [];

    return Object.entries(data).map(([key, value]) => ({
      category: key,
      count: value,
      percentage: `${((value / total) * 100).toFixed(2)}%`,
    }));
  };

  return (
    <div>
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card 
            title="日志统计信息" 
            extra={
              <div>
                <RangePicker 
                  onChange={(dates) => setTimeRange(dates as [Date, Date])} 
                  style={{ marginRight: 8 }}
                />
                <Button type="primary" onClick={fetchStats}>
                  查询
                </Button>
              </div>
            }
          >
            <Spin spinning={loading}>
              {stats ? (
                <Row gutter={[16, 16]}>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="总日志数"
                        value={stats.total_logs}
                        valueStyle={{ color: '#3f8600' }}
                      />
                    </Card>
                  </Col>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="错误数量"
                        value={stats.error_count}
                        valueStyle={{ color: '#cf1322' }}
                      />
                    </Card>
                  </Col>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="API请求数"
                        value={stats.api_count}
                        valueStyle={{ color: '#1890ff' }}
                      />
                    </Card>
                  </Col>
                  <Col span={6}>
                    <Card>
                      <Statistic
                        title="模型调用数"
                        value={stats.model_call_count}
                        valueStyle={{ color: '#722ed1' }}
                      />
                    </Card>
                  </Col>
                  
                  <Col span={12}>
                    <Card title="错误分类统计">
                      <Table 
                        columns={errorColumns} 
                        dataSource={formatStatsForTable(stats.error_by_category, stats.error_count)}
                        pagination={false}
                        size="small"
                      />
                    </Card>
                  </Col>
                  
                  <Col span={12}>
                    <Card title="日志级别分布">
                      <Table 
                        columns={[
                          { title: '级别', dataIndex: 'category', key: 'category' },
                          { title: '数量', dataIndex: 'count', key: 'count' },
                          { title: '占比', dataIndex: 'percentage', key: 'percentage' },
                        ]} 
                        dataSource={formatStatsForTable(stats.logs_by_level, stats.total_logs)}
                        pagination={false}
                        size="small"
                      />
                    </Card>
                  </Col>
                </Row>
              ) : (
                <div style={{ textAlign: 'center', padding: '20px' }}>
                  暂无统计数据
                </div>
              )}
            </Spin>
          </Card>
        </Col>

        <Col span={24}>
          <Card 
            title="系统监控指标" 
            extra={
              <Select 
                value={metricsRange}
                onChange={setMetricsRange}
                style={{ width: 120 }}
              >
                <Option value="15m">15分钟</Option>
                <Option value="30m">30分钟</Option>
                <Option value="1h">1小时</Option>
                <Option value="6h">6小时</Option>
                <Option value="12h">12小时</Option>
                <Option value="24h">24小时</Option>
                <Option value="7d">7天</Option>
              </Select>
            }
          >
            <Spin spinning={metricsLoading}>
              {metrics ? (
                <Row gutter={[16, 16]}>
                  <Col span={12}>
                    <Card title="CPU使用率">
                      <ResponsiveContainer width="100%" height={200}>
                        <LineChart
                          data={formatMetricsForChart(metrics.cpu_usage)}
                          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                        >
                          <CartesianGrid strokeDasharray="3 3" />
                          <XAxis dataKey="time" />
                          <YAxis />
                          <Tooltip />
                          <Legend />
                          <Line type="monotone" dataKey="value" stroke="#8884d8" name="CPU使用率(%)" />
                        </LineChart>
                      </ResponsiveContainer>
                    </Card>
                  </Col>
                  
                  <Col span={12}>
                    <Card title="内存使用率">
                      <ResponsiveContainer width="100%" height={200}>
                        <LineChart
                          data={formatMetricsForChart(metrics.memory_usage)}
                          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                        >
                          <CartesianGrid strokeDasharray="3 3" />
                          <XAxis dataKey="time" />
                          <YAxis />
                          <Tooltip />
                          <Legend />
                          <Line type="monotone" dataKey="value" stroke="#82ca9d" name="内存使用率(%)" />
                        </LineChart>
                      </ResponsiveContainer>
                    </Card>
                  </Col>
                  
                  <Col span={12}>
                    <Card title="API延迟">
                      <ResponsiveContainer width="100%" height={200}>
                        <LineChart
                          data={formatMetricsForChart(metrics.api_latency)}
                          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                        >
                          <CartesianGrid strokeDasharray="3 3" />
                          <XAxis dataKey="time" />
                          <YAxis />
                          <Tooltip />
                          <Legend />
                          <Line type="monotone" dataKey="value" stroke="#1890ff" name="API平均延迟(ms)" />
                        </LineChart>
                      </ResponsiveContainer>
                    </Card>
                  </Col>
                  
                  <Col span={12}>
                    <Card title="模型调用延迟">
                      <ResponsiveContainer width="100%" height={200}>
                        <LineChart
                          data={formatMetricsForChart(metrics.model_latency)}
                          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                        >
                          <CartesianGrid strokeDasharray="3 3" />
                          <XAxis dataKey="time" />
                          <YAxis />
                          <Tooltip />
                          <Legend />
                          <Line type="monotone" dataKey="value" stroke="#722ed1" name="模型平均延迟(ms)" />
                        </LineChart>
                      </ResponsiveContainer>
                    </Card>
                  </Col>
                </Row>
              ) : (
                <div style={{ textAlign: 'center', padding: '20px' }}>
                  暂无监控数据
                </div>
              )}
            </Spin>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default LogStats; 