import { useState, useEffect, useCallback } from 'react';
import { 
  getStorageItem, 
  setStorageItem, 
  removeStorageItem, 
  StorageType 
} from '../utils/storageUtils';

/**
 * 本地存储Hook，提供类似于useState的API来操作本地存储
 * @param key 存储键名
 * @param initialValue 初始值，当存储中没有值时使用
 * @param storageType 存储类型，默认为localStorage
 * @param expiration 可选的过期时间(毫秒)
 * @returns [值, 设置值函数, 移除值函数]
 */
export function useLocalStorage<T>(
  key: string,
  initialValue: T,
  storageType: StorageType = StorageType.LOCAL,
  expiration?: number
): [T, (value: T | ((val: T) => T)) => void, () => void] {
  // 从本地存储获取初始值
  const readValue = useCallback((): T => {
    try {
      const storedValue = getStorageItem<T>(key, storageType, expiration);
      return storedValue !== null ? storedValue : initialValue;
    } catch (error) {
      console.error(`Error reading localStorage key "${key}":`, error);
      return initialValue;
    }
  }, [key, initialValue, storageType, expiration]);

  // 状态用于跟踪当前值
  const [storedValue, setStoredValue] = useState<T>(readValue);

  // 返回一个包装版本的setValue函数
  const setValue = useCallback(
    (value: T | ((val: T) => T)) => {
      try {
        // 允许值是一个函数，类似useState
        const valueToStore = value instanceof Function ? value(storedValue) : value;
        
        // 保存到状态
        setStoredValue(valueToStore);
        
        // 保存到本地存储
        setStorageItem(key, valueToStore, storageType);
      } catch (error) {
        console.error(`Error setting localStorage key "${key}":`, error);
      }
    },
    [key, storedValue, storageType]
  );

  // 移除存储项的函数
  const removeValue = useCallback(() => {
    try {
      // 从本地存储移除
      removeStorageItem(key, storageType);
      
      // 重置状态为初始值
      setStoredValue(initialValue);
    } catch (error) {
      console.error(`Error removing localStorage key "${key}":`, error);
    }
  }, [key, initialValue, storageType]);

  // 监听其他窗口对同一存储的更改
  useEffect(() => {
    // 只有在浏览器环境中才能访问window
    if (typeof window === 'undefined') return;

    // 存储事件处理程序
    const handleStorageChange = (event: StorageEvent) => {
      if (event.key && event.key.includes(key) && event.newValue !== null) {
        try {
          const newValue = JSON.parse(event.newValue);
          setStoredValue(newValue.value);
        } catch (e) {
          console.error(`Error parsing storage event for key "${key}":`, e);
        }
      } else if (event.key && event.key.includes(key) && event.newValue === null) {
        setStoredValue(initialValue);
      }
    };

    // 监听存储事件
    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, [key, initialValue]);

  return [storedValue, setValue, removeValue];
}

/**
 * 会话存储Hook，本地存储Hook的简化版本，使用sessionStorage
 * @param key 存储键名
 * @param initialValue 初始值，当存储中没有值时使用
 * @param expiration 可选的过期时间(毫秒)
 * @returns [值, 设置值函数, 移除值函数]
 */
export function useSessionStorage<T>(
  key: string,
  initialValue: T,
  expiration?: number
): [T, (value: T | ((val: T) => T)) => void, () => void] {
  return useLocalStorage<T>(key, initialValue, StorageType.SESSION, expiration);
} 