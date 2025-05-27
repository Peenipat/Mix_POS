// components/UserList.tsx
import React from "react";
import type { User as UserResponse } from "../../../schemas/userSchema";
import { UserCard } from "./UserCard";

interface UserListProps {
  users: UserResponse[];
}

export function UserList({ users }: UserListProps) {
  if (users.length === 0) {
    return <p className="text-center text-gray-500">No users found.</p>;
  }

  return (
    <div className="grid gap-4">
      {users.map((u) => (
        <UserCard key={u.id} user={u} />
      ))}
    </div>
  );
}
