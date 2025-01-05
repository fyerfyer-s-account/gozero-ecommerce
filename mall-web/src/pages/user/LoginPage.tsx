import { Card } from 'antd'
import LoginForm from '@/components/user/LoginForm'
import { Link } from 'react-router-dom'
import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

const LoginPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (isAuthenticated) {
      navigate(location.state?.from?.pathname || '/', { replace: true });
    }
  }, [isAuthenticated, navigate, location]);

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