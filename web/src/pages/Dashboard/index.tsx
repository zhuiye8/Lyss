import React, { useEffect, useState } from 'react';
import { Row, Col, Typography, Divider } from 'antd';
import StatisticCard from '../../components/Dashboard/StatisticCard';
import UsageTrendChart from '../../components/Dashboard/UsageTrendChart';
import TopAgents from '../../components/Dashboard/TopAgents';
import RecentActivities from '../../components/Dashboard/RecentActivities';
import { 
  getStatistics, 
  getUsageTrend, 
  getTopAgents, 
  getRecentActivities 
} from '../../services/dashboardService';
import { 
  IStatisticData, 
  IUsageData, 
  ITopAgent, 
  IRecentActivity 
} from '../../types/dashboard';

const { Title } = Typography;

const Dashboard: React.FC = () => {
  // 状态管理
  const [statistics, setStatistics] = useState<IStatisticData>({
    agentCount: 0,
    conversationCount: 0,
    userCount: 0,
    tokenUsage: 0
  });
  const [usageData, setUsageData] = useState<IUsageData[]>([]);
  const [topAgents, setTopAgents] = useState<ITopAgent[]>([]);
  const [activities, setActivities] = useState<IRecentActivity[]>([]);
  const [timeRange, setTimeRange] = useState<'7d' | '30d' | '90d'>('7d');
  
  // 加载状态
  const [statisticsLoading, setStatisticsLoading] = useState(true);
  const [usageLoading, setUsageLoading] = useState(true);
  const [agentsLoading, setAgentsLoading] = useState(true);
  const [activitiesLoading, setActivitiesLoading] = useState(true);

  // 获取统计数据
  const fetchStatistics = async () => {
    try {
      setStatisticsLoading(true);
      const data = await getStatistics();
      setStatistics(data);
    } catch (error) {
      console.error('获取统计数据失败:', error);
    } finally {
      setStatisticsLoading(false);
    }
  };

  // 获取使用趋势数据
  const fetchUsageTrend = async (days: number) => {
    try {
      setUsageLoading(true);
      const data = await getUsageTrend(days);
      setUsageData(data);
    } catch (error) {
      console.error('获取使用趋势数据失败:', error);
    } finally {
      setUsageLoading(false);
    }
  };

  // 获取热门智能体
  const fetchTopAgents = async () => {
    try {
      setAgentsLoading(true);
      const data = await getTopAgents(5);
      setTopAgents(data);
    } catch (error) {
      console.error('获取热门智能体失败:', error);
    } finally {
      setAgentsLoading(false);
    }
  };

  // 获取最近活动
  const fetchRecentActivities = async () => {
    try {
      setActivitiesLoading(true);
      const data = await getRecentActivities(10);
      setActivities(data);
    } catch (error) {
      console.error('获取最近活动失败:', error);
    } finally {
      setActivitiesLoading(false);
    }
  };

  // 处理时间范围变化
  const handleTimeRangeChange = (range: '7d' | '30d' | '90d') => {
    setTimeRange(range);
    const days = range === '7d' ? 7 : range === '30d' ? 30 : 90;
    fetchUsageTrend(days);
  };

  // 初始加载数据
  useEffect(() => {
    fetchStatistics();
    fetchUsageTrend(7);
    fetchTopAgents();
    fetchRecentActivities();
  }, []);

  return (
    <div>
      <Title level={3}>系统概览</Title>
      <Divider style={{ margin: '16px 0' }} />
      
      {/* 统计卡片 */}
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={12} md={6}>
          <StatisticCard
            title="智能体总数"
            value={statistics.agentCount}
            type="agents"
            loading={statisticsLoading}
          />
        </Col>
        <Col xs={24} sm={12} md={6}>
          <StatisticCard
            title="对话总数"
            value={statistics.conversationCount}
            type="conversations"
            loading={statisticsLoading}
          />
        </Col>
        <Col xs={24} sm={12} md={6}>
          <StatisticCard
            title="用户总数"
            value={statistics.userCount}
            type="users"
            loading={statisticsLoading}
          />
        </Col>
        <Col xs={24} sm={12} md={6}>
          <StatisticCard
            title="Token 使用量"
            value={statistics.tokenUsage.toLocaleString()}
            type="tokens"
            loading={statisticsLoading}
          />
        </Col>
      </Row>

      {/* 图表和列表 */}
      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <UsageTrendChart
            data={usageData}
            loading={usageLoading}
            timeRange={timeRange}
            onTimeRangeChange={handleTimeRangeChange}
          />
        </Col>
        <Col xs={24} lg={8}>
          <RecentActivities
            activities={activities}
            loading={activitiesLoading}
          />
        </Col>
      </Row>

      {/* 热门智能体 */}
      <Row style={{ marginTop: 16 }}>
        <Col span={24}>
          <TopAgents
            agents={topAgents}
            loading={agentsLoading}
          />
        </Col>
      </Row>
    </div>
  );
};

export default Dashboard; 