import { Product, SearchReq, SearchResp, GetProductReq } from '../types/product';
import { api } from '../utils/api';

export const fetchProducts = async (params: SearchReq): Promise<SearchResp> => {
    const response = await api.get('/api/products/search', { params });
    return response.data;
};

export const fetchProductById = async (id: string): Promise<Product> => {
    const response = await api.get(`/api/products/${id}`);
    return response.data;
};

export const fetchProductSkus = async (id: string): Promise<Sku[]> => {
    const response = await api.get(`/api/products/${id}/skus`);
    return response.data;
};