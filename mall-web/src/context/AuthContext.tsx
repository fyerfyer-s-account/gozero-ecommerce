import React, { createContext, useState, useCallback, useContext, useEffect } from 'react';
import { userApi } from '@/api/user';
import { UserInfo, RegisterReq, AuthState, LoginReq, TokenResp } from '@/types/user';

interface AuthContextType {
  user: UserInfo | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
  login: (username: string, password: string) => Promise<TokenResp>;
  register: (data: RegisterReq) => Promise<void>;
  logout: () => void;
  refreshProfile: () => Promise<void>; // Add this
}

export const AuthContext = createContext<AuthContextType>({} as AuthContextType);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<UserInfo | null>(null);
  const [token, setToken] = useState<string | null>(() => localStorage.getItem('token'));
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => !!localStorage.getItem('token'));
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const loadUserProfile = async () => {
      if (token && !user) {
        try {
          const userInfo = await userApi.getProfile();
          const userWithRole = {
            ...userInfo,
            role: localStorage.getItem('role') || 'user'
          };
          setUser(userWithRole);
        } catch (error) {
          logout();
        }
      }
    };
    loadUserProfile();
  }, [token]);

  const login = useCallback(async (username: string, password: string) => {
    setLoading(true);
    try {
      const response = await userApi.login({ username, password });
      localStorage.setItem('token', response.accessToken);
      setToken(response.accessToken);
      setIsAuthenticated(true);
      
      let role: string | undefined;
      try {
        const payload = JSON.parse(atob(response.accessToken.split('.')[1]));
        role = payload.role || undefined;
        if (role) {
          localStorage.setItem('role', role);
        }
      } catch (e) {
        console.error('Error parsing JWT:', e);
      }
      
      const userInfo = await userApi.getProfile();
      setUser({
        ...userInfo,
        role: role
      });
      
      return response;
    } catch (error) {
      logout();
      throw error;
    } finally {
      setLoading(false);
    }
  }, []);

  const register = useCallback(async (data: RegisterReq) => {
    setLoading(true);
    try {
      await userApi.register(data);
    } finally {
      setLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    try {
      if (token) {
        await userApi.logout();
      }
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      localStorage.removeItem('token');
      localStorage.removeItem('role'); // Add this
      setToken(null);
      setUser(null);
      setIsAuthenticated(false);
    }
  }, [token]);

  const refreshProfile = useCallback(async () => {
    if (token) {
      try {
        const userInfo = await userApi.getProfile();
        setUser(userInfo);
      } catch (error) {
        console.error('Failed to refresh profile:', error);
      }
    }
  }, [token]);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isAuthenticated,
        loading,
        login,
        register,
        logout,
        refreshProfile // Add this
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
