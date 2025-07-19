const mockBarber = {
    id: 1,
    name: "คุณบอย",
    experience: "5 ปี",
    profilePicture:
        "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber3.jpg",
    services: ["ตัดผมชาย", "โกนหนวด", "สระผม"],
    availableTimes: ["10:00 - 12:00", "13:00 - 17:00"],
    averageRating: 4.7,
};

export default function BarberProfile() {
    return (
        <div className="max-w-xl mx-auto p-4 bg-white rounded-lg shadow space-y-4">
            <h2 className="text-2xl font-bold text-center">โปรไฟล์ช่างตัดผม</h2>

            <div className="flex flex-col items-center space-y-2">
                <img
                    src={mockBarber.profilePicture}
                    alt={mockBarber.name}
                    className="w-24 h-24 rounded-full object-cover"
                />
                <h3 className="text-xl font-semibold">{mockBarber.name}</h3>
                <p className="text-gray-500">ประสบการณ์: {mockBarber.experience}</p>
            </div>

            <div>
                <h4 className="font-semibold text-lg">บริการที่ให้</h4>
                <ul className="list-disc ml-6 text-gray-700">
                    {mockBarber.services.map((service, index) => (
                        <li key={index}>{service}</li>
                    ))}
                </ul>
            </div>

            <div>
                <h4 className="font-semibold text-lg">ช่วงเวลาให้บริการ</h4>
                <div className="text-gray-700">
                    {mockBarber.availableTimes.join(", ")}
                </div>
            </div>

            <div>
                <h4 className="font-semibold text-lg">เรตติ้งเฉลี่ย</h4>
                <div className="text-yellow-500 text-lg">
                    {Array.from({ length: 5 }).map((_, idx) => (
                        <span key={idx}>
                            {idx < Math.round(mockBarber.averageRating) ? "★" : "☆"}
                        </span>
                    ))}
                    <span className="text-sm text-gray-600 ml-2">
                        ({mockBarber.averageRating.toFixed(1)})
                    </span>
                </div>
            </div>
        </div>
    );
}
