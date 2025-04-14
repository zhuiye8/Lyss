import React, { ReactNode } from 'react';
import { ConfigProvider, theme } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';

interface XConfigProviderProps {
  children: ReactNode;
}

/**
 * 全局的配置提供器
 * 用于提供统一的主题和配置
 */
const XConfigProvider: React.FC<XConfigProviderProps> = ({ children }) => {
  return (
    <ConfigProvider 
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: '#00b96b',
          borderRadius: 8,
        },
        algorithm: theme.defaultAlgorithm
      }}
    >
      {children}
    </ConfigProvider>
  );
};

export default XConfigProvider; 