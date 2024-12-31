Sure, here's the contents for the file `/mall-web/mall-web/src/utils/api.ts`:

import { useEffect, useState } from 'react';

const API_BASE_URL = 'https://api.example.com'; // Replace with your actual API base URL

export const fetchData = async (endpoint: string, options?: RequestInit) => {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, options);
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
    return response.json();
};

export const useFetch = (endpoint: string, options?: RequestInit) => {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchDataAsync = async () => {
            try {
                const result = await fetchData(endpoint, options);
                setData(result);
            } catch (error) {
                setError(error);
            } finally {
                setLoading(false);
            }
        };

        fetchDataAsync();
    }, [endpoint, options]);

    return { data, loading, error };
};