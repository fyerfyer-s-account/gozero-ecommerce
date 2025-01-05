import ProductForm from '@/components/product/ProductForm';
import { Card } from 'antd';

const CreatePage = () => {
  return (
    <div className="container mx-auto px-4 py-8">
      <Card title="Create New Product">
        <ProductForm />
      </Card>
    </div>
  );
};

export default CreatePage;