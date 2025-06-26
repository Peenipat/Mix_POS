// src/page/ServicePage.tsx
import { useState, useCallback,useEffect } from "react";
import axios from "../lib/axios";
import type { Barber } from "../types/barber";

const dummySlots: TimeSlot[] = [
    { id: 1, date: "2025-06-12", time: "10:00 – 11:00", available: true },
    { id: 2, date: "2025-06-12", time: "11:00 – 12:00", available: false },
    { id: 3, date: "2025-06-13", time: "13:00 – 14:00", available: true },
];

export default function BarberPage() {
const [barbers, setBarbers] = useState<Barber[]>([]);
const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
const [errorBarbers, setErrorBarbers] = useState<string | null>(null);

const loadBarbers = useCallback(async () => {
    setLoadingBarbers(true);
    setErrorBarbers(null);
    try {
        const res = await axios.get<{ status: string; data: Barber[] }>(
            `/barberbooking/tenants/1/barbers/branches/1/barbers`
        );
        if (res.data.status !== "success") {
            throw new Error(res.data.status);
        }
        setBarbers(res.data.data);
    } catch (err: any) {
        setErrorBarbers(err.response?.data?.message || err.message || "Failed to load barbers");
    } finally {
        setLoadingBarbers(false);
    }
}, []);

useEffect(() => {
      loadBarbers();
  }, []);
    const [isModalOpen, setModalOpen] = useState(false);
    return (
        <div className="min-h-screen bg-gradient-to-b from-white via-slate-100 to-slate-200 text-gray-900">
            <div className="container mx-auto py-6 px-6">
                <h1 className="text-4xl font-extrabold mb-8">ช่างของเรา</h1>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                    {loadingBarbers && <p>Loading barbers…</p>}
                    {errorBarbers && <p className="text-red-500">Error loading barbers: {errorBarbers}</p>}
                    {barbers.map((barber) => (
                        <div
                            key={barber.id}
                            className="bg-gray-200 rounded-lg shadow-lg hover:shadow-xl transition flex flex-col h-full"
                        >
                            <img
                                src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${barber.img_path}/${barber.img_name}`}
                                alt={barber.username}
                                className="w-full h-64 object-cover object-top rounded-t-lg"
                            />

                            <div className="p-4 flex flex-col justify-between flex-1">
                                <div>
                                    <h2 className="text-2xl font-semibold mb-2">
                                        {barber.username}
                                    </h2>
                                    <p className="text-gray-400 text-sm mb-2">
                                        {barber.description}
                                    </p>
                                </div>
                                <div className="flex items-center justify-between text-gray-900 mt-4">
                                    <span className="text-sm bg-gray-300 p-1.5 rounded">
                                        ⭐ 4.5
                                    </span>

                                    <button className="px-3 py-1.5 bg-gray-400 rounded text-white hover:bg-blue-700 transition"
                                        onClick={() => setModalOpen(true)}>
                                        ดูคิวของช่าง
                                    </button>

                                    <BarberScheduleModal
                                        isOpen={isModalOpen}
                                        onClose={() => setModalOpen(false)}
                                        slots={dummySlots}
                                        barberName="สมชาย ศรีสุข"
                                    />
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}

interface TimeSlot {
    id: number;
    date: string;       // ex: "2025-06-12"
    time: string;       // ex: "10:00 – 11:00"
    available: boolean; // true = ว่าง, false = ไม่ว่าง
}

interface BarberScheduleModalProps {
    isOpen: boolean;
    onClose: () => void;
    slots: TimeSlot[];
    barberName: string;
}

function BarberScheduleModal({
    isOpen,
    onClose,
    slots,
    barberName,
}: BarberScheduleModalProps) {
    if (!isOpen) return null;

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
            onClick={onClose}
        >
            <div
                className="bg-white rounded-lg w-full max-w-2xl mx-4 shadow-lg overflow-hidden"
                onClick={(e) => e.stopPropagation()}
            >
                {/* Header */}
                <div className="flex justify-between items-center px-6 py-4 border-b text-black">
                    <h3 className="text-xl font-semibold">
                        ตารางงานของ {barberName}
                    </h3>
                    <button
                        className="text-gray-500 hover:text-gray-700"
                        onClick={onClose}
                    >
                        ✕
                    </button>
                </div>

                {/* Body */}
                <div className="p-6 max-h-[60vh] overflow-auto text-black">
                    {slots.length === 0 ? (
                        <p className="text-center text-gray-500">ยังไม่มีข้อมูลตารางงาน</p>
                    ) : (
                        <table className="w-full text-left">
                            <thead>
                                <tr className="bg-gray-100">
                                    <th className="px-4 py-2">วันที่</th>
                                    <th className="px-4 py-2">ช่วงเวลา</th>
                                    <th className="px-4 py-2">สถานะ</th>
                                </tr>
                            </thead>
                            <tbody>
                                {slots.map((slot) => (
                                    <tr key={slot.id} className="border-b last:border-none">
                                        <td className="px-4 py-3">{slot.date}</td>
                                        <td className="px-4 py-3">{slot.time}</td>
                                        <td className="px-4 py-3">
                                            {slot.available ? (
                                                <span className="inline-block px-3 py-1 text-sm bg-green-100 text-green-800 rounded">
                                                    ว่าง
                                                </span>
                                            ) : (
                                                <span className="inline-block px-3 py-1 text-sm bg-red-100 text-red-800 rounded">
                                                    ไม่ว่าง
                                                </span>
                                            )}
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>

                {/* Footer */}
                <div className="flex justify-end px-6 py-4 border-t">
                    <button
                        className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 transition"
                        onClick={onClose}
                    >
                        ปิด
                    </button>
                </div>
            </div>
        </div>
    );
}