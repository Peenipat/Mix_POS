import { z } from "zod"

export const loginSchema = z.object({
    email: z.string().email({message:"Invaild email"}), 
    password: z.string().min(6, { message:"Passwrod too short"}),
})

export type LoginForm = z.infer<typeof loginSchema>;

export const passwordObjectSchema = z.object({
  password: z.string()
    .min(8, { message: "Password must be at least 8 characters long" })
    .regex(/[A-Z]/, { message: "Password must contain at least one uppercase letter" }) 
    .regex(/[a-z]/, { message: "Password must contain at least one lowercase letter" })
    .regex(/[0-9]/, { message: "Password must contain at least one number" })
    .regex(/[^A-Za-z0-9]/, { message: "Password must contain at least one special character" })
});

export const registerSchema = passwordObjectSchema.extend({
  username: z.string().min(3, { message: "Username is required" }),
  email: z.string().email({ message: "Invalid email" }),
  confirmPassword: z.string().min(6, { message: "Confirm Password is required" }),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords do not match",
  path: ["confirmPassword"],
});


  export type RegisterForm = z.infer<typeof registerSchema>;