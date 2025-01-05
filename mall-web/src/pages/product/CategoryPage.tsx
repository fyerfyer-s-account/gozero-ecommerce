import { Card } from 'antd';
import CategoryForm from '@/components/product/CategoryForm';

const CategoryPage = () => {
  return (
    <div className="container mx-auto px-4 py-8">
      <Card title="Create Category">
        <CategoryForm />
      </Card>
    </div>
  );
};

export default CategoryPage;