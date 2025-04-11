import React, { useState } from 'react';
import { Card, Typography, Space } from 'antd';
import { 
  XProvider, 
  Suggestion 
} from '@ant-design/x';

const { Title } = Typography;

interface SuggestionPromptsProps {
  onSelect: (prompt: string) => void;
  suggestions?: string[];
  title?: string;
}

const SuggestionPrompts: React.FC<SuggestionPromptsProps> = ({
  onSelect,
  suggestions = [],
  title = '快捷提示'
}) => {
  // 如果没有提供建议，使用默认建议
  const defaultSuggestions = [
    '你能帮我做什么?',
    '解释一下你的功能',
    '给我一个使用示例',
    '你的知识范围是什么?',
    '你可以接入哪些工具?'
  ];

  const promptList = suggestions.length > 0 ? suggestions : defaultSuggestions;

  return (
    <XProvider>
      <Card title={title} size="small">
        <Space direction="vertical" style={{ width: '100%' }}>
          <Suggestion
            data={promptList.map(text => ({ text }))}
            onSelect={(item) => onSelect(item.text)}
            maxLines={2}
          />
        </Space>
      </Card>
    </XProvider>
  );
};

export default SuggestionPrompts; 