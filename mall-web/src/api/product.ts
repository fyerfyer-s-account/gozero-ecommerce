import { Product, SearchReq, SearchResp } from '../types/product';

const BASE_URL = '/api/products';

export const productApi = {
  search: async (params: SearchReq): Promise<SearchResp> => {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined) {
        queryParams.append(key, value.toString());
      }
    });
    
    const response = await fetch(`${BASE_URL}/search?${queryParams.toString()}`);
    return response.json();
  },

  getProduct: async (id: string): Promise<Product> => {
    const response = await fetch(`${BASE_URL}/${id}`);
    return response.json();
  },

  getProducts: async (): Promise<Product[]> => {
    const response = await fetch(BASE_URL);
    return response.json();
  }
};