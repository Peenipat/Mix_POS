// // src/pages/Login.tsx
// import { useForm } from "react-hook-form";
// import { zodResolver } from "@hookform/resolvers/zod";
// import { loginSchema } from "../schemas/authSchema";
// import { z } from "zod";
// import axios from "../lib/axios";

// type LoginForm = z.infer<typeof loginSchema>;

// export default function Login() {
//   const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>({
//     resolver: zodResolver(loginSchema),
//   });

//   const onSubmit = async (data: LoginForm) => {
//     try {
//         console.log(data)
//       const res = await axios.post("/auth/login", data);
//       alert("Login success: " + res.data.token);
//     } catch (err: any) {
//       alert("Login failed: " + err.response?.data?.error || "Unknown error");
//     }
//   };

//   return (
//     <form onSubmit={handleSubmit(onSubmit)}>
//       <input type="email" placeholder="Email" {...register("email")} />
//       <p>{errors.email?.message}</p>

//       <input type="password" placeholder="Password" {...register("password")} />
//       <p>{errors.password?.message}</p>

//       <button type="submit">Login</button>
//     </form>
//   );
// }
