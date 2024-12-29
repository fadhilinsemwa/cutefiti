// src/api/authApi.ts

const BASE_URL = process.env.MODERN_API_BASE_URL || 'http://localhost:8001';

const api = async (endpoint: string, method: string, body?: any) => {
  const response = await fetch(`${BASE_URL}/${endpoint}`, {
    method,
    headers: {
      'Content-Type': 'application/json', // Change to application/json if your backend expects JSON
    },
    body: body ? JSON.stringify(body) : undefined, // Use JSON.stringify if sending JSON
  });

  if (!response.ok) {
    throw new Error(`Error: ${response.statusText}`);
  }

  return response.json();
};

export const registerUser = (userData: any) => api('register', 'POST', userData);

export const loginUser = (userData: any) => api('login', 'POST', userData);

export const logoutUser = () => api('logout', 'POST'); // Added logout endpoint

export const requestPasswordReset = (email: string) => api('request-password-reset', 'POST', { email });

export const resetPassword = (token: string, newPassword: string) => api('reset-password', 'POST', { reset_token: token, new_password: newPassword });

export const verifyEmail = (token: string) => api('verify-email', 'POST', { token });
