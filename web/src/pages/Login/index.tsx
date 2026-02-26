import { LoginForm, ProFormText } from '@ant-design/pro-components';
import { history, useModel } from '@umijs/max';
import { message } from 'antd';
import { LockOutlined, UserOutlined } from '@ant-design/icons';
import { login } from '@/services/auth';
import { TOKEN_KEY } from '@/constants';

const LoginPage: React.FC = () => {
  const { refresh } = useModel('@@initialState');

  const handleLogin = async (values: { username: string; password: string }) => {
    try {
      const res = await login(values);
      if (res?.data?.token) {
        localStorage.setItem(TOKEN_KEY, res.data.token);
        message.success('Login successful');
        await refresh();
        history.push('/dashboard');
      }
    } catch (error) {
      // Error handled by request interceptor
    }
  };

  return (
    <div style={{ height: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f0f2f5' }}>
      <LoginForm
        title="My Service"
        subTitle="Backend Management System"
        onFinish={handleLogin}
      >
        <ProFormText
          name="username"
          fieldProps={{ size: 'large', prefix: <UserOutlined /> }}
          placeholder="Username"
          rules={[{ required: true, message: 'Please enter username' }]}
        />
        <ProFormText.Password
          name="password"
          fieldProps={{ size: 'large', prefix: <LockOutlined /> }}
          placeholder="Password"
          rules={[{ required: true, message: 'Please enter password' }]}
        />
      </LoginForm>
    </div>
  );
};

export default LoginPage;
