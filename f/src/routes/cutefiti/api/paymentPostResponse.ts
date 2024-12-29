import axios from 'axios';

// Base URL for API
const baseURL = process.env.MODERN_API_BASE_URL2 || 'http://localhost:8002';

// PaymentPostResponse interface representing the structure of payment post response
export interface PaymentPostResponse {
  pprno: number;
  applicant_id: number;
  receipt: string | null;
  oreceipt: string | null;
  status: string | null;
  status_desc: string | null;
  muda: string | null;
  student_id: number | null;
  pay_type: string | null;
  payer_name: string | null;
  invoice_no: string | null;
  bank_ac: string | null;
  amount_tzs: number | null;
  balance_tzs: number | null;
  aina: string | null;
  payment_desc: string | null;
  channeli: string | null;
  transaction_no: string | null;
}

// Function to fetch all payment post responses
export const fetchPaymentPostResponses = async (): Promise<PaymentPostResponse[]> => {
  try {
    const response = await axios.get<PaymentPostResponse[]>(`${baseURL}/paymentpostres`);
    return response.data;
  } catch (error) {
    console.error('Error fetching payment post responses:', error);
    throw error; // Rethrow the error for further handling
  }
};
