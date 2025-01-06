export interface Product {
  id: number;
  name: string;
  brief: string;
  description: string;
  categoryId: number;
  categoryName: string;
  brand: string;
  images: string[];
  price: number;
  stock: number;
  sales: number;
  rating: number;
  status: number;
  createdAt: number;
  skus?: Sku[];
}

export interface CreateProductReq {
  name: string;
  description: string;
  categoryId: number;
  brand: string;
  images: string[];
  price: number;
  skuAttributes: SkuAttributeReq[];
}

export interface CreateProductResp {
  id: number;
}

export interface SkuAttributeReq {
  key: string;
  value: string;
}

export interface Category {
  id: number;
  name: string;
  parentId: number;
  level: number;
  sort: number;
  icon?: string;
  children?: Category[];
}

export interface Sku {
  id: number;
  productId: number;
  name: string;
  code: string;
  price: number;
  stock: number;
  attributes: Record<string, string>;
  status?: number;
}

export interface SearchReq {
  keyword?: string;
  categoryId?: number;
  brandId?: number;
  priceMin?: number;
  priceMax?: number;
  attributes?: string[];
  sort?: string;
  order?: 'asc' | 'desc';
  page?: number;
  pageSize?: number;
}

export interface SearchResp {
  list: Product[];       
  total: number;
  page: number;
  totalPages: number;
}

export interface ProductFilter {
  categories: Category[];
  brands: string[];
  priceRange: {
    min: number;
    max: number;
  };
  attributes: {
    [key: string]: string[];
  };
}

export type ProductStatus = 'active' | 'inactive' | 'deleted';

export interface UpdateProductReq extends Partial<CreateProductReq> {
  id: number;
  status?: ProductStatus;
}

export interface CreateCategoryReq {
  name: string;
  parentId: number;  // Changed from optional
  sort: number;
  icon?: string;
}

export interface CreateCategoryResp {
  id: number;
}

export interface GetCategoriesResp {
  categories: Category[];
}