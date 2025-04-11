import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { Radio, Spin, message } from 'antd';
import ChatAgentBuilder from './ChatAgentBuilder';
import FlowAgentBuilder from './FlowAgentBuilder';
import { getAgent } from '../../services/agentService';

const AgentBuilder: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const [loading, setLoading] = useState(false);
  const [agentType, setAgentType] = useState<'chat' | 'flow'>('chat');

  // 如果是编辑模式，加载智能体类型
  useEffect(() => {
    if (id) {
      fetchAgentType();
    } else {
      // 从URL参数获取初始类型
      const searchParams = new URLSearchParams(location.search);
      const typeParam = searchParams.get('type');
      if (typeParam === 'flow') {
        setAgentType('flow');
      }
    }
  }, [id, location]);

  // 获取智能体类型
  const fetchAgentType = async () => {
    try {
      setLoading(true);
      const data = await getAgent(id!);
      setAgentType(data.type as 'chat' | 'flow');
    } catch (error) {
      console.error('获取智能体详情失败:', error);
      message.error('获取智能体详情失败');
      // 错误时导航回智能体列表
      navigate('/agents');
    } finally {
      setLoading(false);
    }
  };

  // 处理类型变更
  const handleTypeChange = (e: any) => {
    const newType = e.target.value;
    setAgentType(newType);
    
    if (!id) {
      // 更新URL参数
      navigate(`/agents/create?type=${newType}`);
    }
  };

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '60vh' }}>
        <Spin size="large" tip="加载中..." />
      </div>
    );
  }

  return (
    <div>
      {!id && (
        <div style={{ marginBottom: 20 }}>
          <Radio.Group value={agentType} onChange={handleTypeChange} buttonStyle="solid">
            <Radio.Button value="chat">对话式智能体</Radio.Button>
            <Radio.Button value="flow">流程式智能体</Radio.Button>
          </Radio.Group>
        </div>
      )}

      {agentType === 'chat' ? <ChatAgentBuilder /> : <FlowAgentBuilder />}
    </div>
  );
};

export default AgentBuilder; 