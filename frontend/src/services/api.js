import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const register = (userData) => {
  return apiClient.post('/register', userData);
};

export const login = (userData) => {
  return apiClient.post('/login', userData);
};

export const getSpendingSummary = () => {
  return apiClient.get('/spending-summary');
};

export const generateSummary = () => {
  return apiClient.post('/generate-summary');
};
