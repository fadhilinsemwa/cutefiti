// api/pay1.ts
import axios from 'axios';

const baseURL = process.env.MODERN_API_BASE_URL_LIVE_CRDB || 'http://localhost:8080';
const apiClient = axios.create({
  baseURL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interface for the request and response bodies
interface VerificationRequest {
  paymentReference: string;
  token: string;
  checksum: string;
  institutionID: string;
}

interface VerificationResponse {
  status: number;
  statusDesc: string;
  data: {
    payerName: string;
    amount: number;
    amountType: string;
    currency: string;
    paymentReference: string;
    paymentType: string;
    paymentDesc: string;
    payerID: string;
  };
}

interface PaymentPostRequest {
  payerName: string;
  amount: string;
  amountType: string;
  currency: string;
  paymentReference: string;
  paymentType: string;
  payerMobile: string;
  paymentDesc: string;
  payerID: string;
  transactionRef: string;
  transactionChannel: string;
  transactionDate: string;
  token: string;
  checksum: string;
  institutionID: string;
}

interface PaymentPostResponse {
  status: number;
  statusDesc: string;
  data: {
    receipt: string;
  };
}

// Functions to handle the requests
export const verifyPayment = async (data: VerificationRequest): Promise<VerificationResponse> => {
  const response = await apiClient.post<VerificationResponse>('/crdb/v1/verification.php', data);
  return response.data;
};

export const postPayment = async (data: PaymentPostRequest): Promise<PaymentPostResponse> => {
  const response = await apiClient.post<PaymentPostResponse>('/crdb/v1/payment-post.php', data);
  return response.data;
};
