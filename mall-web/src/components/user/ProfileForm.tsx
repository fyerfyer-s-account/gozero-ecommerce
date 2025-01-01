import { Form, Input, Select, Button, message } from 'antd';
import { useAuth } from '@/hooks/useAuth';
import { UpdateProfileReq } from '@/types/user';
import { userApi } from '@/api/user';
import * as yup from 'yup';

const { Option } = Select;

interface ProfileFormData {
  username: string;
  email: string;
  phone?: string; // Make phone optional
}

const schema = yup.object({
  username: yup.string().required('Username is required'),
  email: yup.string().email('Invalid email').required('Email is required'),
  phone: yup.string().optional() // Make phone optional
}).required();

const ProfileForm = () => {
  const { user } = useAuth();
  const [form] = Form.useForm();

  const onSubmit = async (values: UpdateProfileReq) => {
    try {
      await userApi.updateProfile(values);
      message.success('Profile updated successfully');
    } catch (error) {
      message.error('Failed to update profile');
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={user || {}}
      onFinish={onSubmit}
    >
      <Form.Item name="nickname" label="Nickname">
        <Input />
      </Form.Item>
      <Form.Item name="email" label="Email">
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