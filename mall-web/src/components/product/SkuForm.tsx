import { Form, Input, Button, Space } from 'antd';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { SkuAttributeReq } from '@/types/product';

interface SkuFormProps {
  onChange: (attributes: SkuAttributeReq[]) => void;
}

const SkuForm: React.FC<SkuFormProps> = ({ onChange }) => {
  return (
    <Form.List name="skuAttributes">
      {(fields, { add, remove }) => (
        <>
          {fields.map(({ key, name, ...restField }) => (
            <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
              <Form.Item
                {...restField}
                name={[name, 'key']}
                rules={[{ required: true, message: 'Missing key' }]}
              >
                <Input placeholder="Attribute Key (e.g. Color)" />
              </Form.Item>
              <Form.Item
                {...restField}
                name={[name, 'value']}
                rules={[{ required: true, message: 'Missing value' }]}
              >
                <Input placeholder="Attribute Value (e.g. Red)" />
              </Form.Item>
              <MinusCircleOutlined onClick={() => remove(name)} />
            </Space>
          ))}
          <Form.Item>
            <Button 
              type="dashed" 
              onClick={() => add()} 
              block 
              icon={<PlusOutlined />}
            >
              Add SKU Attribute
            </Button>
          </Form.Item>
        </>
      )}
    </Form.List>
  );
};

export default SkuForm;