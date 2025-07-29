import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getCustomerById, CustomerDetail as CustomerDetailType } from "../../api/customer";
import { useAppSelector } from "../../store/hook";
import { GridViewIcon } from "../../components/icons/GridViewIcon";
import { TableCellsIcon } from "../../components/icons/TableCellsIcon";
import { DataTable } from "../../components/DataTable";

export interface Appointment {
    id: number;
    date: string;         // "YYYY-MM-DD"
    time: string;         // "HH:mm"
    customerName: string;
    isMember: boolean;
    barberName: string;
    serviceName: string;
    price: number;
    duration: number;     // minutes
  }

  const mockAppointments: Appointment[] = [
    {
      id: 1,
      date: "2025-07-30",
      time: "10:00",
      customerName: "สมชาย ใจดี",
      isMember: false,
      barberName: "ช่างบอล",
      serviceName: "ตัดผมชาย",
      price: 150,
      duration: 30,
    },
    {
      id: 2,
      date: "2025-08-01",
      time: "13:30",
      customerName: "สมชาย ใจดี",
      isMember: false,
      barberName: "ช่างกอล์ฟ",
      serviceName: "สระผม + ไดร์",
      price: 200,
      duration: 40,
    },
    {
      id: 3,
      date: "2025-08-10",
      time: "11:15",
      customerName: "สมชาย ใจดี",
      isMember: false,
      barberName: "ช่างปุ้ย",
      serviceName: "โกนหนวด",
      price: 100,
      duration: 15,
    },
  ];

  
export function CustomerDetail() {
    const { id } = useParams<{ id: string }>();
    const me = useAppSelector((state) => state.auth.me);
    const tenantId = me?.tenant_ids[0];
    const branchId = me?.branch_id;

    const [customer, setCustomer] = useState<CustomerDetailType | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [viewMode, setViewMode] = useState<"table" | "card">("card");
    const [appointments, setAppointments] = useState<Appointment[]>(mockAppointments);

    useEffect(() => {
        if (!tenantId || !branchId || !id) return;

        const loadCustomer = async () => {
            setLoading(true);
            try {
                const customer = await getCustomerById(Number(tenantId), Number(branchId), Number(id));
                setCustomer(customer);
            } catch (err) {
                setError("ไม่สามารถโหลดข้อมูลลูกค้าได้");
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        loadCustomer();
    }, [tenantId, branchId, id]);

    if (loading) return <p>กำลังโหลดข้อมูลลูกค้า...</p>;
    if (error) return <p className="text-red-500">{error}</p>;
    if (!customer) return <p>ไม่พบข้อมูลลูกค้า</p>;

    return (
        <div className="p-4 max-w-5xl mx-auto space-y-6">
            {/* ข้อมูลลูกค้า */}
            <div className="bg-white shadow rounded-md p-4 max-w-md">
                <h2 className="text-xl font-bold mb-4">ข้อมูลลูกค้า</h2>
                <div className="space-y-2">
                    <div>
                        <span className="font-semibold">ชื่อ:</span> {customer.name}
                    </div>
                    <div>
                        <span className="font-semibold">เบอร์โทร:</span> {customer.phone}
                    </div>
                    <div>
                        <span className="font-semibold">อีเมล:</span> {customer.email}
                    </div>
                </div>
            </div>

            {/* ปุ่มสลับมุมมอง */}
            <div className="flex justify-between items-center">
                <h2 className="text-xl font-bold">ประวัติการจอง</h2>
                <div className="flex gap-2">
                    <button
                        onClick={() => setViewMode("card")}
                        className={`px-4 py-2 rounded-md border ${viewMode === "card"
                            ? "bg-blue-600 text-white"
                            : "bg-white text-gray-800 border-gray-300"
                            }`}
                    >
                        <GridViewIcon className="w-5 h-5" />
                    </button>
                    <button
                        onClick={() => setViewMode("table")}
                        className={`px-4 py-2 rounded-md border ${viewMode === "table"
                            ? "bg-blue-600 text-white"
                            : "bg-white text-gray-800 border-gray-300"
                            }`}
                    >
                        <TableCellsIcon className="w-5 h-5" />
                    </button>
                </div>
            </div>

            {/* ตารางหรือใบจอง */}
            {viewMode === "table" ? (
                <DataTable<Appointment>
                    data={appointments}
                    columns={[
                        { header: "#", accessor: (_r, i) => i + 1 },
                        { header: "วัน", accessor: "date" },
                        { header: "เวลา", accessor: "time" },
                        { header: "บริการ", accessor: "serviceName" },
                        { header: "ราคา", accessor: "price" },
                        { header: "ระยะเวลา", accessor: "duration" },
                        { header: "ช่าง", accessor: "barberName" },
                    ]}
                    showEdit={false}
                    showDelete={false}
                />
            ) : (
                <div className="grid md:grid-cols-2 gap-4">
                    {appointments.map((bookingInfo) => (
                        <div key={bookingInfo.id} className="w-full flex flex-col gap-3 p-4 border rounded-md shadow">
                            <h2 className="text-lg font-semibold text-center">ใบจอง</h2>

                            <div className="flex justify-between text-base px-2">
                                <p>วันที่จอง: <span className="font-medium">{bookingInfo.date}</span></p>
                                <p>เวลาที่จอง: <span className="font-medium">{bookingInfo.time}</span></p>
                            </div>

                            <div className="grid grid-cols-2 gap-y-1 text-base px-2 mt-0">
                                <p>ชื่อลูกค้า:</p>
                                <p className="text-right">{bookingInfo.customerName}</p>

                                <p>สมาชิก:</p>
                                <p className="text-right">{bookingInfo.isMember ? "เป็นสมาชิก" : "ไม่ได้เป็น"}</p>

                                <p>ช่างที่เลือก:</p>
                                <p className="text-right">{bookingInfo.barberName}</p>

                                <p>บริการที่เลือก:</p>
                                <p className="text-right">{bookingInfo.serviceName}</p>

                                <p>ราคา:</p>
                                <p className="text-right">{bookingInfo.price} บาท</p>

                                <p>เวลาโดยประมาณ:</p>
                                <p className="text-right">{bookingInfo.duration} นาที</p>
                            </div>

                            <div className="text-center">
                                <p>ขอบคุณที่ใช้บริการ</p>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );

}
