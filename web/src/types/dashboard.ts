export interface IStatisticData {
  agentCount: number;
  conversationCount: number;
  userCount: number;
  tokenUsage: number;
}

export interface IUsageData {
  date: string;
  tokens: number;
  conversations: number;
}

export interface IAgentData {
  id: string;
  name: string;
  description: string;
  createdAt: string;
  lastAccessed: string;
  status: 'active' | 'inactive' | 'draft';
  type: 'chat' | 'flow';
  usageCount: number;
}

export interface IModelData {
  id: string;
  name: string;
  provider: string;
  type: 'local' | 'cloud';
  status: 'active' | 'inactive';
  lastUsed: string;
  contextLength: number;
  supportsFunctionCalling: boolean;
  apiKey?: string;
  baseUrl?: string;
  parameters?: {
    temperature?: number;
    topP?: number;
    maxTokens?: number;
    [key: string]: any;
  };
}

export interface ISystemSettings {
  siteName: string;
  logoUrl: string;
  apiRateLimit: number;
  allowRegistration: boolean;
  defaultLanguage: string;
  defaultModel: string;
  storageProvider: 'local' | 's3';
  s3Config?: {
    bucket: string;
    region: string;
    accessKey: string;
    secretKey: string;
  };
  emailSettings?: {
    smtpServer: string;
    smtpPort: number;
    smtpUser: string;
    smtpPassword: string;
    senderEmail: string;
  };
}

export interface ITopAgent {
  id: string;
  name: string;
  usageCount: number;
  tokenUsage: number;
}

export interface IRecentActivity {
  id: string;
  type: 'agent_created' | 'agent_updated' | 'conversation' | 'knowledge_updated';
  agentId?: string;
  agentName?: string;
  timestamp: string;
  user: string;
  details: string;
} 