import { defineConfig } from '@umijs/max';

export default defineConfig({
  antd: {},
  access: {},
  model: {},
  initialState: {},
  request: {},
  layout: {
    title: 'My Service',
    locale: false,
  },
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
    '/swagger': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
  },
  routes: [
    {
      path: '/login',
      component: './Login',
      layout: false,
    },
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      name: 'Dashboard',
      path: '/dashboard',
      component: './Dashboard',
      icon: 'DashboardOutlined',
    },
    {
      name: 'Example',
      path: '/example',
      component: './Example',
      icon: 'TableOutlined',
    },
  ],
  outputPath: '../internal/web/dist',
  npmClient: 'npm',
  hash: true,
  jsMinifier: 'terser',
});
