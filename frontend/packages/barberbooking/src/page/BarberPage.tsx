// src/page/ServicePage.tsx
import React from "react";

interface Barber {
    id: number;
    name: string;
    specialization: string;
    experienceYears: number; // ปี
    rating: number; // คะแนน 0–5
    imageUrl: string;
}

const mockBarbers: Barber[] = [
    {
        id: 1,
        name: "สมชาย ศรีสุข",
        specialization: "ตัดผมชายทั่วไป",
        experienceYears: 5,
        rating: 4.9,
        imageUrl:
            "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/barbers/barber1.jpg",
    },
    {
        id: 2,
        name: "วิทยา ตัดทรง",
        specialization: "ตกแต่งหนวด–เครา",
        experienceYears: 3,
        rating: 4.7,
        imageUrl:
            "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/barbers/barber2.jpg",
    },
    {
        id: 3,
        name: "อรทัย นวดผม",
        specialization: "สระ–นวดหนังศีรษะ",
        experienceYears: 4,
        rating: 4.8,
        imageUrl:
            "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/barbers/barber3.jpg",
    },
    {
        id: 4,
        name: "ฐิติพงษ์ ฝีมือดี",
        specialization: "ดัด–ย้อมสีผม",
        experienceYears: 6,
        rating: 4.5,
        imageUrl:
            "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/barbers/barber4.jpg",
    },
];

export default function BarberPage() {
    return (
        <div className="min-h-screen bg-gray-900 text-gray-200">
            <div className="container mx-auto py-12 px-6">
                <h1 className="text-4xl font-extrabold mb-8">ช่างของเรา</h1>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                    {mockBarbers.map((barber) => (
                        <div
                            key={barber.id}
                            className="bg-gray-800 rounded-lg shadow-lg hover:shadow-xl transition flex flex-col h-full"
                        >
                            <img
                                src={barber.imageUrl}
                                alt={barber.name}
                                className="w-full h-64 object-cover object-top rounded-t-lg"
                            />

                            <div className="p-4 flex flex-col justify-between flex-1">
                                <div>
                                    <h2 className="text-2xl font-semibold mb-2">
                                        {barber.name}
                                    </h2>
                                    <p className="text-gray-400 text-sm mb-2">
                                        {barber.specialization}
                                    </p>
                                    <p className="text-gray-400 text-sm">
                                        ประสบการณ์ {barber.experienceYears} ปี
                                    </p>
                                </div>
                                <div className="flex items-center justify-between text-gray-100 mt-4">
                                    <span className="text-sm bg-gray-700 px-2 py-1 rounded">
                                        ⭐ {barber.rating.toFixed(1)}
                                    </span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
}
