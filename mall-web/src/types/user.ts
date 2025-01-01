export interface User {
  id: number;
  username: string;
  email: string;
  phone?: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
}

export interface LoginReq {
  username: string;
  password: string;
}

export interface TokenResp {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  phone: string;
  email: string;
  gender: string;
  memberLevel: number;
  balance: number;
  createdAt: number;
}

export interface RegisterReq {
  username: string;
  password: string;
  phone?: string;
  email?: string;
}

export interface RegisterResp {
  userId: number;
}

export interface UpdateProfileReq {
  nickname?: string;
  avatar?: string;
  gender: string;
  phone?: string;
  email?: string;
}