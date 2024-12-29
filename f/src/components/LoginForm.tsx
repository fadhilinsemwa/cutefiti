// src/components/LoginForm.tsx
import { useState } from 'react';
import { useNavigate, useLocation } from '@modern-js/runtime/router';
import { loginUser } from '../api/authApi';
import { getUser } from '../api/users';
import { getUserRoles } from '../api/roles';

const LoginForm = () => {
  const [formData, setFormData] = useState({ username: '', password: '' });
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const location = useLocation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await loginUser(formData);
      console.log('Login response:', response);

      if (!response.token || !response.user_id) {
        throw new Error('Login response is missing token or user_id');
      }

      localStorage.setItem('token', response.token);
      const user = await getUser(response.user_id);
      console.log('User details:', user);

      const userRoles = await getUserRoles();
      console.log('User roles:', userRoles);

      const role = userRoles[user.username];
      console.log('Selected role:', role);

      if (role && role.role_id) {
        localStorage.setItem('role_id', String(role.role_id));
      } else {
        throw new Error('Role ID not found');
      }

      const redirectPath = location.state?.from || '/';
      setError(null);
      navigate(redirectPath);
    } catch (err) {
      setError('Failed to login');
      console.error('Login error:', err);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-4">
      <h2 className="text-xl mb-4">Login</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Username"
          value={formData.username}
          onChange={(e) => setFormData({ ...formData, username: e.target.value })}
          className="input w-full p-2 mb-4 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <input
          type="password"
          placeholder="Password"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          className="input w-full p-2 mb-6 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>
      <button
        type="submit"
        className="btn-primary w-full p-3 text-white bg-primary-600 rounded hover:bg-primary-500 transition duration-200"
      >
        Login
      </button>
    </form>
  );
};

export default LoginForm;
