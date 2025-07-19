import React, { useState, useEffect } from "react";
import type { UseFormSetValue } from "react-hook-form";
import type { appointmentForm } from "../schemas/appointmentSchema";
import dayjs from "dayjs";
import axios from "../lib/axios";


export type WorkingHourResult = { start: string; end: string } | null;

export const getWorkingHourRangeAxios = async (
  tenantId: number,
  branchId: number,
  date: Date
): Promise<WorkingHourResult> => {
  const yyyyMMdd = dayjs(date).format("YYYY-MM-DD");

  try {
    const overrideRes = await axios.get<{
      id: number;
      branch_id: number;
      work_date: string;
      start_time: string;
      end_time: string;
      IsClosed: boolean;
    }[]>(`/barberbooking/tenants/${tenantId}/branches/${branchId}/working-day-overrides/date?start=${yyyyMMdd}&end=${yyyyMMdd}`);

    const overrides = overrideRes.data;

    if (overrides.length > 0) {
      const override = overrides[0];
      if (!override.IsClosed) {
        return {
          start: override.start_time,
          end: override.end_time,
        };
      } else {
        return null; 
      }
    }
  } catch (err: any) {
    console.warn("⚠️ override API error:", err.message);
  }

  try {
    const workingHourRes = await axios.get<{
      status: string;
      data: {
        Weekday: number;
        start_time: string;
        end_time: string;
        IsClosed: boolean;
      }[];
    }>(`/barberbooking/tenants/${tenantId}/workinghour/branches/${branchId}`);

    const weekday = dayjs(date).day();
    const today = workingHourRes.data.data.find((item) => item.Weekday === weekday);

    if (!today || today.IsClosed) {
      return null;
    }

    return {
      start: dayjs(today.start_time).format("HH:mm"),
      end: dayjs(today.end_time).format("HH:mm"),
    };
  } catch (err: any) {
    console.error("❌ Failed to load default working hours:", err.message);
    return null;
  }
};

type Props = {
  setValue: UseFormSetValue<appointmentForm>;
  date: string;
  disabled?: boolean;
};

export default function TimeSelector({ setValue, date,disabled }: Props) {
  const [hour, setHour] = useState("00");
  const [minute, setMinute] = useState("00");

  const [availableHours, setAvailableHours] = useState<string[]>([]);
  const [availableMinutes, setAvailableMinutes] = useState<string[]>([]);
  const [isClosed, setIsClosed] = useState(false);

  useEffect(() => {
    if (!date) return;

    async function loadWorkingHours() {
      const dateObj = new Date(date); 
      const result = await getWorkingHourRangeAxios(1, 1, dateObj);
      console.log("workingHour", result);

      if (!result) {
        setIsClosed(true);
        setAvailableHours([]);
        setAvailableMinutes([]);
        return;
      }

      setIsClosed(false);

      const [startH, startM] = result.start.split(":").map(Number);
      const [endH, endM] = result.end.split(":").map(Number);

      const hours: string[] = [];
      for (let h = startH; h <= endH; h++) {
        hours.push(String(h).padStart(2, "0"));
      }

      const mins = ["00", "15", "30", "45"];
      setAvailableHours(hours);
      setAvailableMinutes(mins);
      setHour(String(startH).padStart(2, "0"));
      setMinute("00");
    }

    loadWorkingHours();
  }, [date]);

  useEffect(() => {
    const time = `${hour}:${minute}`;
    setValue("time", time);
  }, [hour, minute, setValue]);

  return (
    <div className="max-w-md mx-auto">
      {!date || disabled ? (
        <div className="text-gray-500 font-medium text-center italic">
          กรุณาเลือกวันที่ก่อน
        </div>
      ) : isClosed ? (
        <div className="text-red-500 font-medium text-center">
          วันนี้ร้านปิด
        </div>
      ) : (
        <div className="flex space-x-2">
          <div className="relative w-1/2">
            <select
              id="hour"
              value={hour}
              onChange={(e) => setHour(e.target.value)}
              className="select select-bordered w-full"
              disabled={disabled}
            >
              {availableHours.map((h) => (
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
              disabled={disabled}
            >
              {availableMinutes.map((m) => (
                <option key={m} value={m}>
                  {m}
                </option>
              ))}
            </select>
          </div>
        </div>
      )}
    </div>
  );
}