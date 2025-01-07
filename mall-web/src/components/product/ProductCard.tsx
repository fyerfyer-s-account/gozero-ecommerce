import { Card, Button, Badge, message } from 'antd'
import { Link, useNavigate } from 'react-router-dom'
import { ShoppingCartOutlined, HeartOutlined, EditOutlined } from '@ant-design/icons'
import { useAuth } from '@/context/AuthContext';
import goodsIcon from '@/icon/goods_icon.jpg';

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
  const { user } = useAuth();
  const navigate = useNavigate();

  const actions = [
    <Button
      key="addToCart"
      type="primary"
      icon={<ShoppingCartOutlined />}
      disabled={stock === 0 || status !== 'ON_SALE'}
    />,
    <Button 
      key="favorite"
      icon={<HeartOutlined />} 
    />
  ];

  // Add edit button for admin
  if (user?.role === 'admin') {
    actions.push(
      <Button
        key="edit"
        icon={<EditOutlined />}
        onClick={() => navigate(`/products/${id}/edit`)}
      />
    );
  }

  return (
    <Badge.Ribbon 
      text={status} 
      color={status === 'ON_SALE' ? 'green' : 'red'}
    >
      <Card
        hoverable
        cover={
          <img
            alt={name}
            src={image || goodsIcon}
            className="h-48 object-cover"
          />
        }
        actions={actions}
      >
        <Card.Meta
          title={<div className="text-lg font-bold">{name}</div>}
          description={
            <div>
              <div className="text-red-500 font-bold">¥{price}</div>
              <div className="text-gray-500 truncate">{description}</div>
            </div>
          }
        />
      </Card>
    </Badge.Ribbon>
  );
};

export default ProductCard