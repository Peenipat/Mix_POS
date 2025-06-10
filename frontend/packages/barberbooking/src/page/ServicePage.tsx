// src/page/ServicePage.tsx
import React from "react";

interface Service {
    id: number;
    name: string;
    description: string;
    price: number; // บาท
    duration: string; // ชั่วโมง:นาที
    imageUrl: string;
}

const mockServices: Service[] = [
    {
        id: 1,
        name: "ตัดผมชายทั่วไป",
        description: "ตัดผมทรงคลาสสิค ปรับทรงเรียบร้อย ดูสุภาพเหมาะกับทุกโอกาส",
        price: 250,
        duration: "30",
        imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/service/haircut_classic.jpg",
    },
    {
        id: 2,
        name: "ตกแต่งหนวด–เครา",
        description: "จัดแต่งหนวด เคราให้คมเข้ารูป พร้อมสครับผิวหน้าเบาๆ",
        price: 150,
        duration: "20",
        imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/service/Beard_Trim.jpg",
    },
    {
        id: 3,
        name: "สระ–นวดหนังศีรษะ",
        description: "บำรุงล้ำลึก ลดรังแค ผมเงางาม สุขภาพดี",
        price: 450,
        duration: "60",
        imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/service/Hair_spa.jpg",
    },
    {
        id: 4,
        name: "บริการดัด–ย้อมสีผม",
        description: "จัดแต่งหนวด เคราให้คมเข้ารูป พร้อมสครับผิวหน้าเบาๆ",
        price: 150,
        duration: "20",
        imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/service/Reversing_hair_color.jpg",
    },
    {
        id: 5,
        name: "ตัดผมชาย(เด็ก)",
        description: "จัดแต่งหนวด เคราให้คมเข้ารูป พร้อมสครับผิวหน้าเบาๆ",
        price: 150,
        duration: "20",
        imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/service/Kid_haircut.jpg",
    },
    {
        id: 5,
        name: "ตัดผมชาย(ทรงพิเศษ)",
        description: "จัดแต่งหนวด เคราให้คมเข้ารูป พร้อมสครับผิวหน้าเบาๆ",
        price: 150,
        duration: "20",
        imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/service/Special_haircut.jpg",
    },
];


export default function ServicePage() {
    return (
        <div className="min-h-screen bg-gray-900 text-gray-200">
            <div className="container mx-auto py-12 px-6">
                <h1 className="text-4xl font-extrabold mb-8">บริการของเรา</h1>

                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                    {mockServices.map((service) => (
                        // ใส่ flex flex-col h-full ที่ wrapper การ์ด
                        <div
                            key={service.id}
                            className="bg-gray-800 rounded-lg shadow-lg hover:shadow-xl transition flex flex-col h-full"
                        >
                            {/* รูปคงความสูงเท่านี้ */}
                            <img
                                src={service.imageUrl}
                                alt={service.name}
                                className="w-full h-48 object-cover rounded-t-lg"
                            />

                            {/* ส่วนเนื้อหา ดันราคาลงล่างด้วย flex-1 */}
                            <div className="p-4 flex flex-col justify-between flex-1">
                                <div>
                                    <h2 className="text-2xl font-semibold mb-2">
                                        {service.name}
                                    </h2>
                                    <p className="text-gray-400 text-sm mb-4">
                                        {service.description}
                                    </p>
                                </div>
                                <div className="flex items-center justify-between text-gray-100">
                                    <span className="font-bold text-lg">
                                        ฿{service.price}
                                    </span>
                                    <span className="text-sm bg-gray-700 px-2 py-1 rounded">
                                        {service.duration} นาที
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
