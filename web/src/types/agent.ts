// 智能体类型定义

export interface IPromptTemplate {
  id: string;
  content: string;
  variables: string[];
  description?: string;
  name?: string;
}

export interface IKnowledgeBase {
  id: string;
  name: string;
  description?: string;
  documentCount: number;
  vectorCount: number;
  lastUpdated: string;
}

export interface ITool {
  id: string;
  name: string;
  description: string;
  type: 'api' | 'function' | 'builtin';
  parameters: IToolParameter[];
  returnType: string;
  icon?: string;
}

export interface IToolParameter {
  name: string;
  type: 'string' | 'number' | 'boolean' | 'array' | 'object';
  required: boolean;
  description?: string;
  default?: any;
  enum?: string[];
}

export interface IAgentNode {
  id: string;
  type: 'start' | 'message' | 'condition' | 'tool' | 'end' | 'knowledge';
  position: { x: number; y: number };
  data: {
    prompt?: string;
    condition?: string;
    toolId?: string;
    toolParams?: Record<string, any>;
    knowledgeBaseId?: string;
    message?: string;
  };
}

export interface IAgentEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  condition?: string;
}

export interface IAgentConversation {
  id: string;
  agentId: string;
  messages: IAgentMessage[];
  createdAt: string;
  updatedAt: string;
  variables: Record<string, any>;
}

export interface IAgentMessage {
  id: string;
  role: 'user' | 'assistant' | 'system' | 'tool';
  content: string;
  timestamp: string;
  toolCalls?: IToolCall[];
  toolCallId?: string;
}

export interface IToolCall {
  id: string;
  toolId: string;
  params: Record<string, any>;
  result?: any;
  error?: string;
}

export interface IChatAgentConfig {
  id: string;
  name: string;
  description?: string;
  welcomeMessage: string;
  promptTemplate: IPromptTemplate;
  modelId: string;
  knowledgeBases: string[];
  tools: string[];
  systemPrompt: string;
  parameters: {
    temperature: number;
    topP: number;
    maxTokens: number;
    presencePenalty: number;
    frequencyPenalty: number;
  };
}

export interface IFlowAgentConfig {
  id: string;
  name: string;
  description?: string;
  nodes: IAgentNode[];
  edges: IAgentEdge[];
  modelId: string;
  knowledgeBases: string[];
  tools: string[];
  variables: string[];
}

export type IAgentConfig = IChatAgentConfig | IFlowAgentConfig;

export interface IAgentFormState {
  id?: string;
  name: string;
  description: string;
  type: 'chat' | 'flow';
  visibility: 'public' | 'private';
  modelId: string;
  chatConfig?: Partial<IChatAgentConfig>;
  flowConfig?: Partial<IFlowAgentConfig>;
  status: 'active' | 'inactive' | 'draft';
}

export interface ITestResult {
  id: string;
  agentId: string;
  timestamp: string;
  input: string;
  output: string;
  duration: number;
  tokenUsage: {
    prompt: number;
    completion: number;
    total: number;
  };
  success: boolean;
  error?: string;
  modelId: string;
  conversationId: string;
} 