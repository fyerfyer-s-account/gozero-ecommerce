import { useContext, useState } from 'react';
import { AuthContext } from '../context/AuthContext';
import { userApi } from '@/api/user';
import { RegisterReq } from '@/types/user';

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }

  const [loading, setLoading] = useState(false);

  const register = async (data: RegisterReq) => {
    setLoading(true);
    try {
      const response = await userApi.register(data);
      return response;
    } finally {
      setLoading(false);
    }
  };

  return {
    ...context,
    register,
    loading,
  };
};