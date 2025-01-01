import { Row, Col, Input, Select, Spin, Empty } from 'antd'
import { useProduct } from '../../hooks/useProduct'
import ProductCard from './ProductCard'
import { Product } from '../../types/product'

const { Search } = Input
const { Option } = Select

const ProductList = () => {
  const { products, loading, error, searchProducts } = useProduct()

  const handleSearch = (keyword: string) => {
    searchProducts({ keyword });
  }

  if (loading) return <Spin size="large" className="flex justify-center" />;
  if (error) return <div className="text-red-500 text-center">{error}</div>;
  if (!products || products.length === 0) return (
    <div className="flex justify-center items-center p-8">
      <Empty description="No products found" />
    </div>
  );

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <Search
          placeholder="Search products"
          onSearch={handleSearch}
          style={{ width: 300 }}
        />
      </div>
      <Row gutter={[16, 16]}>
        {products.map((product) => (
          <Col key={product.id} xs={24} sm={12} md={8} lg={6}>
            <ProductCard 
              id={product.id}
              name={product.name}
              price={product.price}
              image={product.images?.[0] || '/default-image.jpg'} 
              description={product.description}
              stock={product.stock}
              status={product.status === 1 ? 'ON_SALE' : 'OFF_SALE'} 
            />
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default ProductList;