// src/api/users.ts
const MODERN_API_BASE_URL = process.env.MODERN_API_BASE_URL || 'http://localhost:8001';

export interface User {
  user_id?: number; // Include user_id in the User interface
  username: string;
  email: string;
  firstname: string;
  lastname: string;
  city: string;
  password: string;
}

// Helper function to check response status and handle errors
async function handleResponse(response: Response) {
  if (!response.ok) {
    const errorMessage = await response.text();
    throw new Error(`HTTP error! status: ${response.status}, message: ${errorMessage}`);
  }
  return await response.json(); // Directly parse JSON response
}

// Create a new user
export async function createUser(user: User): Promise<User> {
  const response = await fetch(`${MODERN_API_BASE_URL}/user`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user),
  });
  return await handleResponse(response);
}

// Get a user by ID
export async function getUser(id: number): Promise<User> {
  const response = await fetch(`${MODERN_API_BASE_URL}/user?id=${id}`);
  return await handleResponse(response);
}

// Get all users
export async function getAllUsers(): Promise<User[]> {
  const response = await fetch(`${MODERN_API_BASE_URL}/users`);
  return await handleResponse(response);
}

// Update a user
export async function updateUser(id: number, user: User): Promise<User> {
  const response = await fetch(`${MODERN_API_BASE_URL}/user?id=${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(user),
  });
  return await handleResponse(response);
}

// Delete a user by ID
export async function deleteUser(id: number): Promise<void> {
  const response = await fetch(`${MODERN_API_BASE_URL}/user?id=${id}`, {
    method: 'DELETE',
  });
  await handleResponse(response);
}

// Delete all users
export async function deleteAllUsers(): Promise<void> {
  const response = await fetch(`${MODERN_API_BASE_URL}/delete-all-users`, {
    method: 'DELETE',
  });
  await handleResponse(response);
}

// Truncate users table
export async function truncateUsersTable(): Promise<void> {
  const response = await fetch(`${MODERN_API_BASE_URL}/truncate-users`, {
    method: 'POST',
  });
  await handleResponse(response);
}
