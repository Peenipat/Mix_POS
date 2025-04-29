import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema, LoginForm } from "../schemas/authSchema";
import { useNavigate } from "react-router-dom"
import { useEffect} from "react";
import { useAppDispatch, useAppSelector } from '../store/hook';
import { loginUser, clearError } from '@/store/authSlice';
import { useTypewriter } from '../components/useTypewriter';
import { navigateByRole } from "@/utils/navigation";
export default function Login() {
  const dispatch = useAppDispatch();
  // content ของ animation ตัวอักษร
  const messages = [
    "Lorem ipsum dolor sit amet  accusamus non! Error voluptatibus dignissimos magnam ",
    "Lorem  accusantium et, solutais deleniti harum ex non. Magni, earum. Cupiditate?",
    "Lorem ipsum dolor sit amet consectetur adipisicing elit.  Error voluptCupiditate?",
  ];
  // ดึงสถานะมาจาก Redux
  const { status, error, user } = useAppSelector(state => state.auth);
  // ตัวจัดการ check ข้อมูลที่เข้ามาใน form
  const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });
  const navigate = useNavigate();

  const onSubmit = (data: LoginForm) => {
    dispatch(loginUser(data)); //ส่งข้อมูล login ไปยัง Redux (loginUser thunk)
  };

  // เมื่อ login สำเร็จ (status === 'succeeded') ให้เก็บ token+navigate
  useEffect(() => {
    if (status === 'succeeded' && user){
      console.log("🚀 user after login: ", user);
      navigateByRole(user.role, navigate);
      setTimeout(() => {
        navigateByRole(user.role, navigate);
      }, 200);
    }
  }, [status, user, navigate]);

  // ถ้ามี error ให้ alert หรือแสดง UI
  useEffect(() => {
    if (status === 'failed' && error) {
      alert(error);
      dispatch(clearError()); //เคลียร์ error ใน Redux หลังแสดงเสร็จ
    }
  }, [status, error, dispatch]);

  const displayText = useTypewriter(messages) // เรียกใช้ function พิมพ์ตัวอักษร
  return (
    <div className="grid grid-cols-2 min-h-screen mx-auto">
      {/* ฝั่งซ้าย Login Form*/ }
      <div className="flex items-center justify-center bg-gray-100 px-12">
        <div className="w-full max-w-md">
          <h1 className="text-5xl font-bold mb-8 text-center text-gray-800">Login</h1>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <input
              type="email"
              placeholder="example@gmail.com"
              className="input input-bordered w-full"
              {...register("email")}
            />
            {errors.email && (
              <p className="text-sm text-red-500">{errors.email.message}</p>
            )}

            <input
              type="password"
              placeholder="Password"
              className="input input-bordered w-full"
              {...register("password")}
            />
            {errors.password && (
              <p className="text-sm text-red-500">{errors.password.message}</p>
            )}

            <div className="grid grid-cols-2 gap-4 pt-2">
              <button className="btn btn-success w-full rounded-lg" type="submit" disabled={status === 'loading'}>
              {status === 'loading' ? 'Logging in…' : 'Login'}
              </button>
              <button
                className="btn btn-warning w-full rounded-lg"
                type="button"
                onClick={() => navigate("/register")}
              >
                Register
              </button>
            </div>
          </form>
        </div>
      </div>

      {/* ฝั่งขวา Image + ตัวอักษร */}
      <div className="relative overflow-hidden">
        <img
          src="./login_img.jpg"
          alt="Login"
          className="w-full h-full object-cover"
        />
        <div className="absolute bottom-0 left-0 w-full flex justify-center">
          <div className="bg-white/30 backdrop-blur-md w-full  p-3 shadow-inner">
            <h2 className="text-xl font-semibold text-black">{displayText}</h2>
          </div>
        </div>
      </div>
    </div>

  );
}

