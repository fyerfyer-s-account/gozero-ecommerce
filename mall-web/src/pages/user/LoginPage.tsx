import { Card } from 'antd'
import LoginForm from '@/components/user/LoginForm'

const LoginPage = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6">Login</h2>
        <LoginForm />
      </Card>
    </div>
  )
}

export default LoginPage