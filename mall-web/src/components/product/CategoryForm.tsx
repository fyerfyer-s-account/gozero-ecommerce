import { Form, Input, InputNumber, Button, message, Select } from 'antd';
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Category, CreateCategoryReq } from '@/types/product';
import { productApi } from '@/api/product';
import { useProduct } from '@/hooks/useProduct';

const { Option } = Select;

const sortOptions = [
  { value: 1, label: 'Highest Priority - Top of List' },
  { value: 2, label: 'High Priority' },
  { value: 3, label: 'Normal Priority' },
  { value: 4, label: 'Low Priority - Bottom of List' }
];

const CategoryForm = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<Category[]>([]);
  const { refreshCategories } = useProduct();

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const data = await productApi.listCategories();
        setCategories(data);
      } catch (error) {
        console.error('Failed to load categories:', error);
        message.error('Failed to load categories');
      }
    };
    fetchCategories();
  }, []);

  const onSubmit = async (values: CreateCategoryReq) => {
    setLoading(true);
    try {
      const submitData: CreateCategoryReq = {
        name: values.name,
        parentId: values.parentId || 0, // Set default to 0 for root category
        sort: values.sort || 1, // Default sort order
        icon: values.icon
      };

      await productApi.createCategory(submitData);
      await refreshCategories();
      message.success('Category created successfully');
      navigate('/');
    } catch (error) {
      console.error('Failed to create category:', error);
      message.error('Failed to create category');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form 
      form={form} 
      layout="vertical" 
      onFinish={onSubmit}
      initialValues={{
        sort: 1, // Default sort value
        parentId: 0 // Default parent ID
      }}
    >
      <Form.Item 
        name="name" 
        label="Category Name" 
        rules={[{ required: true, message: 'Please input category name' }]}
      >
        <Input placeholder="Enter category name" />
      </Form.Item>

      <Form.Item 
        name="parentId" 
        label="Parent Category (Optional)"
        help="Leave empty to create root category"
      >
        <Select 
          allowClear 
          placeholder="Select parent category (optional)"
        >
          {categories.map(category => (
            <Option key={category.id} value={category.id}>
              {category.name}
            </Option>
          ))}
        </Select>
      </Form.Item>

      <Form.Item 
        name="sort" 
        label="Display Priority"
        help="Controls where this category appears in lists"
      >
        <Select placeholder="Select display priority">
          {sortOptions.map(option => (
            <Option key={option.value} value={option.value}>
              {option.label}
            </Option>
          ))}
        </Select>
      </Form.Item>

      <Button type="primary" htmlType="submit" loading={loading}>
        Create Category
      </Button>
    </Form>
  );
};

export default CategoryForm;