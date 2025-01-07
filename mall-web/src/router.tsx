import { createBrowserRouter, Navigate, useLocation } from 'react-router-dom'
import MainLayout from './components/layout/MainLayout'
import LoginPage from './pages/user/LoginPage'
import ProfilePage from './pages/user/ProfilePage'
import ProductListPage from './pages/product/ListPage'
import ProductDetailPage from './pages/product/DetailPage'
import SignUpPage from './pages/user/SignUpPage'
import CreatePage from './pages/product/CreatePage'
import CategoryPage from './pages/product/CategoryPage'
import { useAuth } from './hooks/useAuth'
import EditPage from './pages/product/EditPage'

// Add auth guard component
const RequireAuth = ({ children }: { children: JSX.Element }) => {
  const location = useLocation();
  const { isAuthenticated } = useAuth(); // Use auth context instead of direct token check
  
  if (!isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  
  return children;
};

const PublicRoute = ({ children }: { children: JSX.Element }) => {
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  
  if (isAuthenticated) {
    return <Navigate to={location.state?.from?.pathname || '/'} replace />;
  }
  
  return children;
};

export const router = createBrowserRouter([
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        path: '/',
        element: <ProductListPage />
      },
      {
        path: '/products/:id',
        element: <ProductDetailPage />
      },
      {
        path: '/profile',
        element: <RequireAuth><ProfilePage /></RequireAuth>
      },
      {
        path: '/login',
        element: <PublicRoute><LoginPage /></PublicRoute>
      },
      {
        path: '/register',
        element: <PublicRoute><SignUpPage /></PublicRoute>
      },
      {
        path: '/products/create',
        element: <RequireAuth><CreatePage /></RequireAuth>
      },
      {
        path: '/categories/create',
        element: <RequireAuth><CategoryPage /></RequireAuth>
      },
      {
        path: '/products/:id/edit',
        element: <EditPage />
      }
    ]
  }
]);

export default router