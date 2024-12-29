// usersf/src/components/Roles.tsx
import React, { useEffect, useState } from 'react';
import { createRole, getAllRoles, assignUserRole, getUserRoles } from '../api/roles';
import { getAllUsers } from '../api/users'; // Assuming you have this function to fetch users

interface Role {
  role_id: string;
  role_name: string;
}

interface UserRole {
  email: string;
  role_id: string;
  role_name: string;
}

interface User {
  username: string; // Adjust based on your user data structure
  email: string; // Adjust based on your user data structure
}

const Roles: React.FC = () => {
  const [roles, setRoles] = useState<Role[]>([]);
  const [filteredRoles, setFilteredRoles] = useState<Role[]>([]);
  const [userRoles, setUserRoles] = useState<UserRole[]>([]);
  const [users, setUsers] = useState<User[]>([]); // Added state for users
  const [newRole, setNewRole] = useState({ role_id: '', role_name: '' });
  const [assignRole, setAssignRole] = useState({ username: '', role_id: '' });
  const [activeTab, setActiveTab] = useState<'create' | 'available' | 'assign' | 'userRoles'>('create');
  const [isRoleDropdownOpen, setIsRoleDropdownOpen] = useState(false);
  const [isUserDropdownOpen, setIsUserDropdownOpen] = useState(false); // State for user dropdown
  const [filteredUsers, setFilteredUsers] = useState<User[]>([]); // State for filtered users

  useEffect(() => {
    fetchRoles();
    fetchUserRoles();
    fetchUsers(); // Fetch users on component mount
  }, []);

  const convertObjectToArray = (data: any): any[] => {
    if (typeof data === 'object' && !Array.isArray(data)) {
      return Object.values(data);
    }
    return data;
  };

  const fetchRoles = async () => {
    try {
      let rolesData = await getAllRoles();
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

  const fetchUsers = async () => { // Fetch users function
    try {
      let usersData = await getAllUsers(); // Adjust based on your API
      usersData = convertObjectToArray(usersData);
      if (Array.isArray(usersData)) {
        setUsers(usersData);
        setFilteredUsers(usersData); // Initialize the filtered users list
      } else {
        console.error('Users data is not an array', usersData);
      }
    } catch (error) {
      console.error('Failed to fetch users', error);
    }
  };

  const fetchUserRoles = async () => {
    try {
      let userRolesData = await getUserRoles();
      userRolesData = convertObjectToArray(userRolesData);
      if (Array.isArray(userRolesData)) {
        setUserRoles(userRolesData);
      } else {
        console.error('User roles data is not an array', userRolesData);
      }
    } catch (error) {
      console.error('Failed to fetch user roles', error);
    }
  };

  const handleCreateRole = async () => {
    try {
      const response = await createRole(newRole);
      alert(response.message);
      fetchRoles(); // Refresh roles list
    } catch (error) {
      console.error('Error creating role', error);
    }
  };

  const handleAssignUserRole = async () => {
    try {
      const response = await assignUserRole(assignRole);
      alert(response.message);
      fetchUserRoles(); // Refresh user roles list
    } catch (error) {
      console.error('Error assigning role to user', error);
    }
  };

  const handleRoleSearch = (value: string) => {
    setAssignRole({ ...assignRole, role_id: value });
    setFilteredRoles(
      roles.filter(role =>
        role.role_id.toLowerCase().includes(value.toLowerCase()) || role.role_name.toLowerCase().includes(value.toLowerCase())
      )
    );
  };

  const handleUserSearch = (value: string) => { // Handle user search
    setAssignRole({ ...assignRole, username: value });
    setFilteredUsers(
      users.filter(user =>
        user.username.toLowerCase().includes(value.toLowerCase()) // Adjust to match user data structure
      )
    );
  };

  const selectRole = (role_id: string) => {
    setAssignRole({ ...assignRole, role_id });
    setIsRoleDropdownOpen(false);
  };

  const selectUser = (username: string) => { // Select user from dropdown
    setAssignRole({ ...assignRole, username });
    setIsUserDropdownOpen(false);
  };

  return (
    <div className="roles-management container mx-auto px-4 py-8 text-white">
      <h2 className="text-3xl font-bold text-center text-primary mb-8">Roles Management</h2>

      {/* Horizontal Navigation */}
      <div className="flex justify-center mb-6">
        <button
          onClick={() => setActiveTab('create')}
          className={`px-4 py-2 mx-2 rounded-lg ${activeTab === 'create' ? 'bg-primary text-white' : 'bg-secondary text-gray-300'} hover:bg-primary-dark transition-all`}
        >
          Create New Role
        </button>
        <button
          onClick={() => setActiveTab('available')}
          className={`px-4 py-2 mx-2 rounded-lg ${activeTab === 'available' ? 'bg-primary text-white' : 'bg-secondary text-gray-300'} hover:bg-primary-dark transition-all`}
        >
          Available Roles
        </button>
        <button
          onClick={() => setActiveTab('assign')}
          className={`px-4 py-2 mx-2 rounded-lg ${activeTab === 'assign' ? 'bg-primary text-white' : 'bg-secondary text-gray-300'} hover:bg-primary-dark transition-all`}
        >
          Assign Role to Users
        </button>
        <button
          onClick={() => setActiveTab('userRoles')}
          className={`px-4 py-2 mx-2 rounded-lg ${activeTab === 'userRoles' ? 'bg-primary text-white' : 'bg-secondary text-gray-300'} hover:bg-primary-dark transition-all`}
        >
          User Roles
        </button>
      </div>

      {/* Tab Content */}
      {activeTab === 'create' && (
        <div className="bg-secondary p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-4">Create New Role</h3>
          <div className="space-y-4">
            <input
              type="text"
              placeholder="Role ID"
              value={newRole.role_id}
              onChange={(e) => setNewRole({ ...newRole, role_id: e.target.value })}
              className="w-full py-2 px-4 rounded-lg bg-gray-100 text-black focus:outline-none shadow-sm"
            />
            <input
              type="text"
              placeholder="Role Name"
              value={newRole.role_name}
              onChange={(e) => setNewRole({ ...newRole, role_name: e.target.value })}
              className="w-full py-2 px-4 rounded-lg bg-gray-100 text-black focus:outline-none shadow-sm"
            />
            <button
              onClick={handleCreateRole}
              className="w-full py-2 bg-primary text-white rounded-lg shadow-md hover:bg-primary-dark transition-all"
            >
              Create Role
            </button>
          </div>
        </div>
      )}

      {activeTab === 'available' && (
        <div className="bg-secondary p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-4">Available Roles</h3>
          {roles.length > 0 ? (
            <ul className="space-y-2">
              {roles.map((role) => (
                <li key={role.role_id} className="flex justify-between items-center p-3 bg-gray-800 rounded-lg shadow-sm">
                  <span>Role ID: {role.role_id}</span>
                  <span>Role Name: {role.role_name}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>No roles available.</p>
          )}
        </div>
      )}

      {activeTab === 'assign' && (
        <div className="bg-secondary p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-4">Assign Role to User</h3>
          <div className="space-y-4">
            <div className="relative">
              <input
                type="text"
                placeholder="Username"
                value={assignRole.username}
                onChange={(e) => handleUserSearch(e.target.value)} // Update to filter users
                className="w-full py-2 px-4 rounded-lg bg-gray-100 text-black focus:outline-none shadow-sm"
                onFocus={() => setIsUserDropdownOpen(true)} // Open dropdown on focus
                onBlur={() => setTimeout(() => setIsUserDropdownOpen(false), 200)} // Delay closing dropdown
              />
              {isUserDropdownOpen && (
                <ul className="absolute z-10 bg-white text-primary border border-gray-300 rounded-lg shadow-md max-h-60 overflow-y-auto w-full">
                  {filteredUsers.map((user) => (
                    <li
                      key={user.username}
                      onClick={() => selectUser(user.username)} // Set selected user
                      className="cursor-pointer hover:bg-gray-100 px-4 py-2"
                    >
                      {user.username}
                    </li>
                  ))}
                </ul>
              )}
            </div>
            <div className="relative">
              <input
                type="text"
                placeholder="Role ID"
                value={assignRole.role_id}
                onChange={(e) => handleRoleSearch(e.target.value)} // Update to filter roles
                className="w-full py-2 px-4 rounded-lg bg-gray-100 text-black focus:outline-none shadow-sm"
                onFocus={() => setIsRoleDropdownOpen(true)} // Open dropdown on focus
                onBlur={() => setTimeout(() => setIsRoleDropdownOpen(false), 200)} // Delay closing dropdown
              />
              {isRoleDropdownOpen && (
                <ul className="absolute z-10 bg-white text-primary border border-gray-300 rounded-lg shadow-md max-h-60 overflow-y-auto w-full">
                  {filteredRoles.map((role) => (
                    <li
                      key={role.role_id}
                      onClick={() => selectRole(role.role_id)} // Set selected role
                      className="cursor-pointer hover:bg-gray-100 px-4 py-2"
                    >
                      {role.role_id} - {role.role_name}
                    </li>
                  ))}
                </ul>
              )}
            </div>
            <button
              onClick={handleAssignUserRole}
              className="w-full py-2 bg-primary text-white rounded-lg shadow-md hover:bg-primary-dark transition-all"
            >
              Assign Role
            </button>
          </div>
        </div>
      )}

      {activeTab === 'userRoles' && (
        <div className="bg-secondary p-6 rounded-lg shadow-md">
          <h3 className="text-xl font-semibold mb-4">User Roles</h3>
          {userRoles.length > 0 ? (
            <ul className="space-y-2">
              {userRoles.map((userRole) => (
                <li key={userRole.email} className="flex justify-between items-center p-3 bg-gray-800 rounded-lg shadow-sm">
                  <span>{userRole.email}</span>
                  <span>{userRole.role_name}</span>
                </li>
              ))}
            </ul>
          ) : (
            <p>No user roles assigned.</p>
          )}
        </div>
      )}
    </div>
  );
};

export default Roles;
