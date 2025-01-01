import { Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

const schema = yup.object({
  username: yup.string().required('Username is required'),
  password: yup.string().required('Password is required').min(6),
  email: yup.string().email('Invalid email').required('Email is required'),
  phone: yup.string(),
}).required();

const SignUpForm = () => {
  const navigate = useNavigate();
  const { register: signUp, loading } = useAuth();
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: yupResolver(schema)
  });

  const onSubmit = async (data: any) => {
    try {
      await signUp(data);
      message.success('Sign up successful');
      navigate('/login');
    } catch (error) {
      message.error('Sign up failed');
    }
  };

  return (
    <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
      <Form.Item
        label="Username"
        validateStatus={errors.username ? 'error' : ''}
        help={errors.username?.message}
      >
        <Input {...register('username')} />
      </Form.Item>

      <Form.Item
        label="Password"
        validateStatus={errors.password ? 'error' : ''}
        help={errors.password?.message}
      >
        <Input.Password {...register('password')} />
      </Form.Item>

      <Form.Item
        label="Email"
        validateStatus={errors.email ? 'error' : ''}
        help={errors.email?.message}
      >
        <Input {...register('email')} />
      </Form.Item>

      <Form.Item
        label="Phone (Optional)"
        validateStatus={errors.phone ? 'error' : ''}
        help={errors.phone?.message}
      >
        <Input {...register('phone')} />
      </Form.Item>

      <Button type="primary" htmlType="submit" loading={loading} block>
        Sign Up
      </Button>
    </Form>
  );
};

export default SignUpForm;