import { Form, Input, Select, Button, message } from 'antd';
import { useAuth } from '@/hooks/useAuth';
import { UpdateProfileReq } from '@/types/user';
import { userApi } from '@/api/user';
import * as yup from 'yup';
import { useNavigate } from 'react-router-dom';

const { Option } = Select;

interface ProfileFormData {
  username: string;
  email: string;
  phone?: string; // Make phone optional
}

interface ProfileFormProps {
  onSuccess?: () => void;
}

const schema = yup.object({
  username: yup.string().required('Username is required'),
  email: yup.string().email('Invalid email').required('Email is required'),
  phone: yup.string().optional() // Make phone optional
}).required();

const ProfileForm: React.FC<ProfileFormProps> = ({ onSuccess }) => {
  const { user, refreshProfile } = useAuth();
  const navigate = useNavigate();
  const [form] = Form.useForm();

  const onSubmit = async (values: UpdateProfileReq) => {
    try {
      await userApi.updateProfile({
        ...values,
        gender: values.gender || user?.gender || '', // Ensure gender is never undefined
      });
      await refreshProfile(); // Refresh user data
      message.success('Profile updated successfully');
      onSuccess?.();
      navigate('/profile'); // Redirect to profile page
    } catch (error) {
      console.error('Update profile error:', error);
      message.error('Failed to update profile');
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={{
        nickname: user?.nickname || '',
        email: user?.email || '',          // Ensure email is included
        phone: user?.phone || '',
        gender: user?.gender || ''
      }}
      onFinish={onSubmit}
    >
      <Form.Item name="nickname" label="Nickname">
        <Input />
      </Form.Item>
      <Form.Item name="email" label="Email" rules={[{ type: 'email', message: 'Invalid email' }]}>
        <Input type="email" />
      </Form.Item>
      <Form.Item name="phone" label="Phone">
        <Input />
      </Form.Item>
      <Form.Item name="gender" label="Gender">
        <Select>
          <Option value="male">Male</Option>
          <Option value="female">Female</Option>
          <Option value="other">Other</Option>
        </Select>
      </Form.Item>
      <Button type="primary" htmlType="submit">
        Update Profile
      </Button>
    </Form>
  );
};

export default ProfileForm;