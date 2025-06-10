import React, { useState } from "react";

export default function TimeSelector() {
  const [hour, setHour] = useState("00");
  const [minute, setMinute] = useState("00");

  const hours = Array.from({ length: 24 }, (_, i) =>
    String(i).padStart(2, "0")
  );
  const minutes = ["00", "10", "15", "30", "45", "55"];

  return (
    <form className="max-w-sm mx-auto mt-4">
      <label
        htmlFor="hour"
        className="block mb-2 text-sm font-medium text-gray-900 dark:text-white"
      >
        Select time:
      </label>
      <div className="flex space-x-2">
        {/* Select ชั่วโมง */}
        <div className="relative w-1/2">
          <select
            id="hour"
            value={hour}
            onChange={(e) => setHour(e.target.value)}
            className="appearance-none bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg
                       focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 pr-8
                       dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white
                       dark:focus:ring-blue-500 dark:focus:border-blue-500"
          >
            {hours.map((h) => (
              <option key={h} value={h}>
                {h}
              </option>
            ))}
          </select>
          <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
            <svg
              className="w-4 h-4 text-gray-500 dark:text-gray-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 20 20"
            >
              <path
                stroke="currentColor"
                strokeWidth="2"
                d="M6 8l4 4 4-4"
              />
            </svg>
          </div>
        </div>

        {/* Select นาที */}
        <div className="relative w-1/2">
          <select
            id="minute"
            value={minute}
            onChange={(e) => setMinute(e.target.value)}
            className="appearance-none bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg
                       focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 pr-8
                       dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white
                       dark:focus:ring-blue-500 dark:focus:border-blue-500"
          >
            {minutes.map((m) => (
              <option key={m} value={m}>
                {m}
              </option>
            ))}
          </select>
          <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3">
            <svg
              className="w-4 h-4 text-gray-500 dark:text-gray-400"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 20 20"
            >
              <path
                stroke="currentColor"
                strokeWidth="2"
                d="M6 8l4 4 4-4"
              />
            </svg>
          </div>
        </div>
      </div>

      {/* แสดงค่าที่เลือก */}
      <div className="mt-2 text-gray-700 dark:text-gray-300">
        Selected: {hour}:{minute}
      </div>
    </form>
  );
}
