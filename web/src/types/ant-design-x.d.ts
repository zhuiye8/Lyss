declare module '@ant-design/x' {
  import { ReactNode, FC } from 'react';

  // XProvider props
  interface XProviderProps {
    children: ReactNode;
    theme?: {
      primaryColor?: string;
      borderRadius?: number;
      [key: string]: any;
    };
    settings?: {
      api?: {
        baseURL?: string;
        headers?: Record<string, string>;
        [key: string]: any;
      };
      ui?: {
        avatarShape?: 'circle' | 'square';
        darkMode?: boolean;
        [key: string]: any;
      };
      [key: string]: any;
    };
  }

  // Bubble props
  interface BubbleProps {
    content: ReactNode;
    position?: 'left' | 'right';
    time?: Date | string;
    avatar?: ReactNode;
    hasThoughtChain?: boolean;
    className?: string;
    style?: React.CSSProperties;
    onClick?: () => void;
  }

  // Conversations props
  interface ConversationMessage {
    id: string;
    content: ReactNode;
    position?: 'left' | 'right';
    type?: string;
    avatar?: ReactNode;
    time?: Date | string;
    hasThoughtChain?: boolean;
    [key: string]: any;
  }

  interface ConversationsProps {
    data: ConversationMessage[];
    renderItem?: (message: ConversationMessage) => ReactNode;
    onClickThoughtChain?: (message: ConversationMessage) => void;
    className?: string;
    style?: React.CSSProperties;
  }

  // Sender props
  interface SenderProps {
    value: string;
    onChange: (value: string) => void;
    onSend: () => void;
    sending?: boolean;
    placeholder?: string;
    sendText?: string;
    enterToSend?: boolean;
    disabled?: boolean;
    className?: string;
    style?: React.CSSProperties;
  }

  // Attachments props
  interface AttachmentsProps {
    files?: Array<{
      id: string;
      name: string;
      type: string;
      url?: string;
      size?: number;
      [key: string]: any;
    }>;
    onUpload?: (files: File[]) => void;
    onRemove?: (id: string) => void;
    maxSize?: number;
    accept?: string;
    multiple?: boolean;
    disabled?: boolean;
    className?: string;
    style?: React.CSSProperties;
  }

  // ThoughtChain props
  interface Thought {
    type: string;
    tool?: string;
    params?: Record<string, any>;
    result?: any;
    error?: string;
    [key: string]: any;
  }

  interface ThoughtChainProps {
    open: boolean;
    onClose: () => void;
    title?: ReactNode;
    thoughts: Thought[];
    className?: string;
    style?: React.CSSProperties;
  }

  // Prompts props
  interface PromptVariable {
    name: string;
    description?: string;
    type?: string;
    defaultValue?: any;
  }

  interface PromptTemplate {
    id: string;
    title?: string;
    content: string;
    variables?: PromptVariable[];
    description?: string;
    [key: string]: any;
  }

  interface PromptsProps {
    templates: PromptTemplate[];
    activeId?: string;
    onSelect?: (id: string) => void;
    onChange?: (template: PromptTemplate) => void;
    editable?: boolean;
    className?: string;
    style?: React.CSSProperties;
  }

  // Suggestion props
  interface SuggestionItem {
    text: string;
    icon?: ReactNode;
    [key: string]: any;
  }

  interface SuggestionProps {
    data: SuggestionItem[];
    onSelect: (item: SuggestionItem) => void;
    maxLines?: number;
    className?: string;
    style?: React.CSSProperties;
  }

  // Welcome props
  interface WelcomeAction {
    key: string;
    text: string;
    icon?: ReactNode;
    onClick?: () => void;
  }

  interface WelcomeProps {
    title: ReactNode;
    description?: ReactNode;
    avatar?: ReactNode;
    actions?: WelcomeAction[];
    welcomeText?: string;
    className?: string;
    style?: React.CSSProperties;
  }

  // Export components
  export const XProvider: FC<XProviderProps>;
  export const Bubble: FC<BubbleProps>;
  export const Conversations: FC<ConversationsProps>;
  export const Sender: FC<SenderProps>;
  export const Attachments: FC<AttachmentsProps>;
  export const ThoughtChain: FC<ThoughtChainProps>;
  export const Prompts: FC<PromptsProps>;
  export const Suggestion: FC<SuggestionProps>;
  export const Welcome: FC<WelcomeProps>;
  export const XRequest: any;
  export const XStream: any;
} 