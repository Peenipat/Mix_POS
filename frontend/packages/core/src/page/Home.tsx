import { useNavigate } from "react-router-dom"
import React from 'react'
import { Button } from "@object/shared";
export default function Home() {
  const navigate = useNavigate();
  return (
    <> 
    <h1>Welcome to Home page!</h1>
    <button
      className="btn btn-warning w-full rounded-lg"
      type="button"
      onClick={() => navigate("/register")}
    >
      Register
    </button>

    <button
      className="btn btn-warning w-full rounded-lg"
      type="button"
      onClick={() => navigate("/login")}
    >
      login
    </button>

    <Button>Button</Button>

    </>
  )
}
