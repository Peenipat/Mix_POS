// src/pages/Login.tsx
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
  const navigate = useNavigate()
  const onSubmit = async (data: LoginForm) => {
    try {
      const res = await axios.post("/auth/login", data)
      localStorage.setItem("token", res.data.token)
      navigate("/dashboard") 
    } catch (err) {
      alert("Login failed")
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input type="email" placeholder="Email" {...register("email")} />
      <p>{errors.email?.message}</p>

      <input type="password" placeholder="Password" {...register("password")} />
      <p>{errors.password?.message}</p>

      <button type="submit">Login</button>
      <button type="submit">Register</button>
    </form>
    
  );
}
