// src/components/Logout.tsx

import { useNavigate } from '@modern-js/runtime/router';

const Logout = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    // Clear session data
    localStorage.removeItem('token');
    localStorage.removeItem('role_id');

    // Redirect to login or dashboard
    navigate('/login', { replace: true });
  };

  return (
    <button onClick={handleLogout} className="btn-danger">
      Logout
    </button>
  );
};

export default Logout;
