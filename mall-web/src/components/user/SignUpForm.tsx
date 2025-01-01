import { Form, Input, Button, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { useForm, Controller } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { RegisterReq } from '@/types/user';

const schema = yup.object().shape({
  username: yup.string().required('Username is required'),
  password: yup.string().min(6, 'Password must be at least 6 characters').required('Password is required'),
  email: yup.string().email('Invalid email').optional(),
  phone: yup.string().matches(/^1[3-9]\d{9}$/, 'Invalid phone number').optional(),
});

const SignUpForm = () => {
  const navigate = useNavigate();
  const { register: signUp, loading } = useAuth();
  const { control, handleSubmit, formState: { errors } } = useForm<RegisterReq>({
    resolver: yupResolver(schema),
    mode: 'onBlur'
  });

  const onSubmit = async (data: RegisterReq) => {
    try {
      await signUp(data);
      message.success('Registration successful!');
      navigate('/login');
    } catch (error) {
      message.error(error instanceof Error ? error.message : 'Registration failed');
    }
  };

  return (
    <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
      <Form.Item label="Username" validateStatus={errors.username ? 'error' : ''} help={errors.username?.message}>
        <Controller
          name="username"
          control={control}
          render={({ field }) => <Input {...field} />}
        />
      </Form.Item>
      <Form.Item label="Password" validateStatus={errors.password ? 'error' : ''} help={errors.password?.message}>
        <Controller
          name="password"
          control={control}
          render={({ field }) => <Input.Password {...field} />}
        />
      </Form.Item>
      <Form.Item label="Email" validateStatus={errors.email ? 'error' : ''} help={errors.email?.message}>
        <Controller
          name="email"
          control={control}
          render={({ field }) => <Input {...field} />}
        />
      </Form.Item>
      <Form.Item label="Phone" validateStatus={errors.phone ? 'error' : ''} help={errors.phone?.message}>
        <Controller
          name="phone"
          control={control}
          render={({ field }) => <Input {...field} />}
        />
      </Form.Item>
      <Button type="primary" htmlType="submit" loading={loading} block>
        Sign Up
      </Button>
    </Form>
  );
};

export default SignUpForm;