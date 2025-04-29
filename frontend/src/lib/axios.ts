import axios from "axios";

const instance = axios.create({
  baseURL: "http://localhost:3001/",
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

// Interceptor ดัก response ถ้าเจอ 401 Unauthorized
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      window.location.href = "/"; // เด้งกลับหน้า home
    }
    return Promise.reject(error);
  }
);

export default instance;
