import React, { createContext, useState, useCallback, useEffect } from 'react';
import { Product, SearchReq, Category, CreateProductReq } from '@/types/product';
import { productApi } from '@/api/product';

interface ProductContextType {
  products: Product[];
  currentProduct: Product | null;
  loading: boolean;
  error: string | null;
  searchProducts: (params?: SearchReq) => Promise<void>;
  fetchProductDetails: (id: string) => Promise<void>;
  createProduct: (data: CreateProductReq) => Promise<void>;
  categories: Category[];
  refreshCategories: () => Promise<void>;
}

export const ProductContext = createContext<ProductContextType | null>(null);

export const ProductProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [products, setProducts] = useState<Product[]>([]);
  const [currentProduct, setCurrentProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [categories, setCategories] = useState<Category[]>([]);

  const searchProducts = useCallback(async (params?: SearchReq) => {
    setLoading(true);
    setError(null);
    try {
      const response = await productApi.search(params || {});
      setProducts(response.list);
    } catch (err) {
      setError('Failed to fetch products');
      setProducts([]);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchProductDetails = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const product = await productApi.getProduct(id);
      setCurrentProduct(product);
    } catch (err) {
      setError('Failed to fetch product details');
      setCurrentProduct(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const createProduct = useCallback(async (data: CreateProductReq) => {
    setLoading(true);
    setError(null);
    try {
      await productApi.createProduct(data);
      await searchProducts();  // Refresh product list
    } catch (err) {
      setError('Failed to create product');
      throw err;
    } finally {
      setLoading(false);
    }
  }, [searchProducts]);

  const refreshCategories = useCallback(async () => {
    console.log('Refreshing categories...'); // Debug
    try {
      const data = await productApi.listCategories();
      console.log('Received categories:', data); // Debug
      setCategories(data);
    } catch (err) {
      console.error('Failed to fetch categories:', err);
      setCategories([]);
    }
  }, []);

  useEffect(() => {
    console.log('Initial categories load'); // Debug
    refreshCategories();
  }, []); // Remove refreshCategories from deps

  return (
    <ProductContext.Provider value={{ 
      products, 
      currentProduct,
      loading, 
      error, 
      searchProducts,
      fetchProductDetails,
      createProduct,
      categories,
      refreshCategories
    }}>
      {children}
    </ProductContext.Provider>
  );
};