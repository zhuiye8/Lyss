import React, { useState, useEffect, useRef } from 'react';
import { 
  Card, 
  Input, 
  Button, 
  Spin, 
  message, 
  Space,
  Form,
  Typography
} from 'antd';
import { SendOutlined } from '@ant-design/icons';
import { 
  XProvider, 
  Conversations,
  Bubble, 
  Sender, 
  Attachments,
  ThoughtChain
} from '@ant-design/x';
import { 
  IAgentConversation, 
  IAgentMessage,
  IToolCall
} from '../../types/agent';
import { sendMessage, getConversation } from '../../services/agentService';

const { Title, Text } = Typography;

interface ChatInterfaceProps {
  conversationId: string;
  initialMessages: IAgentMessage[];
  agentName?: string;
  variables?: Record<string, any>;
  onMessageSent?: () => void;
  loading?: boolean;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({
  conversationId,
  initialMessages = [],
  agentName = '智能助手',
  variables = {},
  onMessageSent,
  loading = false
}) => {
  const [value, setValue] = useState('');
  const [sending, setSending] = useState(false);
  const [messages, setMessages] = useState<IAgentMessage[]>(initialMessages);
  const [showThoughtChain, setShowThoughtChain] = useState(false);
  const [activeMessage, setActiveMessage] = useState<IAgentMessage | null>(null);

  // 处理发送消息
  const handleSend = async () => {
    if (!value.trim() || !conversationId) return;

    try {
      setSending(true);
      
      // 将用户消息添加到UI
      const userMessage: IAgentMessage = {
        id: `temp-${Date.now()}`,
        role: 'user',
        content: value,
        timestamp: new Date().toISOString()
      };
      
      setMessages(prev => [...prev, userMessage]);
      
      // 发送消息到API
      await sendMessage(conversationId, value, variables);
      
      // 刷新对话数据
      const updatedConversation = await getConversation(conversationId);
      setMessages(updatedConversation.messages);
      
      // 清空输入框
      setValue('');
      
      // 调用回调
      if (onMessageSent) {
        onMessageSent();
      }
    } catch (error) {
      console.error('发送消息失败:', error);
      message.error('发送消息失败');
    } finally {
      setSending(false);
    }
  };

  // 转换消息格式以适应Ant Design X组件
  const transformMessages = () => {
    return messages.map(msg => {
      const isUser = msg.role === 'user';
      
      return {
        id: msg.id,
        content: msg.content,
        position: isUser ? 'right' : 'left',
        type: isUser ? 'text' : 'text',
        avatar: isUser ? undefined : agentName.charAt(0),
        time: new Date(msg.timestamp),
        hasThoughtChain: !isUser && msg.toolCalls && msg.toolCalls.length > 0,
        toolCalls: msg.toolCalls
      };
    });
  };

  // 处理点击思维链
  const handleClickThoughtChain = (message: any) => {
    const originalMessage = messages.find(m => m.id === message.id);
    if (originalMessage && originalMessage.toolCalls) {
      setActiveMessage(originalMessage);
      setShowThoughtChain(true);
    }
  };

  // 渲染思维链内容
  const renderThoughtChain = () => {
    if (!activeMessage || !activeMessage.toolCalls) return null;

    const thoughts = activeMessage.toolCalls.map((call: IToolCall) => ({
      type: 'tool',
      tool: call.toolId,
      params: call.params,
      result: call.result,
      error: call.error
    }));

    return (
      <ThoughtChain
        open={showThoughtChain}
        title="思维链"
        onClose={() => setShowThoughtChain(false)}
        thoughts={thoughts}
      />
    );
  };
  
  return (
    <XProvider>
      <Card 
        title={<Title level={5}>{agentName}</Title>} 
        style={{ width: '100%', height: '100%', borderRadius: '8px' }}
        bodyStyle={{ padding: '12px', height: 'calc(100% - 56px)', display: 'flex', flexDirection: 'column' }}
        loading={loading}
      >
        <div style={{ flex: 1, overflow: 'auto', marginBottom: '12px' }}>
          <Conversations 
            data={transformMessages()} 
            onClickThoughtChain={handleClickThoughtChain}
            renderItem={(message) => (
              <Bubble
                content={message.content}
                position={message.position as 'left' | 'right'}
                time={message.time}
                avatar={message.avatar}
                hasThoughtChain={message.hasThoughtChain}
              />
            )}
          />
        </div>
        
        <div>
          <Sender
            value={value}
            onChange={setValue}
            onSend={handleSend}
            sending={sending}
            placeholder="输入消息..."
            sendText="发送"
            enterToSend
          />
        </div>

        {renderThoughtChain()}
      </Card>
    </XProvider>
  );
};

export default ChatInterface; 