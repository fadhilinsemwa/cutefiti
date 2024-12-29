// src/pages/Settings.tsx
import { FaUsers, FaUserShield, FaCogs } from 'react-icons/fa'; // FaCogs added for RolesComponents
import { useState } from 'react';
import { useNavigate } from '@modern-js/runtime/router';

const Settings = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();

  const icons = [
    { name: 'User Manager', icon: <FaUsers className="text-5xl" />, action: () => navigate('#users') },
    { name: 'User Role Manager', icon: <FaUserShield className="text-5xl" />, action: () => navigate('#roles') },
    { name: 'Roles-Components Manager', icon: <FaCogs className="text-5xl" />, action: () => navigate('#roles-components') }, // Added RolesComponents Manager
  ];

  const filteredIcons = icons.filter((icon) =>
    icon.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="flex flex-col items-center h-full text-white p-4 overflow-y-auto">
      {/* Search Bar */}
      <div className="w-full max-w-md mb-8 mt-2">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search settings..."
          className="w-full py-2 px-4 rounded-full text-black bg-gray-100 outline-none shadow-sm border border-secondary border-opacity-50 placeholder-gray-500 focus:placeholder-transparent transition-all duration-300"
          style={{ marginTop: '10px' }}
        />
      </div>

      {/* Icon Grid */}
      <div className="grid grid-cols-2 gap-10">
        {filteredIcons.map((icon, index) => (
          <div
            key={index}
            className="flex flex-col items-center cursor-pointer"
            onClick={icon.action}
          >
            <div className="p-6 bg-primary text-white rounded-full hover:bg-[#e68f24] transition-all duration-300">
              {icon.icon}
            </div>
            <span className="mt-4 text-lg text-primary">{icon.name}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Settings;
