import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useProduct } from '../../hooks/useProduct';
import { useAuth } from '../../hooks/useAuth';
import { Loading } from '../../components/common/Loading';
import { Card, Button, Tag, Divider, InputNumber, Radio, Tabs, message, Descriptions } from 'antd';
import { ShoppingCartOutlined, HeartOutlined, EditOutlined } from '@ant-design/icons';
import { Sku } from '@/types/product';
import goodsIcon from '@/icon/goods_icon.jpg';

const { TabPane } = Tabs;

const DetailPage: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const { fetchProductDetails, currentProduct, loading, error } = useProduct();
    const { user } = useAuth();
    const navigate = useNavigate();
    const [selectedSku, setSelectedSku] = useState<Sku | null>(null);
    const [quantity, setQuantity] = useState(1);
    const [mainImage, setMainImage] = useState('');

    useEffect(() => {
        if (id) {
            fetchProductDetails(id);
        }
    }, [id, fetchProductDetails]);

    useEffect(() => {
        if (currentProduct?.images && currentProduct.images.length > 0) {
            setMainImage(currentProduct.images[0]);
        }
    }, [currentProduct]);

    if (loading) return <Loading />;
    if (error) return <div className="text-red-500">{error}</div>;
    if (!currentProduct) return <div>Product not found</div>;

    const handleAddToCart = () => {
        if (!selectedSku) {
            message.warning('Please select product specifications');
            return;
        }
        message.success('Successfully added to cart');
    };

    const renderSkuSelector = () => {
        const skuGroups: { [key: string]: string[] } = {};
        currentProduct.skus?.forEach(sku => {
            Object.entries(sku.attributes).forEach(([key, value]) => {
                if (!skuGroups[key]) {
                    skuGroups[key] = [];
                }
                if (!skuGroups[key].includes(value)) {
                    skuGroups[key].push(value);
                }
            });
        });

        return Object.entries(skuGroups).map(([attrName, values]) => (
            <div key={attrName} className="mb-4">
                <div className="text-gray-600 mb-2">{attrName}</div>
                <Radio.Group 
                    onChange={(e) => {
                        const newSku = currentProduct.skus?.find(
                            sku => sku.attributes[attrName] === e.target.value
                        );
                        setSelectedSku(newSku || null);
                    }}
                >
                    {values.map(value => (
                        <Radio.Button key={value} value={value}>
                            {value}
                        </Radio.Button>
                    ))}
                </Radio.Group>
            </div>
        ));
    };

    return (
        <div className="container mx-auto px-4 py-8">
            <Card 
              extra={
                user?.role === 'admin' && (
                  <Button
                    icon={<EditOutlined />}
                    onClick={() => navigate(`/products/${id}/edit`)}
                  >
                    Edit Product
                  </Button>
                )
              }
            >
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    {/* Image Gallery */}
                    <div className="space-y-4">
                        <div className="border rounded-lg overflow-hidden">
                            <img 
                                src={mainImage || currentProduct.images?.[0] || goodsIcon} 
                                alt={currentProduct.name} 
                                className="w-full h-96 object-contain"
                            />
                        </div>
                        {currentProduct.images && currentProduct.images.length > 0 ? (
                            <div className="grid grid-cols-5 gap-2">
                                {currentProduct.images.map((img, index) => (
                                    <img 
                                        key={index}
                                        src={img || goodsIcon}
                                        alt={`${currentProduct.name}-${index}`}
                                        className={`w-full h-20 object-cover rounded cursor-pointer border 
                                            ${mainImage === img ? 'border-blue-500' : 'border-gray-200'}`}
                                        onClick={() => setMainImage(img)}
                                    />
                                ))}
                            </div>
                        ) : (
                            <div className="grid grid-cols-5 gap-2">
                                <img 
                                    src={goodsIcon}
                                    alt="default"
                                    className="w-full h-20 object-cover rounded cursor-pointer border border-gray-200"
                                />
                            </div>
                        )}
                    </div>

                    {/* Product Info */}
                    <div className="space-y-6">
                        <div>
                            <Tag color="blue">{currentProduct.categoryName}</Tag>
                            <h1 className="text-2xl font-bold mt-2">{currentProduct.name}</h1>
                            <p className="text-gray-500 mt-2">{currentProduct.brief}</p>
                        </div>

                        <Divider />

                        <div>
                            <div className="text-3xl text-red-500 font-bold">
                                ¥{selectedSku ? selectedSku.price : currentProduct.price}
                            </div>
                            <div className="mt-2 text-gray-500">
                                Stock: {selectedSku ? selectedSku.stock : currentProduct.stock}
                            </div>
                        </div>

                        <Divider />

                        {/* SKU Selection */}
                        {currentProduct.skus && currentProduct.skus.length > 0 && (
                            <div className="space-y-4">
                                {renderSkuSelector()}
                            </div>
                        )}

                        <div className="flex items-center space-x-4">
                            <span>Quantity:</span>
                            <InputNumber 
                                min={1} 
                                max={selectedSku?.stock || currentProduct.stock}
                                value={quantity} 
                                onChange={val => setQuantity(val || 1)}
                            />
                        </div>

                        <div className="flex space-x-4">
                            <Button 
                                type="primary" 
                                size="large"
                                icon={<ShoppingCartOutlined />} 
                                onClick={handleAddToCart}
                                disabled={!selectedSku || selectedSku.stock === 0}
                            >
                                Add to Cart
                            </Button>
                            <Button 
                                size="large"
                                icon={<HeartOutlined />}
                            >
                                Add to Wishlist
                            </Button>
                        </div>
                    </div>
                </div>

                <Tabs defaultActiveKey="1" className="mt-8">
                    <TabPane tab="Description" key="1">
                        <div className="prose max-w-none"
                            dangerouslySetInnerHTML={{ __html: currentProduct.description }}
                        />
                    </TabPane>
                    <TabPane tab="Specifications" key="2">
                        <Descriptions bordered>
                            <Descriptions.Item label="Brand">{currentProduct.brand}</Descriptions.Item>
                            <Descriptions.Item label="Category">{currentProduct.categoryName}</Descriptions.Item>
                            <Descriptions.Item label="Status">
                                {currentProduct.status === 1 ? 'In Stock' : 'Out of Stock'}
                            </Descriptions.Item>
                        </Descriptions>
                    </TabPane>
                </Tabs>
            </Card>
        </div>
    );
};

export default DetailPage;