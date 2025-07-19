// src/store/authSlice.ts
import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { loginApi } from '../lib/api/loginApi';
import { getMeApi } from '../api/getMe';
import type { loginResponse } from '../schemas/userSchema';
import type { Me } from '../api/getMe';

export const loginUser = createAsyncThunk<
  loginResponse,
  { email: string; password: string },
  { rejectValue: string }
>(
  '/core/auth/login',
  async (credentials, { rejectWithValue }) => {
    try {
      return await loginApi(credentials);
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.message || 'Login failed');
    }
  }
);

export const loadCurrentUser = createAsyncThunk<
  Me,
  void,
  { rejectValue: string }
>(
  '/core/auth/me',
  async (_, { rejectWithValue }) => {
    try {
      return await getMeApi();
    } catch (err: any) {
      return rejectWithValue(err.response?.data?.message || 'Not authenticated');
    }
  }
);

export interface AuthState {
  loginUser: loginResponse['user'] | null; // ข้อมูลเบื้องต้นจาก /login
  me: Me | null;                            // ข้อมูลโปรไฟล์เต็มจาก /me
  statusLogin: 'idle' | 'loading' | 'succeeded' | 'failed';
  statusMe: 'idle' | 'loading' | 'succeeded' | 'failed';
  errorLogin: string | null;
  errorMe: string | null;
}

const initialState: AuthState = {
  loginUser: null,
  me: null,
  statusLogin: 'idle',
  statusMe: 'idle',
  errorLogin: null,
  errorMe: null,
};

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    logout(state) {
      state.loginUser = null;
      state.me = null;
      state.statusLogin = 'idle';
      state.statusMe = 'idle';
      state.errorLogin = null;
      state.errorMe = null;
    },
  },
  extraReducers: (builder) => {
    // ===== loginUser =====
    builder.addCase(loginUser.pending, (state) => {
      state.statusLogin = 'loading';
      state.errorLogin = null;
    });
    builder.addCase(loginUser.fulfilled, (state, action) => {
      state.statusLogin = 'succeeded';
      state.loginUser = action.payload.user;
    });
    builder.addCase(loginUser.rejected, (state, action) => {
      state.statusLogin = 'failed';
      state.errorLogin = action.payload ?? action.error.message ?? null;
      state.loginUser = null;
    });

    // ===== loadCurrentUser =====
    builder.addCase(loadCurrentUser.pending, (state) => {
      state.statusMe = 'loading';
      state.errorMe = null;
    });
    builder.addCase(loadCurrentUser.fulfilled, (state, action) => {
      state.statusMe = 'succeeded';
      state.me = action.payload;
    });
    builder.addCase(loadCurrentUser.rejected, (state, action) => {
      state.statusMe = 'failed';
      state.errorMe = action.payload ?? action.error.message ?? null;
      state.me = null;
    });
  },
});

export const { logout } = authSlice.actions;
export default authSlice.reducer;
