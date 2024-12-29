// src/components/ErrorResponse.tsx
import React from 'react';
import Modal from './Modal';

interface ErrorResponseProps {
  message: string;
  additionalInfo?: {
    status?: number;
    statusDesc?: string;
    details?: string;
  };
  onClose: () => void;
}

const ErrorResponse: React.FC<ErrorResponseProps> = ({ message, additionalInfo, onClose }) => {
  return (
    <Modal isOpen={true} onClose={onClose}>
      <h3 className="text-lg font-bold text-red-500 mb-4">Error</h3>
      <p>{message}</p>
      {additionalInfo?.statusDesc && (
        <p className="text-red-400 mt-2">
          <strong>Status Description:</strong> {additionalInfo.statusDesc}
        </p>
      )}
      {additionalInfo?.details && (
        <p className="text-red-400 mt-2">
          <strong>Details:</strong> {additionalInfo.details}
        </p>
      )}
    </Modal>
  );
};

export default ErrorResponse;
