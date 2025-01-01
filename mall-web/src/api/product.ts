import { api } from '../utils/api';
import { Product, SearchReq, SearchResp } from '../types/product';

export const productApi = {
  search: async (params: SearchReq): Promise<SearchResp> => {
    const response = await api.get<SearchResp>('/products/search', { params });
    return response.data;
  },

  getProduct: async (id: string): Promise<Product> => {
    const response = await api.get<Product>(`/products/${id}`);
    return response.data;
  },

  getProducts: async (): Promise<Product[]> => {
    const response = await api.get<Product[]>('/products');
    return response.data;
  }
};