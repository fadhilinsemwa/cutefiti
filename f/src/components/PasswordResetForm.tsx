// src/components/PasswordResetForm.tsx

import { useState } from 'react';
import { requestPasswordReset } from '../api/authApi';

const PasswordResetForm = () => {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await requestPasswordReset(email);
      setMessage('Password reset link has been sent to your email');
    } catch (err) {
      setMessage('Failed to send password reset link');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="p-4">
      <h2 className="text-xl mb-4">Reset Password</h2>
      {message && <p>{message}</p>}
      <input
        type="email"
        placeholder="Enter your email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        className="input mb-4"
      />
      <button type="submit" className="btn-primary">Send Reset Link</button>
    </form>
  );
};

export default PasswordResetForm;
