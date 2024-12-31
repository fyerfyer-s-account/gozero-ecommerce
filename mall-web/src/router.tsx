Sure, here's the contents for the file `/mall-web/mall-web/src/router.tsx`:

import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import MainLayout from './components/layout/MainLayout';
import LoginPage from './pages/user/LoginPage';
import ProfilePage from './pages/user/ProfilePage';
import ListPage from './pages/product/ListPage';
import DetailPage from './pages/product/DetailPage';

const AppRouter = () => {
    return (
        <Router>
            <MainLayout>
                <Routes>
                    <Route path="/login" element={<LoginPage />} />
                    <Route path="/profile" element={<ProfilePage />} />
                    <Route path="/products" element={<ListPage />} />
                    <Route path="/products/:id" element={<DetailPage />} />
                </Routes>
            </MainLayout>
        </Router>
    );
};

export default AppRouter;