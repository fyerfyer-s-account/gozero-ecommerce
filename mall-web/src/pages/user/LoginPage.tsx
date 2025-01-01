import { Card } from 'antd'
import LoginForm from '@/components/user/LoginForm'
import { Link } from 'react-router-dom'

const LoginPage = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6">Login</h2>
        <LoginForm />
        <div className="text-center mt-4">
          <Link to="/register">Sign Up</Link>
        </div>
      </Card>
    </div>
  )
}

export default LoginPage