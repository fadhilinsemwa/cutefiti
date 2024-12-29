// src/pages/ApiTracking.tsx
import { useState } from 'react';
import { FaFileInvoiceDollar } from 'react-icons/fa'; // Icon for PaymentPostResponseTable
import { FaMoneyCheckAlt } from 'react-icons/fa'; // Icon for Pay1
import { useNavigate } from '@modern-js/runtime/router';

const ApiTracking = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const navigate = useNavigate();

  const icons = [
    {
      name: 'Payment Post Response',
      icon: <FaFileInvoiceDollar className="text-xl" />, // Adjusted size to 13px
      action: () => navigate('#payment-post-response-table')
    },
    {
      name: 'CRDB Payment Test',
      icon: <FaMoneyCheckAlt className="text-xl" />, // Icon for Pay1
      action: () => navigate('#crdb-payment-test')
    }
  ];

  const filteredIcons = icons.filter((icon) =>
    icon.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="flex flex-col items-center h-full text-black p-4 overflow-y-auto">
      {/* Search Bar */}
      <div className="w-full max-w-md mb-8 mt-2">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Search API tracking..."
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
            <div className="flex flex-col items-center">
              <div className="p-4 bg-primary text-white rounded-full hover:bg-[#e68f24] transition-all duration-300">
                {icon.icon}
              </div>
              <span className="mt-2 text-xs text-primary text-center">{icon.name}</span> {/* Adjusted size to 13px */}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ApiTracking;
