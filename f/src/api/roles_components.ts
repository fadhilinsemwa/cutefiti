// usersf/src/api/roles_components.ts
import axios from 'axios';

const MODERN_API_BASE_URL = process.env.MODERN_API_BASE_URL || 'http://localhost:8001';

const api = axios.create({
  baseURL: MODERN_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const createComponent = async (data: { component_name: string, component_path: string }) => {
  return await api.post('/components', data);
};

export const getAllComponents = async () => {
  return await api.get('/components');
};

export const updateComponent = async (id: number, data: { component_name: string, component_path: string }) => {
  return await api.put(`/components/${id}`, data);
};

export const deleteComponent = async (id: number) => {
  return await api.delete(`/components/${id}`);
};

export const deleteAllComponents = async () => {
  return await api.post('/components/delete_all');
};

export const truncateComponents = async () => {
  return await api.post('/components/truncate');
};

export const assignRoleToComponent = async (data: { role_id: number, component_id: number }) => {
  return await api.post('/roles_components', data);
};

export const getAllRolesComponents = async () => {
  return await api.get('/roles_components');
};

export const updateRoleComponent = async (id: number, data: { role_id: number, component_id: number }) => {
  return await api.put(`/roles_components/${id}`, data);
};

export const deleteRoleComponent = async (id: number) => {
  return await api.delete(`/roles_components/${id}`);
};

export const deleteAllRoleComponents = async () => {
  return await api.post('/roles_components/delete_all');
};

export const truncateRolesComponents = async () => {
  return await api.post('/roles_components/truncate');
};
