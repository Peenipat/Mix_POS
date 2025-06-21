export interface IDCardProps {
    logoSrc?: string;        // URL โลโก้ (ไม่จำเป็นต้องส่งก็ได้)
    avatarSrc?: string;      // URL รูปโพรไฟล์
    avatarAlt?: string;      // alt ของรูป
    name?: string;           // ชื่อ-นามสกุล
    jobTitle?: string;       // ตำแหน่ง
    // idNo?: string;           // รหัสพนักงาน
    // dob?: string;            // วันเกิด (Date of Birth)
    email?: string;          // อีเมล
    phone?: string;          // เบอร์โทรศัพท์
    sex?: string;            // เพศ
    // Callback เมื่อคลิกปุ่ม
    onView?: () => void;
    onEdit?: () => void;
    onDelete?: () => void;
    // ถ้าต้องการส่ง className เพิ่มเติมเข้ามาปรับ styling
    className?: string;
}

export const CardItem: React.FC<IDCardProps> = ({
    logoSrc = "",
    avatarSrc = "",
    avatarAlt = "Avatar",
    name = "Your Name",
    jobTitle = "Job Title",
    // idNo = "000-0000",
    // dob = "01-01-1990",
    email = "youremail@example.com",
    phone = "+1-2345-6789",
    sex = "Male",
    onView = () => { },
    onEdit = () => { },
    onDelete = () => { },
    className = "",
}) => {
    return (
        <div
            className={`max-w-xs mx-auto bg-white rounded-xl shadow-lg overflow-hidden ${className}`}
        >
            {logoSrc && (
                <img
                    src={logoSrc}
                    alt="Logo"
                    className="absolute top-2 left-2 h-12 w-auto object-contain"
                />
            )}
            <div className="h-48 bg-green-100 overflow-hidden">
                {avatarSrc ? (
                    <img
                        src={avatarSrc}
                        alt={avatarAlt}
                        className="w-full h-full object-cover object-top"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center bg-gray-200">
                        <span className="text-gray-500 text-sm">No Photo</span>
                    </div>
                )}
            </div>

            {/* เนื้อหาด้านล่าง (พื้นขาว) */}
            <div className="pt-12 pb-3 px-3 text-center">
                {/* ชื่อ */}
                <h2 className="text-xl font-bold text-gray-800">{name}</h2>
                {/* ตำแหน่ง */}
                {/* <p className="text-sm text-gray-600">{jobTitle}</p> */}

                {/* ข้อมูลต่าง ๆ แบบ Grid สองคอลัมน์ */}
                <div className="mt-4 text-left">
                    <div className="grid grid-cols-3 gap-2 text-sm text-gray-700">
                        <span className="col-span-3">
                            <span className="font-semibold">Email: </span>{email}
                        </span>
                        <span className="col-span-3">
                            <span className="font-semibold">Phone:</span> {phone}
                        </span>
                        <span className="col-span-3">
                            <span className="font-semibold">Sex:</span> {sex}
                        </span>
                    </div>
                </div>
            </div>

            {/* ปุ่ม View / Edit / Delete */}
            <div className="px-6 pb-6">
                <div className="flex space-x-2">
                    <button
                        onClick={onView}
                        className="flex-1 inline-flex justify-center items-center px-3 py-2 text-sm font-medium text-white bg-green-600 rounded-lg hover:bg-green-700 focus:ring-4 focus:outline-none focus:ring-green-300"
                    >
                        View
                    </button>
                    <button
                        onClick={onEdit}
                        className="flex-1 inline-flex justify-center items-center px-3 py-2 text-sm font-medium text-white bg-yellow-500 rounded-lg hover:bg-yellow-600 focus:ring-4 focus:outline-none focus:ring-yellow-300"
                    >
                        Edit
                    </button>
                    <button
                        onClick={onDelete}
                        className="flex-1 inline-flex justify-center items-center px-3 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 focus:ring-4 focus:outline-none focus:ring-red-300"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    );
};
