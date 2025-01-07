import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button, Card, Tabs, message } from 'antd';
import { useAuth } from '@/hooks/useAuth';
import EditProductForm from '@/components/product/EditProductForm';
import EditCategoryForm from '@/components/product/EditCategoryForm';
import EditSkuForm from '@/components/product/EditSkuForm';
import { useProduct } from '@/hooks/useProduct';

const { TabPane } = Tabs;

const EditPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { currentProduct, loading } = useProduct();

  if (!user?.role || user.role !== 'admin') {
    message.error('Admin access required');
    navigate('/');
    return null;
  }

  if (loading) return <div>Loading...</div>;
  if (!currentProduct) return <div>Product not found</div>;

  return (
    <div className="container mx-auto px-4 py-8">
      <Card 
        title="Edit Product"
        extra={
          <Button onClick={() => navigate(`/products/${id}`)}>
            Back to Product
          </Button>
        }
      >
        <Tabs defaultActiveKey="product">
          <TabPane tab="Basic Info" key="product">
            <EditProductForm productId={id!} />
          </TabPane>
          <TabPane tab="Category" key="category">
            <EditCategoryForm categoryId={currentProduct.categoryId} />
          </TabPane>
          <TabPane tab="SKUs" key="skus">
            <EditSkuForm productId={id!} skus={currentProduct.skus || []} />
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default EditPage;