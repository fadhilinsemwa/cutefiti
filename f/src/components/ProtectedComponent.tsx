// src/components/ProtectedComponent.tsx
import React, { useEffect, useState } from 'react';
import { useNavigate, useLocation } from '@modern-js/runtime/router';
import { getUser } from '../api/users';
import { getAllRolesComponents, getAllComponents } from '../api/roles_components';
import { getUserRoles } from '@/api/roles';
import { useAuth } from '../context/AuthContext'; // Import the AuthContext

interface ProtectedComponentProps {
  component: React.ComponentType;
  componentPath: string;
}

const ProtectedComponent: React.FC<ProtectedComponentProps> = ({ component: Component, componentPath }) => {
  const { isAuthenticated, login } = useAuth(); // Get isAuthenticated and login from AuthContext
  const [isAuthorized, setIsAuthorized] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    const checkAccess = async () => {
      try {
        const token = localStorage.getItem('token');
        const storedRoleId = localStorage.getItem('role_id');

        if (!token) {
          console.log("No token found, navigating to login.");
          navigate('#login', { state: { from: location.pathname } }); // Navigate to the login page using hash
          return;
        }

        // Check if role_id is in localStorage
        let role_id = storedRoleId;

        if (!role_id) {
          // If not found in localStorage, fetch the user's roles from API
          const userRoles = await getUserRoles();
          console.log('User Roles:', userRoles);

          const rolesArray = Array.isArray(userRoles) ? userRoles : Object.values(userRoles);

          if (!Array.isArray(rolesArray)) {
            throw new Error('User roles data is not in a valid format');
          }

          // Fetch user details (assuming the user ID is stored or fetched)
          const user = await getUser(1); // Replace '1' with the correct user ID
          console.log("User details:", user);

          // Find the role for the logged-in user
          const userRole = rolesArray.find((role: any) => role.user_id === user.user_id);

          if (!userRole || !userRole.role_id) {
            console.log("User role not found, navigating to unauthorized.");
            navigate('#unauthorized');
            return;
          }

          role_id = userRole.role_id;
        }

        if (!role_id) {
          console.log("Role ID is missing, navigating to unauthorized.");
          setIsAuthorized(false);
          navigate('#unauthorized');
          return;
        }

        // Fetch all role-component mappings
        const rolesComponents = await getAllRolesComponents();
        console.log("Roles-Components mapping:", rolesComponents.data);

        // Filter components assigned to this role
        const componentsForRole = rolesComponents.data.filter(
          (roleComp: any) => roleComp.role_id === role_id
        );

        console.log(`Components assigned to role ${role_id}:`, componentsForRole);

        // Get all components
        const allComponents = await getAllComponents();
        console.log("All components:", allComponents.data); // Fix: Access 'data' property of response

        // Check if the component is assigned to this role
        const isComponentAssigned = componentsForRole.some((roleComp: any) => {
          // Match the component by ID
          const matchedComponent = allComponents.data.find( // Fix: Access 'data' property of response
            (comp: any) => comp.id === roleComp.component_id
            
          );
          console.log(matchedComponent);
          // If a match is found, check if the component path matches
          return matchedComponent && matchedComponent.component_path === componentPath;
          
        });
        console.log(componentPath);
        // If the component is assigned or no restrictions exist, authorize access
        if (isComponentAssigned || componentsForRole.length === 0) {
          setIsAuthorized(true);
          console.log("Access granted to component:", componentPath);
        } else {
          console.log("Access denied, navigating to unauthorized.");
          setIsAuthorized(false);
          navigate('#unauthorized');
        }
      } catch (error) {
        console.error('Error checking access:', error);
        navigate('#error');
      } finally {
        setLoading(false);
      }
    };

    checkAccess();
  }, [componentPath, navigate, location.pathname]);

  useEffect(() => {
    if (isAuthenticated) {
      login(); // Call login when authenticated
    }
  }, [isAuthenticated, login]);

  if (loading) return <div>Loading...</div>;

  return isAuthorized ? <Component /> : null;
};

export default ProtectedComponent;
