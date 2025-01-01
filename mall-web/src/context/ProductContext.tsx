import React, { createContext, useState, useCallback } from 'react';
import { Product, SearchReq } from '../types/product';
import { productApi } from '../api/product';

interface ProductContextType {
  products: Product[];
  loading: boolean;
  error: string | null;
  fetchProducts: () => Promise<void>;
  searchProducts: (params: SearchReq) => Promise<void>;
  currentProduct: Product | null;
  fetchProductDetails: (id: string) => Promise<void>;
}

export const ProductContext = createContext<ProductContextType>({} as ProductContextType);

export const ProductProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [products, setProducts] = useState<Product[]>([]);  // Initialize as empty array
  const [currentProduct, setCurrentProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true);
      const data = await productApi.getProducts();
      setProducts(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch products');
    } finally {
      setLoading(false);
    }
  }, []);

  const searchProducts = useCallback(async (params: SearchReq) => {
    try {
      setLoading(true);
      const data = await productApi.search(params);
      setProducts(data.list); // Changed from data.products to data.list to match API
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to search products');
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchProductDetails = useCallback(async (id: string) => {
    try {
      setLoading(true);
      const data = await productApi.getProduct(id);
      setCurrentProduct(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch product details');
    } finally {
      setLoading(false);
    }
  }, []);

  return (
    <ProductContext.Provider 
      value={{ 
        products, 
        loading, 
        error, 
        fetchProducts, 
        searchProducts,
        currentProduct,
        fetchProductDetails
      }}
    >
      {children}
    </ProductContext.Provider>
  );
};