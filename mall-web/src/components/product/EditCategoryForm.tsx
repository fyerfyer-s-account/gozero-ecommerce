import React from 'react';
import { Form, Input, InputNumber, Button, message } from 'antd';
import { productApi } from '@/api/product';
import { Category } from '@/types/product';

interface Props {
  categoryId: number;
}

const EditCategoryForm: React.FC<Props> = ({ categoryId }) => {
  const [form] = Form.useForm();

  const onFinish = async (values: Partial<Category>) => {
    try {
      await productApi.updateCategory(categoryId, values);
      message.success('Category updated successfully');
    } catch (error) {
      message.error('Failed to update category');
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={onFinish}
    >
      <Form.Item name="name" label="Category Name">
        <Input />
      </Form.Item>
      <Form.Item name="level" label="Level">
        <InputNumber min={1} />
      </Form.Item>
      <Form.Item name="sort" label="Sort Order">
        <InputNumber min={0} />
      </Form.Item>
      <Form.Item name="icon" label="Icon URL">
        <Input />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Update Category
        </Button>
      </Form.Item>
    </Form>
  );
};

export default EditCategoryForm;