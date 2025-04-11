import React, { ReactNode } from 'react';
import { XProvider } from '@ant-design/x';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';

interface XConfigProviderProps {
  children: ReactNode;
}

/**
 * 全局的Ant Design X配置提供器
 * 用于提供统一的主题和配置
 */
const XConfigProvider: React.FC<XConfigProviderProps> = ({ children }) => {
  return (
    <ConfigProvider locale={zhCN}>
      <XProvider
        theme={{
          primaryColor: '#00b96b',
          borderRadius: 8,
        }}
        settings={{
          // API相关设置
          api: {
            baseURL: '/api/v1',
            // 可以配置头部信息等
            headers: {
              'Content-Type': 'application/json',
            },
          },
          // UI相关设置
          ui: {
            avatarShape: 'circle',
            darkMode: false,
          },
        }}
      >
        {children}
      </XProvider>
    </ConfigProvider>
  );
};

export default XConfigProvider; 