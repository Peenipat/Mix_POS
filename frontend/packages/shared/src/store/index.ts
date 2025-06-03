import { configureStore } from '@reduxjs/toolkit';
import authReducer from './authSlice';

//สร้าง Redux store ด้วย configureStore ของ Redux Toolkit
export const store = configureStore({
  reducer: {
    auth: authReducer, // sclice ที่จะเก็บใน store
    // เพิ่ม slice อื่น ๆ ที่นี่
  },
   // เปิดใช้ Redux DevTools เฉพาะในโหมด dev
  devTools: import.meta.env.DEV,
});

//type ของ state ทั้งหมดใน store
export type RootState = ReturnType<typeof store.getState>;
//type ของ dispatch function เพื่อให้ useAppDispatch รู้จัก action
export type AppDispatch = typeof store.dispatch;
