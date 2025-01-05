import { useEffect, useState } from 'react';
import axios from 'axios';

const API_BASE_URL = 'http://0.0.0.0:9000'; // Update with your backend URL

interface ApiResponse<T> {
  data: T;
  message?: string;
  code?: number;
}

export const api = axios.create({
  baseURL: '/api', // Use relative path for proxy. DON'T REMOVE THIS!!! OR YOU WILL GET CORS ERROR!!!
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 5000,
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      // Clear role if not admin
      if (payload.role !== 'admin') {
        localStorage.removeItem('role');
      }
    } catch (e) {
      console.error('Error parsing JWT:', e);
    }
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    console.group('API Error Details');
    console.log('Status:', error.response?.status);
    console.log('URL:', error.config?.url);
    console.log('Method:', error.config?.method);
    console.log('Headers:', error.config?.headers);
    console.log('Response:', error.response?.data);
    console.groupEnd();
    
    if (error.response?.status === 401) {
      if (!error.config?.url.includes('/api/user/')) {
        console.log('Non-auth 401 error - keeping token');
        return Promise.reject(error);
      }
      localStorage.removeItem('token');
      localStorage.removeItem('role');
    }
    return Promise.reject(error);
  }
);

export const fetchData = async (endpoint: string, options?: RequestInit) => {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
  if (!response.ok) {
    throw new Error('Network response was not ok');
  }
  return response.json();
};

export const useFetch = (endpoint: string, options?: RequestInit) => {
  const [data, setData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    const fetchDataAsync = async () => {
      try {
        const result = await fetchData(endpoint, options);
        setData(result);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('Unknown error'));
      } finally {
        setLoading(false);
      }
    };

    fetchDataAsync();
  }, [endpoint, options]);

  return { data, loading, error };
};