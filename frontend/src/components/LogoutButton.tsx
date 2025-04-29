import { useNavigate } from "react-router-dom";
import { useAppDispatch } from "../store/hook" 
import { logout } from "@/store/authSlice";

interface LogoutButtonProps {
  className?: string; 
}

export default function LogoutButton({ className }: LogoutButtonProps) {
  const navigate = useNavigate();
  const dispatch = useAppDispatch()

  const handleLogout = () => {
    dispatch(logout());
    localStorage.removeItem("token");
    localStorage.removeItem("role");
    navigate("/");
  };

  return (
    <button onClick={handleLogout} className={`btn btn-error ${className}`}>
      Logout
    </button>
  );
}
