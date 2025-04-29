import axios from "axios";

const instance = axios.create({
  baseURL: "http://localhost:3001/",
  withCredentials: true,
  headers: { 'Content-Type': 'application/json' },
});

// Interceptor ดัก request ใส่ token ทุกครั้ง
instance.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  console.log("👉 กำลังส่ง token นี้:", token);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor ดัก response ถ้าเจอ 401 Unauthorized
instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      window.location.href = "/"; // เด้งกลับหน้า home
    }
    return Promise.reject(error);
  }
);

export default instance;
