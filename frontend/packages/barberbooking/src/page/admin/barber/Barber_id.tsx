// BarberDetailAdmin.tsx
import { useParams } from "react-router-dom";

const BarberDetail = () => {
  const { id } = useParams<{ id: string }>();
  const mockBarberWithAdminData = {
    id: "1",
    name: "Somchai Haircut",
    role: "ช่างตัดผมทั่วไป",
    rating: 4.9,
    imageUrl: "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/barbers/barber4.jpg",
    identityDocument: {
      type: "บัตรประชาชน",
      number: "1234567890123",
      fileUrl: "/uploads/id-somchai.pdf",
    },
    salary: 15000,
    payslips: [
      { id: "slip1", month: "พฤษภาคม 2025", fileUrl: "/uploads/slip-may.pdf" },
      { id: "slip2", month: "มิถุนายน 2025", fileUrl: "/uploads/slip-june.pdf" },
    ],
  };
  const barber = mockBarberWithAdminData;

  return (
    <div className="p-6 max-w-4xl mx-auto bg-white shadow rounded">
      <div className="flex items-center gap-4">
        <img src={barber.imageUrl} alt={barber.name} className="w-24 h-24 rounded-full  object-cover object-top" />
        <div>
          <h1 className="text-2xl font-bold">{barber.name}</h1>
          <p>ตำแหน่ง: {barber.role}</p>
          <p>คะแนนรีวิว: ⭐ {barber.rating}</p>
        </div>
      </div>

      <hr className="my-4" />

      <div className="flex justify-between gap-6">
        {/* ข้อมูลประจำตัว */}
        <div className="flex flex-col justify-between h-full min-h-[200px] w-1/3 border rounded-lg p-3">
          <div >
            <h2 className="text-xl font-semibold mb-2 text-center">ข้อมูลประจำตัว</h2>
            <p>ประเภทเอกสาร: {barber.identityDocument.type}</p>
            <p>หมายเลข: {barber.identityDocument.number}</p>
            <a
              href={barber.identityDocument.fileUrl}
              target="_blank"
              rel="noreferrer"
              className="text-blue-600 underline"
            >
              ดูเอกสาร
            </a>
          </div>
          <div className="mt-2 flex justify-center">
            <button className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
              เพิ่ม/เปลี่ยนเอกสาร
            </button>
          </div>
        </div>

        {/* ข้อมูลเงินเดือน */}
        <div className="flex flex-col justify-between h-full min-h-[200px] w-1/3  border rounded-lg p-3">
          <div>
            <h2 className="text-xl font-semibold mb-2 text-center">ข้อมูลเงินเดือน</h2>
            <p>เงินเดือนปัจจุบัน: {barber.salary.toLocaleString()} บาท</p>
          </div>
          <div className="mt-2 flex justify-center">
            <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
              แก้ไขเงินเดือน
            </button>
          </div>
        </div>

        {/* ประวัติสลิปเงินเดือน */}
        <div className="flex flex-col justify-between h-full min-h-[200px] w-1/3  border rounded-lg p-3">
          <div>
            <h2 className="text-xl font-semibold mb-2 text-center">ประวัติสลิปเงินเดือน</h2>
            <ul className="list-disc ml-6">
              {barber.payslips.map((slip) => (
                <li key={slip.id}>
                  {slip.month} -{" "}
                  <a href={slip.fileUrl} target="_blank" rel="noreferrer" className="text-blue-500 underline">
                    ดูสลิป
                  </a>
                </li>
              ))}
            </ul>
          </div>
          <div className="mt-2 flex justify-center">
            <button className="bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-600">
              เพิ่มสลิปเงินเดือน
            </button>
          </div>
        </div>
      </div>


    </div>
  );
};

export default BarberDetail;
