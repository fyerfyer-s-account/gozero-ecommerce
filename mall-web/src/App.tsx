Sure, here's the contents for the file `/mall-web/mall-web/src/App.tsx`:

import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { ProductProvider } from './context/ProductContext';
import RouterConfig from './router';

const App = () => {
  return (
    <AuthProvider>
      <ProductProvider>
        <Router>
          <RouterConfig />
        </Router>
      </ProductProvider>
    </AuthProvider>
  );
};

export default App;