import axios, { AxiosResponse, InternalAxiosRequestConfig } from 'axios';

// 创建axios实例
const api = axios.create({
  // 使用相对路径，让Vite代理正确工作
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // 获取token，需要额外检查token格式，确保与后端预期一致
    const token = localStorage.getItem('token');
    console.log('请求拦截器获取token:', token);
    console.log('请求URL:', config.url || '');
    console.log('完整请求URL:', (config.baseURL || '') + (config.url || ''));
    
    if (token) {
      // 确保token格式正确，这里我们直接使用token作为Bearer token
      // 根据后端API的实际要求，可能需要调整
      config.headers.Authorization = `Bearer ${token}`;
      console.log('发送Authorization头:', config.headers.Authorization);
    }
    return config;
  },
  (error) => {
    console.error('请求拦截器错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log('API响应成功:', response.status, response.config.url);
    return response.data;
  },
  (error) => {
    console.error('API响应错误:', error.message);
    
    // 详细记录错误信息以便调试
    if (error.response) {
      console.error('错误状态码:', error.response.status);
      console.error('错误数据:', error.response.data);
      console.error('请求URL:', error.config.url);
      console.error('请求方法:', error.config.method);
      console.error('请求头:', JSON.stringify(error.config.headers));
      
      // 处理401错误，重定向到登录页面
      if (error.response.status === 401) {
        console.log('收到401错误，清除token并重定向');
        localStorage.removeItem('token');
        sessionStorage.removeItem('isLoggedIn');
        window.location.href = '/login';
      }
    } else if (error.request) {
      console.error('请求发送但没有收到响应', error.request);
    }
    
    return Promise.reject(error);
  }
);

export default api; 