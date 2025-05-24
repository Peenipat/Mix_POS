// src/pages/Login.tsx
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerSchema, RegisterForm } from "../schemas/authSchema";
import { useNavigate } from "react-router-dom";
import axios from "../lib/axios";
import React from "react";
import { useTypewriter } from "../hooks/useTypewriter";

export default function Register() {
    const messages = [
        "Lorem ipsum dolor sit amet  accusamus non! Error voluptatibus dignissimos magnam ",
        "Lorem  accusantium et, solutais deleniti harum ex non. Magni, earum. Cupiditate?",
        "Lorem ipsum dolor sit amet consectetur adipisicing elit.  Error voluptCupiditate?",
      ];
    const displayText = useTypewriter(messages) // เรียกใช้ function พิมพ์ตัวอักษร
    const navigate = useNavigate();
    // ฟอร์ม register เชื่อมกับ schema ตรวจข้อมูลก่อนส่ง
    const { register, handleSubmit, formState: { errors } ,reset} = useForm<RegisterForm>({
        resolver: zodResolver(registerSchema),
    });
    // register function
    const onSubmit = async (data: RegisterForm) => {
        try {
            await axios.post("/auth/register", data)
            reset()
            navigate("/dashboard")
        } catch (error) {
            console.error(error);
            alert("Register failed");
        }
    }

    return (
        <div>
            <div className="grid grid-cols-2 min-h-screen mx-auto">
                {/* ฝั่งซ้าย: ฟอร์มสมัครสมาชิก */}
                <div className="flex items-center justify-center bg-gray-100 px-12">
                    <div className="w-full max-w-md">
                        <h1 className="text-5xl font-bold mb-8 text-center text-gray-800">Register</h1>
                        <form onSubmit={handleSubmit(onSubmit)}>
                            <input type="text" placeholder="Username" className="input input-bordered w-full my-2" {...register("username")} />
                            <input type="email" placeholder="Email" className="input input-bordered w-full my-2" {...register("email")} />
                            <input type="password" placeholder="Password" className="input input-bordered w-full my-2" {...register("password")} /> {errors.password && <p className="text-red-500 text-sm">{errors.password.message}</p>}
                            <input type="password" placeholder="Confirm password" className="input input-bordered w-full my-2" {...register("confirmPassword")} />{errors.confirmPassword && <p className="text-red-500 text-sm">{errors.confirmPassword.message}</p>}
                            {/* errors.confirmPassword?.message จะขึ้น ถ้าไม่ตรงกัน */}
                            <div className="pt-2 ">
                                <button
                                    className="btn btn-warning w-full rounded-lg"
                                    type="submit"
                                // onClick={}
                                >
                                    Register
                                </button>
                            </div>
                        </form>
                    </div>
                </div>

                {/* Right: Image + ตัวอักษร */}
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

        </div>

    );
}
