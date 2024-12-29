// src/components/UserProfileDropdown.tsx
import React, { useState, useEffect } from 'react';
import { FaUser } from 'react-icons/fa';
import { logoutUser } from '../api/authApi';
import { useNavigate, useLocation } from '@modern-js/runtime/router';
import Modal from './Modal';
import LoginForm from './LoginForm';
import Signup from './Signup';
import { useAuth } from '../context/AuthContext';

const UserProfileDropdown = () => {
  const { isAuthenticated, login, logout } = useAuth(); // Added login for future use if needed
  const [isOpen, setIsOpen] = useState(false);
  const [showLoginForm, setShowLoginForm] = useState(false);
  const [showSignupForm, setShowSignupForm] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  const toggleDropdown = () => {
    setIsOpen((prev) => !prev);
  };

  const handleLoginClick = () => {
    navigate('#login', { state: { from: location.hash } });
    setIsOpen(false);
  };

  const handleLogout = async () => {
    try {
      await logoutUser();
      localStorage.removeItem('token');
      localStorage.removeItem('role_id');
      logout(); // Update context state
      setIsOpen(false);
      navigate('#dashboard', { replace: true });
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const handleSignup = () => {
    setShowSignupForm(true);
    setIsOpen(false);
  };

  // Keep auth state updated based on local storage
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      login(); // Update auth context to logged in if token exists
    }
  }, [login]);

  return (
    <div className="relative">
      <button
        onClick={toggleDropdown}
        className="flex items-center justify-center w-8 h-8 bg-[#e68f24] text-white rounded-full shadow hover:bg-[#cc6e00] focus:outline-none"
      >
        <FaUser className="text-lg text-white" />
      </button>
      {isOpen && (
        <div
          className="absolute right-0 mt-2 w-48 bg-white border border-gray-200 rounded shadow-lg"
          style={{ zIndex: 9999 }}
        >
          <ul className="py-2">
            {isAuthenticated ? (
              <>
                <li className="px-4 py-2 hover:bg-gray-100 cursor-pointer">Profile</li>
                <li className="px-4 py-2 hover:bg-gray-100 cursor-pointer" onClick={handleLogout}>
                  Logout
                </li>
              </>
            ) : (
              <>
                <li className="px-4 py-2 hover:bg-gray-100 cursor-pointer" onClick={handleLoginClick}>
                  Login
                </li>
                <li className="px-4 py-2 hover:bg-gray-100 cursor-pointer" onClick={handleSignup}>
                  Sign Up
                </li>
              </>
            )}
          </ul>
        </div>
      )}

      {/* Modal for login */}
      <Modal isOpen={showLoginForm} onClose={() => setShowLoginForm(false)}>
        <LoginForm />
      </Modal>

      {/* Modal for signup */}
      <Modal isOpen={showSignupForm} onClose={() => setShowSignupForm(false)}>
        <Signup />
      </Modal>
    </div>
  );
};

export default UserProfileDropdown;
