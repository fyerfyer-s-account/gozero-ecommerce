export interface Product {
  id: number;
  name: string;
  brief: string;
  description: string;
  categoryId: number;
  brand: string;
  images: string[];
  price: number;
  marketPrice: number;
  stock: number;
  sales: number;
  rating: number;
  status: number;
  createdAt: number;
}

export interface Category {
  id: string;
  name: string;
  description: string;
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
}

export interface SearchResp {
  products: Product[];
  total: number;
  page: number;
  totalPages: number;
}