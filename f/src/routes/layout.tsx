// src/routes/layout.tsx
import './index.css';
import { NavLink, useLocation, useNavigate } from '@modern-js/runtime/router';
import { useEffect, useState } from 'react';
import { FaHome, FaCog, FaTable, FaServer } from 'react-icons/fa'; // Import FaServer for ApiTracking
import Dashboard from './dashboard/Dashboard';
import Users from '../components/Users';
import Roles from '../components/Roles';
import Settings from '../pages/Settings';
import UserProfileDropdown from '../components/UserProfileDropdown';
import SidePanel from '../components/SidePanel';
import RolesComponents from '../components/RolesComponents';
import LoginForm from '../components/LoginForm';
import Modal from '../components/Modal';
import ProtectedComponent from '../components/ProtectedComponent';
import PaymentPostResponseTable from './mupayapi/components/PaymentPostResponseTable';
import ApiTracking from './mupayapi/pages/ApiTracking'; // Import ApiTracking
import CRDBPaymentTest from './mupayapi/components/pay1'; // Import CRDB Payment Test
import { AuthProvider } from '../context/AuthContext';

const Layout = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [initialLoad, setInitialLoad] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);

  useEffect(() => {
    if (initialLoad) {
      if (!location.hash) {
        navigate('#dashboard', { replace: true });
      }
      setInitialLoad(false);
    }

    // Open the modal when navigating to #login
    if (location.hash === '#login') {
      setIsModalOpen(true);
    } else {
      setIsModalOpen(false);
    }
  }, [initialLoad, location.hash, navigate]);

  const handleProtectedRouteAccess = (hash: string, event: React.MouseEvent<HTMLAnchorElement, MouseEvent>) => {
    event.preventDefault();
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('#login', { state: { from: hash } });
    } else {
      navigate(hash);
    }
  };

  const getActiveClass = (hash: string) =>
    location.hash === hash
      ? 'text-[#e68f24] border-[#911841]'
      : 'text-[#911841] border-transparent';

  return (
    <AuthProvider>
      <div className="relative min-h-screen overflow-hidden">
        {/* Navigation Bar */}
        <nav className="fixed top-0 left-0 right-0 bg-white shadow-md z-50">
          <div className="flex justify-between items-center py-1 px-3 relative">
            {/* SidePanel Toggle Icon */}
            <SidePanel />

            <div className="flex items-center space-x-3">
              {/* Dashboard and Settings Links */}
              <NavLink
                to="#dashboard"
                className={`flex flex-col items-center ${getActiveClass('#dashboard')}`}
              >
                <FaHome className="text-lg" />
                <span className="text-xs">Dashboard</span>
              </NavLink>
              <NavLink
                to="#settings"
                className={`flex flex-col items-center ${getActiveClass('#settings')}`}
                onClick={(event) => handleProtectedRouteAccess('#settings', event)}
              >
                <FaCog className="text-lg" />
                <span className="text-xs">Settings</span>
              </NavLink>
              <NavLink
                to="#api-tracking"
                className={`flex flex-col items-center ${getActiveClass('#api-tracking')}`}
                onClick={(event) => handleProtectedRouteAccess('#api-tracking', event)}
              >
                <FaServer className="text-lg" />
                <span className="text-xs">API Tracking</span>
              </NavLink>

            </div>

            {/* UserProfileDropdown */}
            <div className="relative">
              <UserProfileDropdown />
            </div>
          </div>
        </nav>

        {/* Page Content */}
        <main className="pt-16 h-full overflow-auto bg-gray-100">
          {location.hash === '#dashboard' && <Dashboard />}
          {location.hash === '#users' && (
            <ProtectedComponent component={Users} componentPath="src/components/Users.tsx" />
          )}
          {location.hash === '#roles' && (
            <ProtectedComponent component={Roles} componentPath="src/components/Roles.tsx" />
          )}
          {location.hash === '#roles-components' && (
            <ProtectedComponent component={RolesComponents} componentPath="src/components/RolesComponents.tsx" />
          )}
          {location.hash === '#payment-post-response-table' && (
            <ProtectedComponent component={PaymentPostResponseTable} componentPath="src/routes/mupayapi/components/PaymentPostResponseTable.tsx" />
          )}
          {location.hash === '#api-tracking' && (
            <ProtectedComponent component={ApiTracking} componentPath="src/routes/mupayapi/pages/ApiTracking.tsx" />
          )}
          {location.hash === '#crdb-payment-test' && ( // Render CRDB Payment Test
            <ProtectedComponent component={CRDBPaymentTest} componentPath="src/routes/mupayapi/components/pay1.tsx" />
          )}
          {location.hash === '#settings' && (
            <ProtectedComponent component={Settings} componentPath="src/pages/Settings.tsx" />
          )}
          {location.hash === '#unauthorized' && <div>Unauthorized Access</div>}
          {location.hash === '#error' && <div>An Error Occurred</div>}

          {/* Modal to display the LoginForm */}
          <Modal isOpen={isModalOpen} onClose={() => navigate('#dashboard')}>
            <LoginForm />
          </Modal>
        </main>
      </div>
    </AuthProvider>
  );
};

export default Layout;
