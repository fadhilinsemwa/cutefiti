// usersf/src/components/RolesComponents.tsx
import React, { useEffect, useState } from 'react';
import {
  createComponent,
  getAllComponents,
  updateComponent,
  deleteComponent,
  deleteAllComponents,
  truncateComponents,
  assignRoleToComponent,
  getAllRolesComponents,
  updateRoleComponent,
  deleteRoleComponent,
  deleteAllRoleComponents,
  truncateRolesComponents,
} from '../api/roles_components';
import 'tailwindcss/tailwind.css';

import { getAllRoles } from '../api/roles'; 

interface Role {
  role_id: number;
  role_name: string;
}

interface Component {
  id: number;
  component_name: string;
  component_path: string;
  created_at: string;
  updated_at: string;
}

interface RoleComponentMapping {
  id: number;
  role_id: number;
  component_id: number;
  component_path: string;
  created_at: string;
  updated_at: string;
}

const RolesComponents = () => {
  const [components, setComponents] = useState<Component[] | null>(null);
  const [rolesComponents, setRolesComponents] = useState<RoleComponentMapping[] | null>(null);
  const [roles, setRoles] = useState<Role[]>([]);
  const [filteredRoles, setFilteredRoles] = useState<Role[]>([]);
  const [newComponent, setNewComponent] = useState({ component_name: '', component_path: '' });
  const [selectedRole, setSelectedRole] = useState({ role_id: '', component_id: '' });
  const [activeTab, setActiveTab] = useState<'create' | 'list' | 'assign' | 'mappings'>('create');
  const [showDropdown, setShowDropdown] = useState(false);

  // Function to convert objects to arrays (if needed)
  const convertObjectToArray = (obj: Record<string, any>): Role[] => {
    return Object.keys(obj).map((key) => ({ id: key, ...obj[key] }));
  };

  useEffect(() => {
    fetchComponents();
    fetchRolesComponents();
    fetchRoles(); // Ensure roles are fetched
  }, []);

  // Fetch Roles logic remains unchanged
  const fetchRoles = async () => {
    try {
      let rolesData = await getAllRoles();
      console.log("Fetched Roles Data:", rolesData); // Log the fetched data
      rolesData = convertObjectToArray(rolesData);
      if (Array.isArray(rolesData)) {
        setRoles(rolesData);
        setFilteredRoles(rolesData); // Initialize the filtered list
      } else {
        console.error('Roles data is not an array', rolesData);
      }
    } catch (error) {
      console.error('Failed to fetch roles', error);
    }
  };

  const fetchComponents = async () => {
    try {
      const result = await getAllComponents();
      setComponents(result.data || []);
    } catch (error) {
      console.error('Failed to fetch components:', error);
      setComponents([]); // Fallback to empty array in case of error
    }
  };

  const fetchRolesComponents = async () => {
    try {
      const result = await getAllRolesComponents();
      setRolesComponents(result.data || []);
    } catch (error) {
      console.error('Failed to fetch roles-components mappings:', error);
      setRolesComponents([]); // Fallback to empty array in case of error
    }
  };

  // Map component ID's to names
  const getComponentNameById = (id: number): string => {
    if (!components) return 'Unknown Component';
    const component = components?.find(comp => comp.id === id);
    return component ? component.component_name : 'Unknown Component';
  };
  

  // Hardcoded component names and paths
  const hardcodedComponents = [
    { name: 'Dashboard', path: 'src/routes/dashboard/Dashboard.tsx' },
    { name: 'Profile', path: 'src/components/UserProfileDropdown.tsx' },
    { name: 'Settings', path: 'src/pages/Settings.tsx' },
    { name: 'Roles', path: 'src/components/Roles.tsx' },
    { name: 'Roles Components', path: 'src/components/RolesComponents.tsx' },
    { name: 'Users', path: 'src/components/Users.tsx' },
    { name: 'Payment Post Response', path: 'src/routes/mupayapi/components/PaymentPostResponseTable.tsx' },
    { name: 'CRDB Payment Test', path: 'src/routes/mupayapi/components/pay1.tsx' },
    { name: 'API Tracking', path: 'src/routes/mupayapi/pages/ApiTracking.tsx' },
    
  ];

  const handleComponentNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setNewComponent({ ...newComponent, component_name: value });
    setShowDropdown(value.length > 0); // Show dropdown if input is not empty
  };

  const handleComponentSelect = (name: string, path: string) => {
    setNewComponent({ component_name: name, component_path: path });
    setShowDropdown(false); // Hide dropdown after selection
  };

  const handleCreateComponent = async () => {
    await createComponent(newComponent);
    fetchComponents();
  };

  const handleUpdateComponent = async (id: number) => {
    await updateComponent(id, newComponent);
    fetchComponents();
  };

  const handleDeleteComponent = async (id: number) => {
    await deleteComponent(id);
    fetchComponents();
  };

  const handleDeleteAllComponents = async () => {
    await deleteAllComponents();
    fetchComponents();
  };

  const handleTruncateComponents = async () => {
    await truncateComponents();
    fetchComponents();
  };

  const handleAssignRoleToComponent = async () => {
    await assignRoleToComponent({
      role_id: parseInt(selectedRole.role_id),
      component_id: parseInt(selectedRole.component_id),
    });
    fetchRolesComponents();
  };

  const handleUpdateRoleComponent = async (id: number) => {
    await updateRoleComponent(id, {
      role_id: parseInt(selectedRole.role_id),
      component_id: parseInt(selectedRole.component_id),
    });
    fetchRolesComponents();
  };

  const handleDeleteRoleComponent = async (id: number) => {
    await deleteRoleComponent(id);
    fetchRolesComponents();
  };

  const handleDeleteAllRoleComponents = async () => {
    await deleteAllRoleComponents();
    fetchRolesComponents();
  };

  const handleTruncateRolesComponents = async () => {
    await truncateRolesComponents();
    fetchRolesComponents();
  };

  {console.log('Roles State:', roles)}
  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Roles & Components Management</h1>

      {/* Tab Navigation */}
      <div className="flex border-b mb-4">
        <button
          className={`px-4 py-2 focus:outline-none ${activeTab === 'create' ? 'border-b-2 border-blue-500' : ''}`}
          onClick={() => setActiveTab('create')}
        >
          Create a Component
        </button>
        <button
          className={`px-4 py-2 focus:outline-none ${activeTab === 'list' ? 'border-b-2 border-blue-500' : ''}`}
          onClick={() => setActiveTab('list')}
        >
          Components List
        </button>
        <button
          className={`px-4 py-2 focus:outline-none ${activeTab === 'assign' ? 'border-b-2 border-blue-500' : ''}`}
          onClick={() => setActiveTab('assign')}
        >
          Assign Role to Component
        </button>
        <button
          className={`px-4 py-2 focus:outline-none ${activeTab === 'mappings' ? 'border-b-2 border-blue-500' : ''}`}
          onClick={() => setActiveTab('mappings')}
        >
          Roles-Components Mappings
        </button>
      </div>

      {/* Content */}
      {activeTab === 'create' && (
        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-2">Create a Component</h2>
          <input
            className="border p-2 rounded mb-2 w-full"
            placeholder="Component Name"
            value={newComponent.component_name}
            onChange={handleComponentNameChange}
            onFocus={() => setShowDropdown(true)}
            onBlur={() => setShowDropdown(false)}
          />
          {showDropdown && (
            <div className="relative">
              <ul className="absolute z-10 bg-white text-primary border border-gray-300 rounded-lg shadow-md mt-1 w-full">
                {hardcodedComponents
                  .filter(({ name }) =>
                    name.toLowerCase().includes(newComponent.component_name.toLowerCase())
                  )
                  .map(({ name, path }) => (
                    <li
                      key={name}
                      className="p-2 hover:bg-blue-500 hover:text-white cursor-pointer"
                      onMouseDown={() => handleComponentSelect(name, path)} // Using onMouseDown to prevent blur
                    >
                      {name}
                    </li>
                  ))}
              </ul>
            </div>
          )}


          <input
            className="border p-2 rounded mb-2 w-full"
            placeholder="Component Path"
            value={newComponent.component_path}
            readOnly
          />
          <button
            className="bg-blue-500 text-white p-2 rounded"
            onClick={handleCreateComponent}
          >
            Create Component
          </button>
        </div>
      )}

      {activeTab === 'list' && (
        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-2">Components List</h2>
          <div className="space-y-4">
            {components && components.map((component) => (
              <div key={component.id} className="border p-4 rounded">
                <p>Name: {component.component_name}</p>
                <p>Path: {component.component_path}</p>
                <button
                  className="bg-yellow-500 text-white p-2 rounded mt-2"
                  onClick={() => handleUpdateComponent(component.id)}
                >
                  Update
                </button>
                <button
                  className="bg-red-500 text-white p-2 rounded mt-2"
                  onClick={() => handleDeleteComponent(component.id)}
                >
                  Delete
                </button>
              </div>
            ))}
          </div>
          <button
            className="bg-red-600 text-white p-2 rounded mt-4"
            onClick={handleDeleteAllComponents}
          >
            Delete All Components
          </button>
          <button
            className="bg-gray-500 text-white p-2 rounded mt-2"
            onClick={handleTruncateComponents}
          >
            Truncate Components
          </button>
        </div>
      )}

      {activeTab === 'assign' && (
        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-2">Assign Role to Component</h2>

          <select
            className="border p-2 rounded mb-2 w-full"
            value={selectedRole.role_id}
            onChange={(e) => setSelectedRole({ ...selectedRole, role_id: e.target.value })}
          >
            <option value="">Select Role</option>
            {roles.length > 0 ? roles.map((role) => (
              <option key={role.role_id} value={role.role_id}>
                {role.role_name}
              </option>
            )) : <option disabled>No roles available</option>}
          </select>
          <select
            className="border p-2 rounded mb-2 w-full"
            value={selectedRole.component_id}
            onChange={(e) => setSelectedRole({ ...selectedRole, component_id: e.target.value })}
          >
            <option value="">Select Component</option>
            {components && components.map((component) => (
              <option key={component.id} value={component.id}>
                {component.component_name}
              </option>
            ))}
          </select>
          <button
            className="bg-blue-500 text-white p-2 rounded"
            onClick={handleAssignRoleToComponent}
          >
            Assign Role to Component
          </button>
        </div>
      )}

      {activeTab === 'mappings' && (
        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-2">Roles-Components Mappings</h2>
          <div className="space-y-4">
          {rolesComponents && rolesComponents.map((mapping) => (
            <div key={mapping.id} className="border p-4 rounded">
              <p>Role ID: {mapping.role_id}</p> {/* role id */}
              <p>Role Name: {roles.find(role => role.role_id === mapping.role_id)?.role_name || 'Unknown Role'}</p> {/* role name*/}
              <p>Component Name: {getComponentNameById(mapping.component_id)}</p> {/* component name */}
              <p>Component Path: {components?.find(comp => comp.id === mapping.component_id)?.component_path || 'Unknown Path'}</p>
              <button
                className="bg-yellow-500 text-white p-2 rounded mt-2"
                onClick={() => handleUpdateRoleComponent(mapping.id)}
              >
                Update Mapping
              </button>
              <button
                className="bg-red-500 text-white p-2 rounded mt-2"
                onClick={() => handleDeleteRoleComponent(mapping.id)}
              >
                Delete Mapping
              </button>
            </div>
          ))}

          </div>
          <button
            className="bg-red-600 text-white p-2 rounded mt-4"
            onClick={handleDeleteAllRoleComponents}
          >
            Delete All Mappings
          </button>
          <button
            className="bg-gray-500 text-white p-2 rounded mt-2"
            onClick={handleTruncateRolesComponents}
          >
            Truncate Mappings
          </button>
        </div>
      )}
    </div>
  );
};

export default RolesComponents;
