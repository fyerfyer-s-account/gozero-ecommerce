import { api } from '@/utils/api';
import { RegisterReq, LoginReq, TokenResp, UserInfo, UpdateProfileReq, User } from '@/types/user';

export const userApi = {
  register: async (data: RegisterReq): Promise<User> => {
    const response = await api.post<User>('/api/user/register', data);
    return response.data;
  },

  login: async (data: LoginReq): Promise<TokenResp> => {
    const response = await api.post<TokenResp>('/api/user/login', data);
    return response.data;
  },

  getProfile: async (): Promise<UserInfo> => {
    const response = await api.get<UserInfo>('/api/user/profile');
    return response.data;
  },

  updateProfile: async (data: UpdateProfileReq): Promise<UserInfo> => {
    const response = await api.put<UserInfo>('/api/user/profile', data);
    return response.data;
  }
};