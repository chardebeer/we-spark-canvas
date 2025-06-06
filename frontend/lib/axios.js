// frontend/lib/axios.js
import axios from "axios";

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  headers: {
    "Content-Type": "application/json",
  },
  timeout: 10000, // 10s timeout
});

// Interceptors for authentication
apiClient.interceptors.request.use(
  (config) => {
    // Only add token in browser environment
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem("token");
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for handling auth errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Redirect to login on 401 unauthorized errors
    if (error.response && error.response.status === 401 && typeof window !== 'undefined') {
      // Clear any existing auth data
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // Get the current path to redirect back after login
      const currentPath = window.location.pathname;
      if (currentPath !== '/login') {
        window.location.href = `/login?redirect=${currentPath}`;
      }
    }
    return Promise.reject(error);
  }
);

export default apiClient;
