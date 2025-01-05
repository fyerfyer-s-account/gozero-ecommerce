import { Form, Input, InputNumber, Button, message, Select } from 'antd';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useProduct } from '@/hooks/useProduct';
import SkuForm from './SkuForm';
import { CreateProductReq, SkuAttributeReq, Category, ProductStatus } from '@/types/product';
import goodsIcon from '@/icon/goods_icon.jpg';

const { Option } = Select;

interface ProductFormData {
  name: string;
  description: string;
  categoryId: number;
  brand: string;
  price: number;
  skuAttributes: SkuAttributeReq[];
}

const ProductForm = () => {
  const { createProduct, loading, categories, refreshCategories } = useProduct();
  const navigate = useNavigate();
  const [form] = Form.useForm();

  useEffect(() => {
    const loadCategories = async () => {
      try {
        await refreshCategories();
      } catch (error) {
        console.error('Failed to load categories:', error);
        message.error('Failed to load categories');
      }
    };
    loadCategories();
  }, []); // Only load once on mount

  useEffect(() => {
    console.log('Current categories:', categories); // Debug
  }, [categories]);

  const onSubmit = async (values: ProductFormData) => {
    try {
      if (!values.categoryId) {
        message.error('Please select a category');
        return;
      }

      const productData: CreateProductReq = {
        name: values.name,
        description: values.description,
        categoryId: values.categoryId,
        brand: values.brand,
        price: values.price,
        images: [goodsIcon],
        skuAttributes: values.skuAttributes || []
      };

      await createProduct(productData);
      message.success('Product created successfully');
      navigate('/');
    } catch (error) {
      console.error('Create product error:', error);
      message.error('Failed to create product');
    }
  };

  return (
    <Form form={form} layout="vertical" onFinish={onSubmit}>
      <Form.Item name="name" label="Product Name" rules={[{ required: true }]}>
        <Input />
      </Form.Item>

      <Form.Item name="categoryId" label="Category" rules={[{ required: true }]}>
        <Select placeholder="Select a category">
          {categories?.length > 0 ? (
            categories.map(category => (
              <Option key={category.id} value={category.id}>
                {category.name} (ID: {category.id})
              </Option>
            ))
          ) : (
            <Option disabled>No categories available</Option>
          )}
        </Select>
      </Form.Item>

      <Form.Item name="description" label="Description" rules={[{ required: true }]}>
        <Input.TextArea />
      </Form.Item>

      <Form.Item name="brand" label="Brand" rules={[{ required: true }]}>
        <Input />
      </Form.Item>

      <Form.Item name="price" label="Price" rules={[{ required: true }]}>
        <InputNumber 
          min={0} 
          precision={2} 
          style={{ width: '200px' }}
        />
      </Form.Item>

      <Form.Item label="SKU Attributes">
        <SkuForm onChange={attributes => form.setFieldsValue({ skuAttributes: attributes })} />
      </Form.Item>

      <Button type="primary" htmlType="submit" loading={loading}>
        Create Product
      </Button>
    </Form>
  );
};

export default ProductForm;