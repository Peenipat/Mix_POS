export const mockReviews = [
    {
        id: 1,
        customerName: "คุณบอย",
        date: "2025-07-01",
        rating: 5,
        comment: "ช่างตัดดีมาก เป็นกันเองครับ 😊",
        appointmentTime: "10:00 - 10:30",
        serviceName: "ตัดผมชาย",
    },
    {
        id: 2,
        customerName: "คุณเมย์",
        date: "2025-07-02",
        rating: 4,
        comment: "ช่างมือเบา พูดจาดีค่ะ แต่รอนิดนึง",
        appointmentTime: "13:00 - 14:00",
        serviceName: "ย้อมสีผม",
    },
    {
        id: 3,
        customerName: "คุณตูน",
        date: "2025-07-03",
        rating: 3,
        comment: "ผมโอเค แต่ร้านร้อนนิดหน่อย",
        appointmentTime: "11:00 - 11:45",
        serviceName: "ตัด+เซ็ตผม",
    },
    {
        id: 4,
        customerName: "คุณเอิร์ธ",
        date: "2025-07-03",
        rating: 5,
        comment: "ผมทรงนี้หล่อเลยครับ ชอบมาก!",
        appointmentTime: "16:00 - 16:30",
        serviceName: "ตัดผมชาย",
    },
    {
        id: 5,
        customerName: "คุณหญิง",
        date: "2025-07-04",
        rating: 4,
        comment: "บริการดี ร้านสะอาดค่ะ",
        appointmentTime: "14:00 - 14:45",
        serviceName: "ทำทรีทเมนต์ผม",
    },
    {
        id: 6,
        customerName: "คุณมาร์ค",
        date: "2025-07-05",
        rating: 2,
        comment: "ตัดช้ากว่าที่คิด แต่ก็โอเค",
        appointmentTime: "09:00 - 09:30",
        serviceName: "ตัดผมชาย",
    },
    {
        id: 7,
        customerName: "คุณจูน",
        date: "2025-07-05",
        rating: 5,
        comment: "เซ็ตผมสวยมาก ถ่ายรูปขึ้นเลยค่ะ ❤️",
        appointmentTime: "10:30 - 11:00",
        serviceName: "เซ็ตผม",
    },
    {
        id: 8,
        customerName: "คุณแบงค์",
        date: "2025-07-06",
        rating: 4,
        comment: "ดีครับ ช่างมีประสบการณ์",
        appointmentTime: "12:00 - 12:30",
        serviceName: "โกนหนวด",
    },
    {
        id: 9,
        customerName: "คุณแนต",
        date: "2025-07-06",
        rating: 3,
        comment: "พอใช้ได้ แต่ยังไม่ถูกใจทรงผม",
        appointmentTime: "13:30 - 14:00",
        serviceName: "ตัดผมชาย",
    },
    {
        id: 10,
        customerName: "คุณแนน",
        date: "2025-07-07",
        rating: 5,
        comment: "รักเลย! สีผมสวยมากกก 🧡",
        appointmentTime: "15:00 - 16:00",
        serviceName: "ย้อมสีผม",
    },
    {
        id: 11,
        customerName: "คุณดิว",
        date: "2025-07-07",
        rating: 4,
        comment: "ร้านอยู่ใกล้ เดินทางสะดวก",
        appointmentTime: "17:00 - 17:30",
        serviceName: "ตัด+เซ็ตผม",
    },
    {
        id: 12,
        customerName: "คุณหมวย",
        date: "2025-07-08",
        rating: 5,
        comment: "บริการประทับใจ ให้คำแนะนำดีมากค่ะ",
        appointmentTime: "10:00 - 10:45",
        serviceName: "ตัดผมสั้นผู้หญิง",
    },
    {
        id: 13,
        customerName: "คุณวิน",
        date: "2025-07-08",
        rating: 2,
        comment: "ผมไม่เท่ากันนิดหน่อยครับ",
        appointmentTime: "11:00 - 11:30",
        serviceName: "ตัดผมชาย",
    },
    {
        id: 14,
        customerName: "คุณอิ้งค์",
        date: "2025-07-09",
        rating: 5,
        comment: "ดีงามค่ะ ไปอีกแน่นอน 😍",
        appointmentTime: "13:00 - 14:00",
        serviceName: "สระไดร์",
    },
    {
        id: 15,
        customerName: "คุณพีท",
        date: "2025-07-09",
        rating: 3,
        comment: "กลางๆ ครับ ไม่แย่แต่ยังไม่ว้าว",
        appointmentTime: "16:00 - 16:30",
        serviceName: "ตัดผมชาย",
    },
];
import { useState } from "react";

