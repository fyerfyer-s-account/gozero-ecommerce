import { AuthProvider } from './context/AuthContext';
import { ProductProvider } from './context/ProductContext';
import { RouterProvider } from 'react-router-dom';
import router from './router';

const App = () => {
  return (
    <AuthProvider>
      <ProductProvider>
        <RouterProvider router={router} />
      </ProductProvider>
    </AuthProvider>
  );
};

export default App;