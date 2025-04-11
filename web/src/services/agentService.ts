import api from './api';
import { 
  IAgentFormState, 
  IChatAgentConfig, 
  IFlowAgentConfig, 
  IPromptTemplate,
  IKnowledgeBase,
  ITool,
  IAgentConversation,
  IAgentMessage,
  ITestResult
} from '../types/agent';

// 获取单个智能体
export const getAgent = async (id: string): Promise<IAgentFormState> => {
  const response = await api.get(`/agents/${id}`);
  return response.data;
};

// 创建新智能体
export const createAgent = async (agent: IAgentFormState): Promise<IAgentFormState> => {
  const response = await api.post('/agents', agent);
  return response.data;
};

// 更新智能体
export const updateAgent = async (id: string, agent: Partial<IAgentFormState>): Promise<IAgentFormState> => {
  const response = await api.put(`/agents/${id}`, agent);
  return response.data;
};

// 获取提示词模板列表
export const getPromptTemplates = async (): Promise<IPromptTemplate[]> => {
  const response = await api.get('/prompt-templates');
  return response.data;
};

// 创建提示词模板
export const createPromptTemplate = async (template: Partial<IPromptTemplate>): Promise<IPromptTemplate> => {
  const response = await api.post('/prompt-templates', template);
  return response.data;
};

// 更新提示词模板
export const updatePromptTemplate = async (id: string, template: Partial<IPromptTemplate>): Promise<IPromptTemplate> => {
  const response = await api.put(`/prompt-templates/${id}`, template);
  return response.data;
};

// 删除提示词模板
export const deletePromptTemplate = async (id: string): Promise<void> => {
  await api.delete(`/prompt-templates/${id}`);
};

// 获取知识库列表
export const getKnowledgeBases = async (): Promise<IKnowledgeBase[]> => {
  const response = await api.get('/knowledge-bases');
  return response.data;
};

// 获取工具列表
export const getTools = async (): Promise<ITool[]> => {
  const response = await api.get('/tools');
  return response.data;
};

// 创建对话
export const createConversation = async (agentId: string): Promise<IAgentConversation> => {
  const response = await api.post(`/agents/${agentId}/conversations`);
  return response.data;
};

// 获取对话历史
export const getConversation = async (conversationId: string): Promise<IAgentConversation> => {
  const response = await api.get(`/conversations/${conversationId}`);
  return response.data;
};

// 发送消息
export const sendMessage = async (
  conversationId: string, 
  message: string, 
  variables?: Record<string, any>
): Promise<IAgentMessage> => {
  const response = await api.post(`/conversations/${conversationId}/messages`, {
    content: message,
    variables
  });
  return response.data;
};

// 获取智能体对话历史列表
export const getAgentConversations = async (
  agentId: string,
  page: number = 1,
  pageSize: number = 10
): Promise<{
  data: IAgentConversation[];
  total: number;
}> => {
  const response = await api.get(`/agents/${agentId}/conversations`, {
    params: { page, pageSize }
  });
  return response.data;
};

// 测试智能体
export const testAgent = async (
  agentId: string,
  input: string,
  variables?: Record<string, any>
): Promise<ITestResult> => {
  const response = await api.post(`/agents/${agentId}/test`, {
    input,
    variables
  });
  return response.data;
};

// 获取测试历史
export const getTestHistory = async (
  agentId: string,
  page: number = 1,
  pageSize: number = 10
): Promise<{
  data: ITestResult[];
  total: number;
}> => {
  const response = await api.get(`/agents/${agentId}/tests`, {
    params: { page, pageSize }
  });
  return response.data;
};

// 发布智能体
export const publishAgent = async (agentId: string): Promise<IAgentFormState> => {
  const response = await api.post(`/agents/${agentId}/publish`);
  return response.data;
};

// 导出智能体配置
export const exportAgentConfig = async (agentId: string): Promise<IChatAgentConfig | IFlowAgentConfig> => {
  const response = await api.get(`/agents/${agentId}/export`);
  return response.data;
};

// 导入智能体配置
export const importAgentConfig = async (config: IChatAgentConfig | IFlowAgentConfig): Promise<IAgentFormState> => {
  const response = await api.post('/agents/import', config);
  return response.data;
}; 