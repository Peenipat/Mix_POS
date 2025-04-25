import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginSchema, LoginForm } from "../schemas/authSchema";
import { useNavigate } from "react-router-dom"
import axios from "../lib/axios";
import { useEffect, useState } from "react";

export default function Login() {
  const [displayText, setDisplayText] = useState("");
  const [messageIndex, setMessageIndex] = useState(0);
  const [charIndex, setCharIndex] = useState(0);
  const [deleting, setDeleting] = useState(false);
  const messages = [
    "Lorem ipsum dolor sit amet  accusamus non! Error voluptatibus dignissimos magnam ",
    "Lorem  accusantium et, solutais deleniti harum ex non. Magni, earum. Cupiditate?",
    "Lorem ipsum dolor sit amet consectetur adipisicing elit.  Error voluptCupiditate?",
  ];
  const { register, handleSubmit, formState: { errors } } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
  });
  const navigate = useNavigate();

  const onSubmit = async (data: LoginForm) => {
    try {
      const res = await axios.post("/auth/login", data);
      const { token, user } = res.data;
      localStorage.setItem("token", token);
      localStorage.setItem("role", user.role);

      if (user.role === "SUPER_ADMIN") navigate("/admin/dashboard");
      else if (user.role === "BRANCH_ADMIN") navigate("/branch/dashboard");
      else if (user.role === "STAFF") navigate("/staff/dashboard");
      else navigate("/dashboard");
    } catch (err) {
      alert("Login failed");
    }
  };

  useEffect(() => {
    const currentMessage = messages[messageIndex];
    const delay = deleting ? 40 : 70;

    const timeout = setTimeout(() => {
      if (!deleting) {
        setDisplayText(currentMessage.slice(0, charIndex + 1));
        setCharIndex(charIndex + 1);

        if (charIndex + 1 === currentMessage.length) {
          setTimeout(() => setDeleting(true), 1500); 
        }
      } else {
        setDisplayText(currentMessage.slice(0, charIndex - 1));
        setCharIndex(charIndex - 1);

        if (charIndex === 0) {
          setDeleting(false);
          setMessageIndex((messageIndex + 1) % messages.length); // วนลูปข้อความ
        }
      }
    }, delay);

    return () => clearTimeout(timeout);
  }, [charIndex, deleting]);
  return (
    <div className="grid grid-cols-2 min-h-screen mx-auto">
  {/* Left: Login Form */}
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
          <button className="btn btn-success w-full rounded-lg" type="submit">
            Login
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

  {/* Right: Image + Mix POS */}
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

