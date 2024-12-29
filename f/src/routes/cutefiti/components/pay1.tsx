// components/pay1.tsx
import React, { useState } from 'react';
import { verifyPayment, postPayment } from '../api/pay1';
import ErrorResponse from './ErrorResponse';

const Pay1Component: React.FC = () => {
  const initialVerificationData = {
    paymentReference: 'MU82',
    token: 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJtd2FuemF1bml2ZXJzaXR5IiwiaWF0IjoxNjM0MDIxNTU1fQ.CpmghOT2tMX3L6GCZ-rpGS10A2xoBNUMPKKEsAdUxAQ',
    checksum: 'eec8eaf4b9e22eac040cd922338f2b80ead4688a',
    institutionID: '60891604930813',
  };
  const initialPaymentPostData = {
    payerName: 'HENRY, James (M)',
    amount: '500',
    amountType: 'FLEXIBLE',
    currency: 'TZS',
    paymentReference: 'MU82',
    paymentType: '50',
    payerMobile: '0787000000',
    paymentDesc: 'Tuition Fee',
    payerID: 'E300/90',
    transactionRef: 'Fuyet355534233550tl1',
    transactionChannel: 'IB',
    transactionDate: '2017-06-30 14:48:42',
    token: 'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJtd2FuemF1bml2ZXJzaXR5IiwiaWF0IjoxNjM0MDIxNTU1fQ.CpmghOT2tMX3L6GCZ-rpGS10A2xoBNUMPKKEsAdUxAQ',
    checksum: 'eec8eaf4b9e22eac040cd922338f2b80ead4688a',
    institutionID: '8143',
  };

  const [verificationData, setVerificationData] = useState(initialVerificationData);
  const [paymentPostData, setPaymentPostData] = useState(initialPaymentPostData);
  const [showVerificationFields, setShowVerificationFields] = useState(false);
  const [editVerification, setEditVerification] = useState(false);
  const [showPaymentFields, setShowPaymentFields] = useState(false);
  const [editPayment, setEditPayment] = useState(false);
  const [verificationResponse, setVerificationResponse] = useState<any>(null);
  const [paymentResponse, setPaymentResponse] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);
  const [errorInfo, setErrorInfo] = useState<{ status?: number; statusDesc?: string; details?: string } | undefined>(undefined);
  const [loading, setLoading] = useState({ verify: false, post: false });

  const handleVerifyPayment = async () => {
    setLoading((prev) => ({ ...prev, verify: true }));
    try {
      const response = await verifyPayment(verificationData);
      setVerificationResponse(response);
      setError(null);
      setErrorInfo(undefined);
    } catch (err: any) {
      console.error('Verification Error:', err);
      setError('Verification Error');
      setErrorInfo({
        status: err.response?.data?.status,
        statusDesc: err.response?.data?.statusDesc,
        details: err.response?.data?.message,
      });
    } finally {
      setLoading((prev) => ({ ...prev, verify: false }));
    }
  };

  const handlePostPayment = async () => {
    setLoading((prev) => ({ ...prev, post: true }));
    try {
      const response = await postPayment(paymentPostData);
      setPaymentResponse(response);
      setError(null);
      setErrorInfo(undefined);
    } catch (err: any) {
      console.error('Payment Post Error:', err);
      setError('Payment Post Error');
      setErrorInfo({
        status: err.response?.data?.status,
        statusDesc: err.response?.data?.statusDesc,
        details: err.response?.data?.message,
      });
    } finally {
      setLoading((prev) => ({ ...prev, post: false }));
    }
  };

  return (
    <div className="p-4 max-w-lg mx-auto bg-gray-50 rounded-lg shadow-md">
      <h1 className="text-2xl md:text-3xl font-bold mb-4 text-center text-gray-800">CRDB Payment Testing</h1>

      {/* Verification Section */}
      <div className="mb-6">
        <h2 className="text-lg md:text-xl font-semibold mb-2 text-gray-700">Verify Payment</h2>
        <p className="text-xs md:text-sm text-gray-600">Reference: {verificationData.paymentReference}</p>
        <div className="flex gap-2 mt-2">
          <button
            className="flex-1 bg-primary text-white py-1 md:py-2 rounded text-sm md:text-base hover:bg-primary-dark transition duration-200"
            onClick={() => setShowVerificationFields(!showVerificationFields)}
          >
            Show
          </button>
          <button
            className="flex-1 bg-secondary text-white py-1 md:py-2 rounded text-sm md:text-base hover:bg-secondary-dark transition duration-200"
            onClick={() => setEditVerification(true)}
          >
            Edit
          </button>
          <button
            className="flex-1 bg-blue-500 text-white py-1 md:py-2 rounded text-sm md:text-base hover:bg-blue-600 transition duration-200"
            onClick={() => {
              handleVerifyPayment();
              setEditVerification(false);
            }}
            disabled={loading.verify}
          >
            {loading.verify ? 'Loading...' : 'Verify'}
          </button>
        </div>

        {showVerificationFields && (
          <div className="mt-4">
            {Object.keys(initialVerificationData).map((key) => (
              <div key={key} className="mb-2">
                <label className="block font-semibold text-xs md:text-sm text-gray-700">{key}</label>
                {editVerification ? (
                  <input
                    className="border p-1 md:p-2 w-full rounded text-xs md:text-sm"
                    type="text"
                    value={verificationData[key as keyof typeof verificationData]}
                    onChange={(e) => setVerificationData({ ...verificationData, [key]: e.target.value })}
                  />
                ) : (
                  <p className="text-gray-600 overflow-x-auto whitespace-pre-line break-words text-xs md:text-sm">
                    {verificationData[key as keyof typeof verificationData]}
                  </p>
                )}
              </div>
            ))}
          </div>
        )}
        {verificationResponse && (
          <pre className="mt-4 bg-gray-200 p-2 rounded overflow-auto text-xs md:text-sm">{JSON.stringify(verificationResponse, null, 2)}</pre>
        )}
      </div>

      {/* Payment Section */}
      <div>
        <h2 className="text-lg md:text-xl font-semibold mb-2 text-gray-700">Post Payment</h2>
        <p className="text-xs md:text-sm text-gray-600">Reference: {paymentPostData.paymentReference}</p>
        <div className="flex gap-2 mt-2">
          <button
            className="flex-1 bg-primary text-white py-1 md:py-2 rounded text-sm md:text-base hover:bg-primary-dark transition duration-200"
            onClick={() => setShowPaymentFields(!showPaymentFields)}
          >
            Show
          </button>
          <button
            className="flex-1 bg-secondary text-white py-1 md:py-2 rounded text-sm md:text-base hover:bg-secondary-dark transition duration-200"
            onClick={() => setEditPayment(true)}
          >
            Edit
          </button>
          <button
            className="flex-1 bg-green-500 text-white py-1 md:py-2 rounded text-sm md:text-base hover:bg-green-600 transition duration-200"
            onClick={() => {
              handlePostPayment();
              setEditPayment(false);
            }}
            disabled={loading.post}
          >
            {loading.post ? 'Loading...' : 'Post Payment'}
          </button>
        </div>

        {showPaymentFields && (
          <div className="mt-4">
            {Object.keys(initialPaymentPostData).map((key) => (
              <div key={key} className="mb-2">
                <label className="block font-semibold text-xs md:text-sm text-gray-700">{key}</label>
                {editPayment ? (
                  <input
                    className="border p-1 md:p-2 w-full rounded text-xs md:text-sm"
                    type="text"
                    value={paymentPostData[key as keyof typeof paymentPostData]}
                    onChange={(e) => setPaymentPostData({ ...paymentPostData, [key]: e.target.value })}
                  />
                ) : (
                  <p className="text-gray-600 overflow-x-auto whitespace-pre-line break-words text-xs md:text-sm">
                    {paymentPostData[key as keyof typeof paymentPostData]}
                  </p>
                )}
              </div>
            ))}
          </div>
        )}
        {paymentResponse && (
          <pre className="mt-4 bg-gray-200 p-2 rounded overflow-auto text-xs md:text-sm">{JSON.stringify(paymentResponse, null, 2)}</pre>
        )}
      </div>

      {/* Error Modal */}
      {error && (
        <ErrorResponse
          message={error}
          additionalInfo={errorInfo}
          onClose={() => {
            setError(null);
            setErrorInfo(undefined);
          }}
        />
      )}
    </div>
  );
};

export default Pay1Component;
