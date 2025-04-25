import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema } from "../schemas/authSchema";
import { useNavigate } from "react-router-dom"
import { z } from "zod";
import axios from "../lib/axios";

type LoginForm = z.infer<typeof loginSchema>;

export default function Login() {
  const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });
  const navigate = useNavigate();

  const onSubmit = async (data: LoginForm) => {
    try {
      const res = await axios.post("/auth/login", data);
      const { token, user } = res.data;

      // เก็บ token และ role
      localStorage.setItem("token", token);
      localStorage.setItem("role", user.role);

      // Redirect ตาม role
      if (user.role === "SUPER_ADMIN") {
        navigate("/admin/dashboard"); // หน้าของ หน้า super admin
      } else if (user.role === "BRANCH_ADMIN") {
        navigate("/branch/dashboard"); // หน้าของ admin แต่ละร้าน
      } else if (user.role === "STAFF") {
        navigate("/staff/dashboard"); // หน้าของ staff
      }else {
        navigate("/dashboard"); // default หรือหน้า user ทั่วไป
      }
    } catch (err) {
      alert("Login failed");
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input type="email" placeholder="Email" {...register("email")} />
      <p>{errors.email?.message}</p>

      <input type="password" placeholder="Password" {...register("password")} />
      <p>{errors.password?.message}</p>

      <button type="submit">Login</button>
      <button type="button" onClick={() => navigate("/register")}>Register</button>
    </form>
  );
}
