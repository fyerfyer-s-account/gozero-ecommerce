Sure, here's the contents for the file `/mall-web/mall-web/src/api/user.ts`:

import { User, RegisterReq, LoginReq, UpdateProfileReq } from '../types/user';
import { api } from '../utils/api';

export const registerUser = async (data: RegisterReq): Promise<User> => {
    const response = await api.post('/api/user/register', data);
    return response.data;
};

export const loginUser = async (data: LoginReq): Promise<string> => {
    const response = await api.post('/api/user/login', data);
    return response.data.token;
};

export const getUserProfile = async (): Promise<User> => {
    const response = await api.get('/api/user/profile');
    return response.data;
};

export const updateUserProfile = async (data: UpdateProfileReq): Promise<void> => {
    await api.put('/api/user/profile', data);
};