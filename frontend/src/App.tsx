import { Routes, Route } from "react-router-dom";
import Login from "./page/Login";
import Dashboard from "./page/Dashboard";
import ProtectedRoute from "./components/ProtectedRoute"

function App() {
  return (
    <Routes>
      <Route path="/" element={<Login />} />
      <Route path="/dashboard" element={
        <ProtectedRoute>
          <Dashboard />
        </ProtectedRoute>} />
    </Routes>
  );
}

export default App;