export default function CustomerReview() {
    const [filterRating, setFilterRating] = useState<number | null>(null);
    const [currentPage, setCurrentPage] = useState(1);

    const filteredReviews =
        filterRating === null
            ? mockReviews
            : mockReviews.filter((r) => r.rating === filterRating);

    const itemsPerPage = 6;
    const totalPages = Math.ceil(filteredReviews.length / itemsPerPage);
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const paginatedReviews = filteredReviews.slice(startIndex, endIndex);

    return (
        <div className="space-y-4">
            <h2 className="text-xl font-bold mb-4">รีวิวจากลูกค้า</h2>

            {/* ฟิลเตอร์ตามดาว */}
            <div className="flex gap-2 flex-wrap">
                <button
                    onClick={() => {
                        setFilterRating(null);
                        setCurrentPage(1);
                    }}
                    className={`px-3 py-1 rounded ${filterRating === null
                        ? "bg-blue-600 text-white"
                        : "bg-gray-200"
                        }`}
                >
                    ทั้งหมด
                </button>
                {[5, 4, 3, 2, 1].map((rating) => (
                    <button
                        key={rating}
                        onClick={() => {
                            setFilterRating(rating);
                            setCurrentPage(1);
                        }}
                        className={`px-3 py-1 rounded ${filterRating === rating
                            ? "bg-blue-600 text-white"
                            : "bg-gray-200"
                            }`}
                    >
                        {rating} ดาว
                    </button>
                ))}
            </div>

            {/* รีวิว */}
            {paginatedReviews.length > 0 ? (
                <div className="grid grid-cols-2 gap-2">
                    {paginatedReviews.map((review) => (
                        <div
                            key={review.id}
                            className="p-4 border rounded-lg bg-white shadow-sm space-y-1"
                        >
                            <div className="flex justify-between items-center">
                                <span className="text-xl font-semibold">{review.customerName}</span>
                                <span className="text-md text-gray-500">{review.date}</span>
                            </div>
                            <div className="text-sm text-gray-600">
                                เวลา: {review.appointmentTime} | บริการ: {review.serviceName}
                            </div>
                            <div className="text-yellow-500 text-sm">
                                {Array.from({ length: 5 }).map((_, idx) => (
                                    <span key={idx}>{idx < review.rating ? "★" : "☆"}</span>
                                ))}
                            </div>
                            <div className="text-gray-700">{review.comment}</div>
                        </div>
                    ))}
                </div>
            ) : (
                <p className="text-gray-500">ไม่มีรีวิวในช่วงนี้</p>
            )}

            {/* Pagination */}
            {filteredReviews.length > 0 && (
                <div className="flex flex-col items-center mt-4 space-y-1">
                    <p className="text-sm text-gray-600">
                        รีวิวทั้งหมด {filteredReviews.length} รายการ | กำลังแสดง{" "}
                        {filteredReviews.length === 0
                            ? 0
                            : `${startIndex + 1}–${Math.min(endIndex, filteredReviews.length)}`}
                    </p>

                    {totalPages > 1 && (
                        <div className="flex gap-2">
                            {Array.from({ length: totalPages }).map((_, idx) => {
                                const page = idx + 1;
                                return (
                                    <button
                                        key={page}
                                        onClick={() => setCurrentPage(page)}
                                        className={`px-3 py-1 rounded ${currentPage === page
                                            ? "bg-blue-600 text-white"
                                            : "bg-gray-200"
                                            }`}
                                    >
                                        {page}
                                    </button>
                                );
                            })}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
}
 