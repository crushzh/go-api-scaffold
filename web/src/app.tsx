import { history, RequestConfig, RunTimeLayoutConfig } from '@umijs/max';
import { message } from 'antd';
import { TOKEN_KEY, API_PREFIX } from './constants';

// Runtime request configuration
export const request: RequestConfig = {
  baseURL: API_PREFIX,
  timeout: 10000,
  requestInterceptors: [
    (config: any) => {
      const token = localStorage.getItem(TOKEN_KEY);
      if (token) {
        config.headers = {
          ...config.headers,
          Authorization: `Bearer ${token}`,
        };
      }
      return config;
    },
  ],
  responseInterceptors: [
    (response: any) => {
      const { data } = response;
      if (data?.code !== undefined && data.code !== 0) {
        message.error(data.message || 'Request failed');
        return Promise.reject(new Error(data.message));
      }
      return response;
    },
  ],
  errorConfig: {
    errorHandler: (error: any) => {
      if (error?.response?.status === 401) {
        localStorage.removeItem(TOKEN_KEY);
        history.push('/login');
        return;
      }
      message.error(error?.message || 'Network error');
    },
  },
};

// Get initial state (called on app load)
export async function getInitialState(): Promise<{
  currentUser?: API.CurrentUser;
}> {
  const token = localStorage.getItem(TOKEN_KEY);
  if (!token) {
    if (location.pathname !== '/login') {
      history.push('/login');
    }
    return {};
  }
  // Token exists, return basic user info from token
  // In production, replace with API call: GET /api/v1/auth/profile
  try {
    const payload = JSON.parse(atob(token.split('.')[1]));
    return {
      currentUser: {
        username: payload.username || 'admin',
        role: payload.role || 'admin',
      },
    };
  } catch {
    return {};
  }
}

// Layout runtime configuration
export const layout: RunTimeLayoutConfig = ({ initialState }) => {
  return {
    rightContentRender: false,
    waterMarkProps: { content: '' },
    onPageChange: () => {
      const token = localStorage.getItem(TOKEN_KEY);
      if (!token && location.pathname !== '/login') {
        history.push('/login');
      }
    },
    logout: () => {
      localStorage.removeItem(TOKEN_KEY);
      history.push('/login');
    },
    menuHeaderRender: undefined,
  };
};

// Type definitions
declare namespace API {
  interface CurrentUser {
    username: string;
    role: string;
  }
}
