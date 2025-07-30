import React, { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { loginSchema } from '../schemas/authSchema';
import type { LoginForm } from '../schemas/authSchema';
import { useNavigate } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from '../store/hook';
import { loginUser, loadCurrentUser, logout } from '../store/authSlice';
import { navigateByRole } from '../utils/navigation';
import Modal from '@object/shared/components/Modal';
import { FiInfo } from 'react-icons/fi';

export default function Login() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();

  const statusLogin = useAppSelector((state) => state.auth.statusLogin);
  const errorLogin = useAppSelector((state) => state.auth.errorLogin);
  const loginUserData = useAppSelector((state) => state.auth.loginUser);
  const statusMe = useAppSelector((state) => state.auth.statusMe);
  const me = useAppSelector((state) => state.auth.me);
  const errorMe = useAppSelector((state) => state.auth.errorMe);

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isRegisterForm, setIsRegisterForm] = useState(false)

  const handleClose = () => {
    setIsModalOpen(false)
  }
  const handleOpen = () => {
    // setIsModalOpen(true)
  }

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<LoginForm>({ resolver: zodResolver(loginSchema) });

  const onSubmit = (data: LoginForm) => {
    dispatch(loginUser(data));
  };

  // 1) เมื่อ login สำเร็จ → fetch /me
  useEffect(() => {
    if (statusLogin === 'succeeded' && loginUserData) {
      dispatch(loadCurrentUser());
    }
  }, [statusLogin, loginUserData, dispatch]);

  // 2) เมื่อ fetch /me สำเร็จ → redirect ตาม role
  useEffect(() => {
    if (statusMe === 'succeeded' && me) {
      navigateByRole(me.role, navigate);
    }
  }, [statusMe, me, navigate]);

  // 3) error handling
  useEffect(() => {
    if (statusLogin === 'failed' && errorLogin) {
      alert(errorLogin);
      dispatch(logout());
    }
    if (statusMe === 'failed' && errorMe) {
      alert(errorMe);
      dispatch(logout());
    }
  }, [statusLogin, errorLogin, statusMe, errorMe, dispatch]);

  const loginAs = (role: 'admin' | 'barber') => {
    const account = {
      email: import.meta.env[`VITE_${role.toUpperCase()}_EMAIL`],
      password: import.meta.env[`VITE_${role.toUpperCase()}_PASSWORD`],
    };

    setValue('email', account.email || '');
    setValue('password', account.password || '');

    handleSubmit(onSubmit)(); 
  };


  return (
    <div className="grid grid-cols-2 min-h-screen mx-auto">
      {/* ฝั่งซ้าย: Login Form */}
      <div className="flex items-center justify-center bg-gray-100 px-12">
        <div className="w-full max-w-md">
          <h1 className='text-red-500 text-2xl text-center'>ตอนนี้ยังเป็น Version ที่ยังไม่สมบูรณ์บาง function อาจจะใช้งานไม่ได้</h1>
          <h1 className="text-5xl font-bold mb-8 text-center text-gray-800">Login</h1>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <input
              type="email"
              placeholder="example@gmail.com"
              className="input input-bordered w-full"
              {...register('email')}
            />
            {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}

            <input
              type="password"
              placeholder="Password"
              className="input input-bordered w-full"
              {...register('password')}
            />
            {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}

            <div className="grid grid-cols-2 gap-4 pt-2">
              <button
                className="btn btn-success w-full rounded-lg"
                type="submit"
                disabled={statusLogin === 'loading' || statusMe === 'loading'}
              >
                {(statusLogin === 'loading' || statusMe === 'loading') ? 'Loading…' : 'Login'}
              </button>
              <button
                className="btn btn-warning w-full rounded-lg"
                type="button"
                onClick={() => handleOpen()}
              >
                Register
              </button>
            </div>
          </form>

          <div className='flex gap-5 justify-center mt-5'>
          <button type="button" className='p-2 bg-green-500 text-white rounded-md' onClick={() => navigate("/")}>
             กลับหน้าหลัก
            </button>
            <button type="button" className='p-2 bg-green-500 text-white rounded-md' onClick={() => loginAs('barber')}>
              ทดสอบในมุม ช่าง
            </button>

            <button type="button" className='p-2 bg-green-500 text-white rounded-md' onClick={() => loginAs('admin')}>
              ทดสอบในมุม แอดมิน
            </button>
          </div>

        </div>
      </div>

      <div className="relative overflow-hidden">
        <img src="" alt="Login" className="w-full h-full object-cover" />
        <div className="absolute bottom-0 left-0 w-full flex justify-center">
          <div className="bg-white/30 backdrop-blur-md w-full p-3 shadow-inner" />
          <img src="" alt="dde" />
        </div>
      </div>


      <Modal isOpen={isModalOpen} onClose={handleClose} title={isRegisterForm ? "ลงทะเบียน" : "เข้าสู่ระบบ"} blurBackground>
        <div className='p-12 pt-0'>
          {isRegisterForm ? (
            <form className="space-y-5">
              <div>
                <label className="flex items-end text-sm font-medium text-gray-700 mb-1">
                  ชื่อลูกค้า
                  {/* <CustomTooltip
                  id="tooltip-cusname"
                  content="แนะนำเป็นภาษาไทย"
                  trigger="hover"
                  placement="top"
                  bgColor="bg-gray-200"
                  textColor="text-gray-900"
                  textSize="text-sm"
                  className="ml-1"
                >
                  <span><FiInfo /></span>
                </CustomTooltip> */}
                </label>
                <input
                  type="text"
                  placeholder="กรุณากรอกชื่อ"
                  className="input input-bordered w-full"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">เบอร์โทร</label>
                <input
                  type="text"
                  placeholder="กรุณากรอกเบอร์โทรศัพท์"
                  className="input input-bordered w-full"
                />
              </div>

              <div>
                <label className="flex items-end text-sm font-medium text-gray-700 mb-1 ">
                  รหัสผ่าน
                  {/* <CustomTooltip
                  id="tooltip-password"
                  content="รหัสผ่านสำหรับเข้าใช้งานครั้งถัดไป"
                  trigger="hover"
                  placement="top"
                  bgColor="bg-gray-200"
                  textColor="text-gray-900"
                  textSize="text-sm"
                  className="ml-1"
                >
                  <span><FiInfo /></span>
                </CustomTooltip> */}
                </label>
                <input
                  type="password"
                  placeholder="ตั้งรหัสผ่าน"
                  className="input input-bordered w-full"
                />
              </div>

              <div className="flex gap-4 pt-4">
                <button type="submit" className="w-1/2 bg-green-600 hover:bg-green-700 text-white py-2 rounded">
                  ลงทะเบียน
                </button>
                <button
                  type="button"
                  onClick={() => setIsRegisterForm(false)} // สลับไป login
                  className="w-1/2 bg-gray-400 hover:bg-gray-500 text-white py-2 rounded"
                >
                  มีบัญชีอยู่แล้ว?
                </button>
              </div>
            </form>
          ) : (
            <form className="space-y-5">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">เบอร์โทร</label>
                <input
                  type="text"
                  placeholder="กรุณากรอกเบอร์โทรศัพท์"
                  className="input input-bordered w-full"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">รหัสผ่าน</label>
                <input
                  type="password"
                  placeholder="กรุณากรอกรหัสผ่าน"
                  className="input input-bordered w-full"
                />
              </div>

              <div className="flex gap-4 pt-4">
                <button type="submit" className="w-1/2 bg-blue-600 hover:bg-blue-700 text-white py-2 rounded">
                  เข้าสู่ระบบ
                </button>
                <button
                  type="button"
                  onClick={() => setIsRegisterForm(true)} // สลับกลับไป register
                  className="w-1/2 bg-gray-400 hover:bg-gray-500 text-white py-2 rounded"
                >
                  สมัครสมาชิก
                </button>
              </div>
            </form>
          )}
        </div>
      </Modal>
    </div>
  );
}


