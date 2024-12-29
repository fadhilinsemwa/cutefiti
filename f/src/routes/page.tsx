// src/routes/page.tsx
import Layout from '@/routes/layout';
import './index.css';
import { AuthProvider } from '../context/AuthContext';
import { overrideConsoleLog } from '../utils/log';

// Call the function to override console.log behavior
overrideConsoleLog();

const Index = () => {
  return (
    <div className="container-box">
      <AuthProvider>
        <Layout />
      </AuthProvider>
    </div>
  );
};

export default Index;
