import { createSlice, createAsyncThunk,  } from '@reduxjs/toolkit';
import { loginApi } from '../lib/api/loginApi';
import type { loginResponse } from '../schemas/userSchema';


interface AuthState {
  user: loginResponse['user'] | null // ข้อมูล user หลัง login
  status: 'idle' | 'loading' | 'succeeded' | 'failed'; // async process สถานะการเรียก API
  error: string | null;
}

const initialState: AuthState = { // state เริ่มต้น
  user: null,
  status: 'idle',
  error: null,
};

// ส่งข้อมูลให้ Backend
export const loginUser = createAsyncThunk<
  loginResponse,  // รูปแบบข้อมูลที่คืนกลับ
  { email: string; password: string }, // รูปแบบ arguments ที่รับเข้า
  { rejectValue: string } 
>(
  '/core/auth/login',
  async (credentials, { rejectWithValue }) => {
    try {
     return await loginApi(credentials)
    } catch (err: any) {
        // ส่งต่อ error จาก backend
      return rejectWithValue(err.response?.data?.message || 'Login failed');
    }
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
     // action สำหรับ logout เคลียร์ข้อมูลผู้ใช้และสถานะต่าง ๆ
    logout(state) {
      state.user = null;
      state.status = 'idle';
      state.error = null;
    },
    // action สำหรับล้าง error
    clearError(state) {
      state.error = null;
    },
  },
  extraReducers: builder => {
    builder
      // รอ login
      .addCase(loginUser.pending, state => {
        state.status = 'loading';
        state.error = null;
      })
      // login เสร็จ
      .addCase(loginUser.fulfilled, (state, action) => {
        state.status = 'succeeded';
        state.user = action.payload.user; // เก็บข้อมูลผู้ใช้ที่ได้รับกลับมา
      })
      // login ไม่ผ่าน
      .addCase(loginUser.rejected, (state, action) => {
        state.status = 'failed';
        // เก็บ error
        state.error = action.payload ?? action.error.message ?? null;
      });
  },
});

// action สำหรับเอาไปใช้ใน commpent
export const { logout, clearError } = authSlice.actions;
// reducer ไว้ต่อกับ store
export default authSlice.reducer;
