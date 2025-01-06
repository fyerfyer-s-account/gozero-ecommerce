import { 
  Product, 
  SearchReq, 
  SearchResp, 
  CreateProductReq, 
  CreateProductResp,
  Category,
  UpdateProductReq,
  ProductFilter,
  CreateCategoryReq,
  CreateCategoryResp,
  GetCategoriesResp
} from '@/types/product';
import { api } from '@/utils/api';

export const productApi = {
  search: async (params: SearchReq): Promise<SearchResp> => {
    const response = await api.get<SearchResp>('/api/products/search', { params });
    return response.data;
  },

  getProduct: async (id: string): Promise<Product> => {
    const response = await api.get<Product>(`/api/products/${id}`);
    return response.data;
  },

  createProduct: async (data: CreateProductReq): Promise<CreateProductResp> => {
    const response = await api.post<CreateProductResp>('/api/admin/products', data);
    return response.data;
  },

  listCategories: async (): Promise<Category[]> => {
    try {
      const response = await api.get<GetCategoriesResp>('/api/categories');
      console.log('Categories API Response:', response.data);
      return response.data.categories;
    } catch (error) {
      console.error('Categories API Error:', error);
      throw error;
    }
  },

  createCategory: async (data: CreateCategoryReq): Promise<CreateCategoryResp> => {
    const response = await api.post<CreateCategoryResp>('/api/admin/categories', data);
    return response.data;
  },
};