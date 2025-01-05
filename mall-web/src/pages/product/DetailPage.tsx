import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useProduct } from '../../hooks/useProduct';
import { Loading } from '../../components/common/Loading';

const DetailPage: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const { fetchProductDetails, currentProduct, loading, error } = useProduct();

    useEffect(() => {
        if (id) {
            fetchProductDetails(id);
        }
    }, [id, fetchProductDetails]);

    if (loading) {
        return <Loading />;
    }

    if (error) {
        return <div className="text-red-500">{error}</div>;
    }

    if (!currentProduct) {
        return <div>Product not found</div>;
    }

    return (
        <div className="container mx-auto px-4 py-8">
            <h1 className="text-2xl font-bold mb-4">{currentProduct.name}</h1>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div>
                    <img src={currentProduct.images[0]} alt={currentProduct.name} className="w-full" />
                </div>
                <div>
                    <p className="text-xl font-bold text-red-500">¥{currentProduct.price}</p>
                    <p className="mt-4">{currentProduct.description}</p>
                </div>
            </div>
        </div>
    );
};

export default DetailPage;