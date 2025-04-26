import { useNavigate } from "react-router-dom";

interface LogoutButtonProps {
  className?: string; 
}

export default function LogoutButton({ className }: LogoutButtonProps) {
  const navigate = useNavigate();

  const handleLogout = () => {
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
