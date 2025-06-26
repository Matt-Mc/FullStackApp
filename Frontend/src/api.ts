import type { Bill, MP } from './types';

const API_BASE_URL = 'http://localhost:1500/api';

export const fetchBills = async (): Promise<Bill[]> => {
  const response = await fetch(`${API_BASE_URL}/bills`);
  if (!response.ok) {
    throw new Error('Failed to fetch bills');
  }
  const data = await response.json();
  return data.bills;
};

export const fetchMPs = async (): Promise<MP[]> => {
  const response = await fetch(`${API_BASE_URL}/mps`);
  if (!response.ok) {
    throw new Error('Failed to fetch MPs');
  }
  const data = await response.json();
  return data.MPs;
};
