import React, { useState, useEffect } from "react";
import type { UseFormSetValue } from "react-hook-form";
import type { appointmentForm } from "../schemas/appointmentSchema"; // แก้ path ให้ถูกต้อง

type Props = {
  setValue: UseFormSetValue<appointmentForm>;
};

export default function TimeSelector({ setValue }: Props) {
  const [hour, setHour] = useState("00");
  const [minute, setMinute] = useState("00");

  const hours = Array.from({ length: 24 }, (_, i) =>
    String(i).padStart(2, "0")
  );
  const minutes = ["00", "10", "15", "30", "45", "55"];

  useEffect(() => {
    const time = `${hour}:${minute}`;
    setValue("time", time); 
  }, [hour, minute, setValue]);

  return (
    <div className="max-w-md mx-auto">
      <div className="flex space-x-2">
        <div className="relative w-1/2">
          <select
            id="hour"
            value={hour}
            onChange={(e) => setHour(e.target.value)}
            className="select select-bordered w-full"
          >
            {hours.map((h) => (
              <option key={h} value={h}>
                {h}
              </option>
            ))}
          </select>
        </div>
        <div className="relative w-1/2">
          <select
            id="minute"
            value={minute}
            onChange={(e) => setMinute(e.target.value)}
            className="select select-bordered w-full"
          >
            {minutes.map((m) => (
              <option key={m} value={m}>
                {m}
              </option>
            ))}
          </select>
        </div>
      </div>
    </div>
  );
}
