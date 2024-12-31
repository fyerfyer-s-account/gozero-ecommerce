import React, { useEffect, useContext } from 'react';
import { useParams } from 'react-router-dom';
import { ProductContext } from '../../context/ProductContext';
import { Loading } from '../../components/common/Loading';
import { ProductCard } from '../../components/product/ProductCard';

const DetailPage: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const { fetchProductDetails, product, loading } = useContext(ProductContext);

    useEffect(() => {
        if (id) {
            fetchProductDetails(id);
        }
    }, [id, fetchProductDetails]);

    if (loading) {
        return <Loading />;
    }

    return (
        <div className="p-4">
            {product ? <ProductCard product={product} /> : <p>Product not found.</p>}
        </div>
    );
};

export default DetailPage;