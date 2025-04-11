/**
 * 本地存储工具
 * 提供对localStorage和sessionStorage的封装
 */

// 缓存前缀，防止命名冲突
const STORAGE_PREFIX = 'agent_platform_';

/**
 * 本地存储类型
 */
export enum StorageType {
  LOCAL = 'localStorage',
  SESSION = 'sessionStorage',
}

/**
 * 设置缓存
 * @param key 缓存键名
 * @param value 缓存值
 * @param type 存储类型，默认为localStorage
 */
export const setStorageItem = <T>(
  key: string,
  value: T,
  type: StorageType = StorageType.LOCAL
): void => {
  try {
    const prefixedKey = `${STORAGE_PREFIX}${key}`;
    const storageValue = JSON.stringify({
      value,
      timestamp: new Date().getTime(),
    });
    window[type].setItem(prefixedKey, storageValue);
  } catch (error) {
    console.error('Error setting storage item:', error);
  }
};

/**
 * 获取缓存
 * @param key 缓存键名
 * @param type 存储类型，默认为localStorage
 * @param expiration 过期时间(毫秒)，默认不过期
 * @returns 缓存值或null(如果不存在或已过期)
 */
export const getStorageItem = <T>(
  key: string,
  type: StorageType = StorageType.LOCAL,
  expiration?: number
): T | null => {
  try {
    const prefixedKey = `${STORAGE_PREFIX}${key}`;
    const item = window[type].getItem(prefixedKey);

    if (!item) return null;

    const { value, timestamp } = JSON.parse(item);

    // 检查是否过期
    if (expiration && new Date().getTime() - timestamp > expiration) {
      removeStorageItem(key, type);
      return null;
    }

    return value as T;
  } catch (error) {
    console.error('Error getting storage item:', error);
    return null;
  }
};

/**
 * 移除缓存
 * @param key 缓存键名
 * @param type 存储类型，默认为localStorage
 */
export const removeStorageItem = (
  key: string,
  type: StorageType = StorageType.LOCAL
): void => {
  try {
    const prefixedKey = `${STORAGE_PREFIX}${key}`;
    window[type].removeItem(prefixedKey);
  } catch (error) {
    console.error('Error removing storage item:', error);
  }
};

/**
 * 清除所有缓存
 * @param type 存储类型，默认为localStorage
 */
export const clearStorage = (type: StorageType = StorageType.LOCAL): void => {
  try {
    const storage = window[type];
    const keys = Object.keys(storage);
    
    // 只清除带有前缀的项
    keys.forEach((key) => {
      if (key.startsWith(STORAGE_PREFIX)) {
        storage.removeItem(key);
      }
    });
  } catch (error) {
    console.error('Error clearing storage:', error);
  }
};

/**
 * 获取所有缓存的键
 * @param type 存储类型，默认为localStorage
 * @returns 所有缓存的键(不含前缀)
 */
export const getStorageKeys = (
  type: StorageType = StorageType.LOCAL
): string[] => {
  try {
    const storage = window[type];
    const keys = Object.keys(storage);
    
    return keys
      .filter((key) => key.startsWith(STORAGE_PREFIX))
      .map((key) => key.replace(STORAGE_PREFIX, ''));
  } catch (error) {
    console.error('Error getting storage keys:', error);
    return [];
  }
}; 