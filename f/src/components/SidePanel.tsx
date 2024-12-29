// src/components/SidePanel.tsx
import { useState } from 'react';
import { FaBars, FaTimes } from 'react-icons/fa';

const SidePanel = () => {
  const [isOpen, setIsOpen] = useState(false);

  const togglePanel = () => {
    setIsOpen(!isOpen);
  };

  return (
    <>
      {/* SidePanel */}
      <div
        className={`fixed top-0 left-0 h-screen bg-secondary text-white transition-transform duration-300 transform ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        } w-64 md:w-1/4 lg:w-1/5 z-40`}
        style={{
          willChange: 'transform',
        }}
      >
        <div className="p-4 pt-12"> {/* Added top padding */}
          <h2 className="text-2xl mb-6">Navigation</h2>
          <ul>
            <li className="mb-3">
              <a href="#dashboard" className="text-lg">
                Dashboard
              </a>
            </li>
            <li className="mb-3">
              <a href="#settings" className="text-lg">
                Settings
              </a>
            </li>
          </ul>
        </div>
      </div>

      {/* Toggler Icon */}
      <div
        onClick={togglePanel}
        className="p-2 bg-primary text-white rounded-md focus:outline-none z-50 cursor-pointer mt-2 ml-2"
      >
        {isOpen ? <FaTimes className="text-lg" /> : <FaBars className="text-lg" />}
      </div>

      {/* Overlay for mobile to close panel */}
      {isOpen && (
        <div
          onClick={togglePanel}
          className="fixed top-0 left-0 w-full h-full bg-black opacity-50 z-30 md:hidden"
        ></div>
      )}
    </>
  );
};

export default SidePanel;
