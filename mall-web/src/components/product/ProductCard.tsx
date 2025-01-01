import { Card, Button, Badge, message } from 'antd'
import { Link } from 'react-router-dom'
import { ShoppingCartOutlined, HeartOutlined } from '@ant-design/icons'

interface ProductCardProps {
  id: string | number
  name: string
  price: number
  image: string
  description: string
  stock: number
  status: string
}

const ProductCard: React.FC<ProductCardProps> = ({ 
  id, 
  name, 
  price, 
  image, 
  description,
  stock,
  status 
}) => {
  const handleAddToCart = () => {
    message.success('Added to cart')
  }

  return (
    <Badge.Ribbon 
      text={status} 
      color={status === 'ON_SALE' ? 'green' : 'red'}
    >
      <Card
        hoverable
        cover={
          <Link to={`/products/${id}`}>
            <img 
              alt={name} 
              src={image} 
              className="h-48 w-full object-cover"
            />
          </Link>
        }
        actions={[
          <Button 
            key="addToCart"
            type="primary" 
            icon={<ShoppingCartOutlined />}
            onClick={handleAddToCart}
            disabled={stock === 0 || status !== 'ON_SALE'}
          >
            Add to Cart
          </Button>,
          <Button 
            key="favorite"
            icon={<HeartOutlined />} 
          />
        ]}
      >
        <Card.Meta
          title={
            <Link to={`/products/${id}`} className="text-lg font-medium">
              {name}
            </Link>
          }
          description={
            <div className="mt-2">
              <p className="text-gray-500 line-clamp-2">{description}</p>
              <div className="mt-4 flex justify-between items-center">
                <span className="text-lg font-bold text-red-500">
                  ${price.toFixed(2)}
                </span>
                <span className={`text-sm ${stock > 0 ? 'text-green-500' : 'text-red-500'}`}>
                  {stock > 0 ? `${stock} in stock` : 'Out of stock'}
                </span>
              </div>
            </div>
          }
        />
      </Card>
    </Badge.Ribbon>
  )
}

export default ProductCard