/**
 * 会话存储工具
 * 用于在localStorage中缓存会话数据
 */

const PREFIX = 'lyss_app_';

interface StorageItem<T> {
  data: T;
  timestamp: number;
  expiry?: number; // 过期时间（毫秒）
}

/**
 * 保存数据到本地存储
 * @param key 存储键名
 * @param data 要存储的数据
 * @param expiryMs 过期时间（毫秒），如果不设置则不过期
 */
export const saveToStorage = <T>(key: string, data: T, expiryMs?: number): void => {
  try {
    const storageKey = `${PREFIX}${key}`;
    const item: StorageItem<T> = {
      data,
      timestamp: Date.now(),
      expiry: expiryMs ? Date.now() + expiryMs : undefined
    };
    localStorage.setItem(storageKey, JSON.stringify(item));
  } catch (error) {
    console.error('Failed to save to localStorage:', error);
  }
};

/**
 * 从本地存储中获取数据
 * @param key 存储键名
 * @returns 存储的数据，如果不存在或已过期则返回null
 */
export const getFromStorage = <T>(key: string): T | null => {
  try {
    const storageKey = `${PREFIX}${key}`;
    const storedValue = localStorage.getItem(storageKey);
    
    if (!storedValue) {
      return null;
    }
    
    const item: StorageItem<T> = JSON.parse(storedValue);
    
    // 检查是否过期
    if (item.expiry && item.expiry < Date.now()) {
      localStorage.removeItem(storageKey);
      return null;
    }
    
    return item.data;
  } catch (error) {
    console.error('Failed to get from localStorage:', error);
    return null;
  }
};

/**
 * 从本地存储中删除数据
 * @param key 存储键名
 */
export const removeFromStorage = (key: string): void => {
  try {
    const storageKey = `${PREFIX}${key}`;
    localStorage.removeItem(storageKey);
  } catch (error) {
    console.error('Failed to remove from localStorage:', error);
  }
};

/**
 * 清除与应用相关的所有本地存储
 */
export const clearAllStorage = (): void => {
  try {
    Object.keys(localStorage).forEach(key => {
      if (key.startsWith(PREFIX)) {
        localStorage.removeItem(key);
      }
    });
  } catch (error) {
    console.error('Failed to clear all storage:', error);
  }
};

/**
 * 会话专用存储 - 专门用于存储对话数据
 */
export const conversationStorage = {
  /**
   * 保存对话数据
   * @param conversationId 对话ID
   * @param data 对话数据
   */
  saveConversation: <T>(conversationId: string, data: T): void => {
    saveToStorage(`conversation_${conversationId}`, data);
  },
  
  /**
   * 获取对话数据
   * @param conversationId 对话ID
   * @returns 对话数据
   */
  getConversation: <T>(conversationId: string): T | null => {
    return getFromStorage<T>(`conversation_${conversationId}`);
  },
  
  /**
   * 删除对话数据
   * @param conversationId 对话ID
   */
  removeConversation: (conversationId: string): void => {
    removeFromStorage(`conversation_${conversationId}`);
  },
  
  /**
   * 获取所有对话ID列表
   * @returns 对话ID列表
   */
  getAllConversationIds: (): string[] => {
    const ids: string[] = [];
    try {
      Object.keys(localStorage).forEach(key => {
        if (key.startsWith(`${PREFIX}conversation_`)) {
          const id = key.replace(`${PREFIX}conversation_`, '');
          ids.push(id);
        }
      });
    } catch (error) {
      console.error('Failed to get all conversation IDs:', error);
    }
    return ids;
  }
}; 