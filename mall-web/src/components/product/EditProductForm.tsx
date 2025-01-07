import React from 'react';
import { Form, Input, InputNumber, Upload, Button, message } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import { useProduct } from '@/hooks/useProduct';
import { UpdateProductReq } from '@/types/product';
import { productApi } from '@/api/product';

interface Props {
  productId: string;
}

const EditProductForm: React.FC<Props> = ({ productId }) => {
  const { currentProduct, fetchProductDetails } = useProduct();
  const [form] = Form.useForm();

  const onFinish = async (values: UpdateProductReq) => {
    try {
      await productApi.updateProduct(parseInt(productId), values);
      message.success('Product updated successfully');
      fetchProductDetails(productId);
    } catch (error) {
      message.error('Failed to update product');
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      initialValues={currentProduct}
      onFinish={onFinish}
    >
      <Form.Item name="name" label="Product Name">
        <Input />
      </Form.Item>
      <Form.Item name="description" label="Description">
        <Input.TextArea rows={4} />
      </Form.Item>
      <Form.Item name="price" label="Price">
        <InputNumber min={0} />
      </Form.Item>
      <Form.Item name="brand" label="Brand">
        <Input />
      </Form.Item>
      <Form.Item name="images" label="Images">
        <Upload>
          <Button icon={<UploadOutlined />}>Select Images</Button>
        </Upload>
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Update Product
        </Button>
      </Form.Item>
    </Form>
  );
};

export default EditProductForm;