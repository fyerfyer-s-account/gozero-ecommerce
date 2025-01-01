import { Form, Input, Button } from 'antd'
import { useForm } from 'react-hook-form'
import { yupResolver } from '@hookform/resolvers/yup'
import * as yup from 'yup'

interface ProfileFormData {
  username: string;
  email: string;
  phone?: string; // Make phone optional
}

const schema = yup.object({
  username: yup.string().required('Username is required'),
  email: yup.string().email('Invalid email').required('Email is required'),
  phone: yup.string().optional() // Make phone optional
}).required()

const ProfileForm = () => {
  const { register, handleSubmit, formState: { errors } } = useForm<ProfileFormData>({
    resolver: yupResolver(schema)
  })

  const onSubmit = (data: ProfileFormData) => {
    console.log(data)
  }

  return (
    <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
      <Form.Item label="Username" validateStatus={errors.username ? 'error' : ''}>
        <Input {...register('username')} />
      </Form.Item>
      <Form.Item label="Email" validateStatus={errors.email ? 'error' : ''}>
        <Input {...register('email')} type="email" />
      </Form.Item>
      <Form.Item label="Phone" validateStatus={errors.phone ? 'error' : ''}>
        <Input {...register('phone')} />
      </Form.Item>
      <Button type="primary" htmlType="submit">
        Update Profile
      </Button>
    </Form>
  )
}

export default ProfileForm