import axios, { AxiosResponse, AxiosError, InternalAxiosRequestConfig } from 'axios';
import useAuthStore from '../store/useAuthStore';

// 创建axios实例
const request = axios.create({
  baseURL: process.env.REACT_APP_API_URL || '/api',
  timeout: 15000,
});

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = useAuthStore.getState().token;
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    return response.data;
  },
  (error: AxiosError) => {
    if (error.response) {
      const { status } = error.response;
      
      // 处理401错误
      if (status === 401) {
        // 清除用户登录状态
        useAuthStore.getState().logout();
        // 跳转到登录页
        window.location.href = '/login';
      }
      
      // 显示错误信息
      const errorMessage = 
        (error.response.data as any)?.message || '服务器错误，请稍后再试';
      console.error(errorMessage);
    } else {
      console.error('网络错误，请检查您的网络连接');
    }
    
    return Promise.reject(error);
  }
);

export default request; 