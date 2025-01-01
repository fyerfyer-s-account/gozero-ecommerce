import { createBrowserRouter } from 'react-router-dom'
import MainLayout from './components/layout/MainLayout'
import LoginPage from './pages/user/LoginPage'
import ProfilePage from './pages/user/ProfilePage'
import ProductListPage from './pages/product/ListPage'
import ProductDetailPage from './pages/product/DetailPage'
import SignUpPage from './pages/user/SignUpPage'

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />
  },
  {
    path: '/register', 
    element: <SignUpPage />
  },
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
        element: <ProfilePage />
      }
    ]
  }
])

export default router