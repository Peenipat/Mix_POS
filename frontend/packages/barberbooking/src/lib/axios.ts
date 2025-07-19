import axios from "axios";
// import { store } from "@/store"        // path ไปยัง Redux store
// import { RootState } from "@/store";

const instance = axios.create({
  baseURL: `${import.meta.env.VITE_API_BASE_URL}`,
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
