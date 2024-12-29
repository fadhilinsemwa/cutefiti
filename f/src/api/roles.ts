import axios from 'axios';

const API_BASE_URL = process.env.MODERN_API_BASE_URL3 || 'http://localhost:8001/api';

// Create a role
export const createRole = async (role: { role_id: string; role_name: string }) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/roles`, role);
    return response.data;
  } catch (error) {
    console.error("Error creating role:", error);
    throw error;
  }
};

// Get all roles
export const getAllRoles = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/roles`);
    return response.data;
  } catch (error) {
    console.error("Error fetching roles:", error);
    throw error;
  }
};

// Assign role to user
export const assignUserRole = async (userRole: { username: string; role_id: string }) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/assign-role`, userRole);
    return response.data;
  } catch (error) {
    console.error("Error assigning role to user:", error);
    throw error;
  }
};

// Get all user roles
export const getUserRoles = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/user-roles`);
    return response.data;
  } catch (error) {
    console.error("Error fetching user roles:", error);
    throw error;
  }
};
