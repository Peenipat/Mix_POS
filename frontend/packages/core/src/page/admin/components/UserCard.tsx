// components/UserCard.tsx
import React from "react";
import type { User } from "../../../schemas/userSchema";

interface UserCardProps {
  user: User;
}

export function UserCard({ user }: UserCardProps) {
  return (
    <div className="flex items-center space-x-4 p-4 bg-white rounded-lg shadow">
      {user.img_path && (
        <img
          src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${user.img_path}/${user.img_name}`}
          alt={`${user.username} avatar`}
          className="w-12 h-12 rounded-full object-cover"
        />
      )}
      <div>
        <div className="font-medium text-gray-900">{user.username}</div>
        <div className="text-sm text-gray-500">{user.email}</div>
        <div className="text-xs text-gray-400">{user.role}</div>
      </div>
    </div>
  );
}
