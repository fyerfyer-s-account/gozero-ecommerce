import { Link } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import { useAuth } from '@/hooks/useAuth';

const Header = () => {
  const { user, logout } = useAuth();
  
  return (
    <Layout.Header className="bg-white shadow">
      <div className="container mx-auto flex justify-between items-center">
        <Link to="/" className="text-xl font-bold">Mall</Link>
        <Menu mode="horizontal" className="border-0">
          <Menu.Item key="home">
            <Link to="/">Home</Link>
          </Menu.Item>
          {user ? (
            <>
              <Menu.Item key="profile">
                <Link to="/profile">Profile</Link>
              </Menu.Item>
              <Menu.Item key="logout" onClick={logout}>
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

export default Header