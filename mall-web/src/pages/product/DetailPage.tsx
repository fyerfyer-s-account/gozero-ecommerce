import React, { useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useProduct } from '../../hooks/useProduct';
import { Loading } from '../../components/common/Loading';

const DetailPage: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const { fetchProductDetails, loading } = useProduct();

    useEffect(() => {
        if (id) {
            fetchProductDetails(id);
        }
    }, [id, fetchProductDetails]);

    if (loading) {
        return <Loading />;
    }

    return (
        <div className="container mx-auto px-4 py-8">
            {/* Product detail content */}
        </div>
    );
};

export default DetailPage;