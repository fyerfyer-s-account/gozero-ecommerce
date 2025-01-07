import React from 'react';
import { Form, InputNumber, Button, Space, message } from 'antd';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { Sku } from '@/types/product';
import { productApi } from '@/api/product';

interface Props {
  productId: string;
  skus: Sku[];
}

const EditSkuForm: React.FC<Props> = ({ productId, skus }) => {
  const [form] = Form.useForm();

  const onFinish = async (values: { skus: Sku[] }) => {
    try {
      const updatePromises = values.skus.map(sku => 
        productApi.updateSku(sku.id, {
          price: sku.price,
          stock: sku.stock
        })
      );
      
      await Promise.all(updatePromises);
      message.success('SKUs updated successfully');
    } catch (error) {
      message.error('Failed to update SKUs');
    }
  };

  return (
    <Form form={form} onFinish={onFinish} initialValues={{ skus }}>
      <Form.List name="skus">
        {(fields, { add, remove }) => (
          <>
            {fields.map(({ key, name, ...restField }) => (
              <Space key={key} align="baseline">
                <Form.Item
                  {...restField}
                  name={[name, 'price']}
                  label="Price"
                >
                  <InputNumber min={0} />
                </Form.Item>
                <Form.Item
                  {...restField}
                  name={[name, 'stock']}
                  label="Stock"
                >
                  <InputNumber min={0} />
                </Form.Item>
                <MinusCircleOutlined onClick={() => remove(name)} />
              </Space>
            ))}
            <Form.Item>
              <Button 
                type="dashed" 
                onClick={() => add()} 
                icon={<PlusOutlined />}
              >
                Add SKU
              </Button>
            </Form.Item>
          </>
        )}
      </Form.List>
      <Form.Item>
        <Button type="primary" htmlType="submit">
          Update SKUs
        </Button>
      </Form.Item>
    </Form>
  );
};

export default EditSkuForm;