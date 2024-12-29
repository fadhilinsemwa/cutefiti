// src/components/Users.tsx

import { useState, useEffect } from 'react';
import {
  createUser,
  updateUser,
  deleteUser,
  deleteAllUsers,
  truncateUsersTable,
  getAllUsers,
  User,
} from '../api/users';

const Users = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [editingUserId, setEditingUserId] = useState<number | null>(null);
  const [newUser, setNewUser] = useState<User | null>(null);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  useEffect(() => {
    loadUsers(); // Load all users on component mount
  }, []);

  const loadUsers = async () => {
    const allUsers = await getAllUsers();
    setUsers(allUsers);
  };

  const handleEditUser = (user: User) => {
    setEditingUserId(user.user_id ?? null); // Use user_id or null
    setEditingUser({ ...user });
  };

  const handleSaveUser = async () => {
    if (editingUser && editingUserId !== null) {
      await updateUser(editingUserId, editingUser);
      alert('User updated successfully!');
      setEditingUserId(null);
      setEditingUser(null);
      loadUsers();
    }
  };

  const handleDeleteUser = async (id: number) => {
    await deleteUser(id);
    alert('User deleted successfully!');
    loadUsers();
  };

  const handleAddUserRow = () => {
    setNewUser({
      user_id: undefined, // Set to undefined; modify as necessary
      username: '',
      email: '',
      firstname: '',
      lastname: '',
      city: '',
      password: '',
    });
  };

  const handleSaveNewUser = async () => {
    if (newUser) {
      await createUser(newUser);
      alert('New user added successfully!');
      setNewUser(null);
      loadUsers();
    }
  };

  const handleDeleteAllUsers = async () => {
    await deleteAllUsers();
    alert('All users deleted successfully!');
    loadUsers();
  };

  const handleTruncateUsersTable = async () => {
    await truncateUsersTable();
    alert('Users table truncated successfully!');
    loadUsers();
  };

  return (
    <div className="p-6 bg-gray-100 rounded-lg shadow-md">
      <h1 className="text-3xl font-bold mb-4">Users Management</h1>

      {/* Add User Button at the Top */}
      <div className="flex space-x-4 mb-4">
        <button onClick={handleAddUserRow} className="bg-primary text-white py-2 px-4 rounded">
          Add User
        </button>
      </div>

      {/* Users List */}
      <ul className="mb-4">
        {users.map((user) => (
          <li
            key={user.user_id ?? -1} // Use user_id as the key, fallback to -1 if undefined
            className="flex flex-col sm:grid sm:grid-cols-5 gap-2 justify-between items-center p-2 border-b border-gray-300"
          >
            {editingUserId === user.user_id ? (
              <>
                <div className="flex-1 mb-2 sm:mb-0 sm:col-span-4">
                  <input
                    type="text"
                    value={editingUser?.username || ''}
                    onChange={(e) =>
                      setEditingUser({ ...editingUser!, username: e.target.value })
                    }
                    className="input p-2 border border-gray-300 rounded mb-2 w-full"
                    placeholder="Username"
                  />
                  <input
                    type="email"
                    value={editingUser?.email || ''}
                    onChange={(e) =>
                      setEditingUser({ ...editingUser!, email: e.target.value })
                    }
                    className="input p-2 border border-gray-300 rounded mb-2 w-full"
                    placeholder="Email"
                  />
                  <input
                    type="text"
                    value={editingUser?.firstname || ''}
                    onChange={(e) =>
                      setEditingUser({ ...editingUser!, firstname: e.target.value })
                    }
                    className="input p-2 border border-gray-300 rounded mb-2 w-full"
                    placeholder="First Name"
                  />
                  <input
                    type="text"
                    value={editingUser?.lastname || ''}
                    onChange={(e) =>
                      setEditingUser({ ...editingUser!, lastname: e.target.value })
                    }
                    className="input p-2 border border-gray-300 rounded mb-2 w-full"
                    placeholder="Last Name"
                  />
                  <input
                    type="text"
                    value={editingUser?.city || ''}
                    onChange={(e) =>
                      setEditingUser({ ...editingUser!, city: e.target.value })
                    }
                    className="input p-2 border border-gray-300 rounded mb-2 w-full"
                    placeholder="City"
                  />
                </div>
                {/* Move buttons to the bottom and align to the left */}
                <div className="flex sm:col-span-5 space-x-2">
                  <button onClick={handleSaveUser} className="bg-secondary text-white py-2 px-4 rounded">
                    Save
                  </button>
                  <button
                    onClick={() => {
                      setEditingUserId(null);
                      setEditingUser(null);
                    }}
                    className="bg-gray-500 text-white py-2 px-4 rounded"
                  >
                    Cancel
                  </button>
                  <button onClick={() => handleDeleteUser(user.user_id ?? -1)} className="bg-red-500 text-white py-2 px-4 rounded">
                    Delete
                  </button>
                </div>
              </>
            ) : (
              <>
                <div className="flex-1 mb-2 sm:mb-0 sm:col-span-4">
                  <span>
                    {user.username} ({user.email})
                  </span>
                </div>
                <div className="flex sm:col-span-1 space-x-2">
                  <button onClick={() => handleEditUser(user)} className="bg-secondary text-white py-2 px-4 rounded">
                    Edit
                  </button>
                  <button onClick={() => handleDeleteUser(user.user_id ?? -1)} className="bg-red-500 text-white py-2 px-4 rounded">
                    Delete
                  </button>
                </div>
              </>
            )}
          </li>
        ))}

        {/* Add New User Row */}
        {newUser && (
          <li
            key="new-user"
            className="flex flex-col sm:grid sm:grid-cols-5 gap-2 justify-between items-center p-2 border-b border-gray-300"
          >
            <div className="flex-1 mb-2 sm:mb-0 sm:col-span-4">
              <input
                type="text"
                value={newUser.username}
                onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
                className="input p-2 border border-gray-300 rounded mb-2 w-full"
                placeholder="Username"
              />
              <input
                type="email"
                value={newUser.email}
                onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
                className="input p-2 border border-gray-300 rounded mb-2 w-full"
                placeholder="Email"
              />
              <input
                type="text"
                value={newUser.firstname}
                onChange={(e) => setNewUser({ ...newUser, firstname: e.target.value })}
                className="input p-2 border border-gray-300 rounded mb-2 w-full"
                placeholder="First Name"
              />
              <input
                type="text"
                value={newUser.lastname}
                onChange={(e) => setNewUser({ ...newUser, lastname: e.target.value })}
                className="input p-2 border border-gray-300 rounded mb-2 w-full"
                placeholder="Last Name"
              />
              <input
                type="text"
                value={newUser.city}
                onChange={(e) => setNewUser({ ...newUser, city: e.target.value })}
                className="input p-2 border border-gray-300 rounded mb-2 w-full"
                placeholder="City"
              />
              <input
                type="password"
                value={newUser.password}
                onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
                className="input p-2 border border-gray-300 rounded mb-2 w-full"
                placeholder="Password"
              />
            </div>
            <div className="flex sm:col-span-1 space-x-2">
              <button onClick={handleSaveNewUser} className="bg-secondary text-white py-2 px-4 rounded">
                Save
              </button>
              <button onClick={() => setNewUser(null)} className="bg-gray-500 text-white py-2 px-4 rounded">
                Cancel
              </button>
            </div>
          </li>
        )}
      </ul>

      {/* Add User Button at the Bottom */}
      <div className="flex space-x-4 mb-4">
        <button onClick={handleAddUserRow} className="bg-primary text-white py-2 px-4 rounded">
          Add User
        </button>
      </div>

      {/* Action Buttons at the Bottom */}
      <div className="flex space-x-4 mb-4">
        <button onClick={handleDeleteAllUsers} className="bg-red-500 text-white py-2 px-4 rounded">
          Delete All Users
        </button>
        <button onClick={handleTruncateUsersTable} className="bg-red-600 text-white py-2 px-4 rounded">
          Truncate Users Table
        </button>
      </div>
    </div>
  );
};

export default Users;
