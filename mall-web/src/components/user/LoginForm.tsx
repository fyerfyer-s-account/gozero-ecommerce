import { Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { LoginReq } from '@/types/user';

const LoginForm = () => {
  const navigate = useNavigate();
  const { login, loading } = useAuth();
  const [form] = Form.useForm();

  const onSubmit = async (values: LoginReq) => {
    try {
      const response = await login(values.username, values.password);
      message.success('Login successful');
      navigate('/', { replace: true });  // Remove setTimeout, direct navigation
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Login failed');
    }
  };

  return (
    <Form 
      form={form}
      layout="vertical" 
      onFinish={onSubmit}
      autoComplete="off"
    >
      <Form.Item 
        name="username"
        label="Username"
        rules={[{ required: true, message: 'Username is required' }]}
      >
        <Input />
      </Form.Item>
      
      <Form.Item 
        name="password"
        label="Password"
        rules={[{ required: true, message: 'Password is required' }]}
      >
        <Input.Password />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" loading={loading} block>
          Login
        </Button>
      </Form.Item>
    </Form>
  );
};

export default LoginForm;