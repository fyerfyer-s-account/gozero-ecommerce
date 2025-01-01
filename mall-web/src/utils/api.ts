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
});

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    const message = error.response?.data?.message || 'An error occurred';
    return Promise.reject(new Error(message));
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