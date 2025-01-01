import { LoginReq, RegisterReq, TokenResp, UserInfo, User } from '../types/user';
import { api } from '../utils/api';

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
  }
};