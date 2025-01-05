import { Layout, Menu } from 'antd';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import React from 'react';

const Header: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();
  console.log('Current user in header:', user); // Debug log
  const isAdmin = user?.role === 'admin';
  console.log('Is admin?', isAdmin); // Debug log

  return (
    <Layout.Header className="bg-white shadow">
      <div className="container mx-auto flex justify-between items-center">
        <div className="flex items-center">
          <Link to="/" className="text-xl font-bold mr-4">Mall</Link>
          {isAdmin && (
            <span className="bg-blue-500 text-white px-2 py-1 rounded text-xs">
              Admin
            </span>
          )}
        </div>
        <Menu 
          mode="horizontal" 
          className="border-0"
          style={{ minWidth: '300px' }} // Add minimum width
        >
          <Menu.Item key="home">
            <Link to="/">Home</Link>
          </Menu.Item>
          {isAuthenticated ? (
            <>
              {isAdmin && (
                <>
                  <Menu.Item key="create-product">
                    <Link to="/products/create">Create Product</Link>
                  </Menu.Item>
                  <Menu.Item key="create-category">
                    <Link to="/categories/create">Create Category</Link>
                  </Menu.Item>
                </>
              )}
              <Menu.Item key="profile">
                <Link to="/profile">Profile</Link>
              </Menu.Item>
              <Menu.Item key="logout" onClick={() => {
                logout();
                navigate('/');
              }}>
                Logout
              </Menu.Item>
            </>
          ) : (
            <>
              <Menu.Item key="login">
                <Link to="/login">Login</Link>
              </Menu.Item>
              <Menu.Item key="signup">
                <Link to="/register">Sign Up</Link>
              </Menu.Item>
            </>
          )}
        </Menu>
      </div>
    </Layout.Header>
  );
};

export default Header;