/**
 * 数据同步工具
 * 处理本地数据与服务器之间的同步
 */
import { getStorageItem, setStorageItem, StorageType } from './storageUtils';

// 同步状态类型
export enum SyncStatus {
  PENDING = 'pending',
  SYNCING = 'syncing',
  SYNCED = 'synced',
  FAILED = 'failed',
}

// 同步队列项类型
export interface SyncQueueItem<T> {
  id: string;
  type: string;
  data: T;
  timestamp: number;
  status: SyncStatus;
  retryCount: number;
}

// 同步队列存储键
const SYNC_QUEUE_KEY = 'sync_queue';

// 最大重试次数
const MAX_RETRY_COUNT = 3;

/**
 * 获取同步队列
 * @returns 同步队列
 */
export const getSyncQueue = <T>(): SyncQueueItem<T>[] => {
  return getStorageItem<SyncQueueItem<T>[]>(SYNC_QUEUE_KEY, StorageType.LOCAL) || [];
};

/**
 * 保存同步队列
 * @param queue 同步队列
 */
export const saveSyncQueue = <T>(queue: SyncQueueItem<T>[]): void => {
  setStorageItem(SYNC_QUEUE_KEY, queue, StorageType.LOCAL);
};

/**
 * 添加项到同步队列
 * @param type 操作类型
 * @param data 数据
 * @param id 可选的ID，如果不提供则使用时间戳
 * @returns 队列项ID
 */
export const addToSyncQueue = <T>(
  type: string,
  data: T,
  id?: string
): string => {
  const queue = getSyncQueue<T>();
  const itemId = id || `${type}_${Date.now()}`;
  
  // 创建新的队列项
  const newItem: SyncQueueItem<T> = {
    id: itemId,
    type,
    data,
    timestamp: Date.now(),
    status: SyncStatus.PENDING,
    retryCount: 0,
  };
  
  // 添加到队列
  queue.push(newItem);
  saveSyncQueue(queue);
  
  return itemId;
};

/**
 * 更新同步队列中的项
 * @param id 队列项ID
 * @param updates 更新内容
 */
export const updateSyncQueueItem = <T>(
  id: string,
  updates: Partial<SyncQueueItem<T>>
): void => {
  const queue = getSyncQueue<T>();
  const index = queue.findIndex(item => item.id === id);
  
  if (index !== -1) {
    queue[index] = { ...queue[index], ...updates };
    saveSyncQueue(queue);
  }
};

/**
 * 移除同步队列中的项
 * @param id 队列项ID
 */
export const removeSyncQueueItem = (id: string): void => {
  const queue = getSyncQueue();
  const updatedQueue = queue.filter(item => item.id !== id);
  saveSyncQueue(updatedQueue);
};

/**
 * 同步单个队列项
 * @param item 队列项
 * @param syncFunction 同步函数，接收数据并返回Promise
 * @returns Promise<boolean> 是否同步成功
 */
export const syncQueueItem = async <T>(
  item: SyncQueueItem<T>,
  syncFunction: (type: string, data: T) => Promise<any>
): Promise<boolean> => {
  try {
    // 更新状态为同步中
    updateSyncQueueItem(item.id, { status: SyncStatus.SYNCING });
    
    // 调用同步函数
    await syncFunction(item.type, item.data);
    
    // 同步成功，从队列中移除
    removeSyncQueueItem(item.id);
    return true;
  } catch (error) {
    console.error(`Sync error for item ${item.id}:`, error);
    
    // 增加重试计数
    const newRetryCount = item.retryCount + 1;
    
    // 如果超过最大重试次数，标记为失败
    if (newRetryCount > MAX_RETRY_COUNT) {
      updateSyncQueueItem(item.id, { 
        status: SyncStatus.FAILED,
        retryCount: newRetryCount
      });
    } else {
      // 否则标记为等待重试
      updateSyncQueueItem(item.id, { 
        status: SyncStatus.PENDING,
        retryCount: newRetryCount
      });
    }
    
    return false;
  }
};

/**
 * 同步所有等待同步的队列项
 * @param syncFunction 同步函数，接收类型和数据并返回Promise
 * @returns Promise<{success: number, failed: number}> 同步结果统计
 */
export const syncAllPendingItems = async <T>(
  syncFunction: (type: string, data: T) => Promise<any>
): Promise<{success: number, failed: number}> => {
  const queue = getSyncQueue<T>();
  const pendingItems = queue.filter(item => 
    item.status === SyncStatus.PENDING || 
    (item.status === SyncStatus.FAILED && item.retryCount < MAX_RETRY_COUNT)
  );
  
  let success = 0;
  let failed = 0;
  
  // 依次同步每个项
  for (const item of pendingItems) {
    const result = await syncQueueItem(item, syncFunction);
    if (result) {
      success++;
    } else {
      failed++;
    }
  }
  
  return { success, failed };
}; 