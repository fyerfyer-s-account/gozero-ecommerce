import { Card } from 'antd';
import SignUpForm from '@/components/user/SignUpForm';

const SignUpPage = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <Card className="w-full max-w-md">
        <h2 className="text-2xl font-bold text-center mb-6">Sign Up</h2>
        <SignUpForm />
      </Card>
    </div>
  );
};

export default SignUpPage;