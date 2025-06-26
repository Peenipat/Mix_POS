import {
  BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer,
  LineChart, Line, CartesianGrid, Legend,
} from 'recharts';

export default function AdminDashboard() {
  const summary = {
    todayBookings: 5,
    monthBookings: 120,
    todayRevenue: 1500,
    monthRevenue: 32000,
    todayExpenses: 400,
    newCustomers: 3,
  };

  const upcomingAppointments = [
    { time: "10:00", customer: "สมชาย", barber: "ช่างเอก", service: "ตัดผมชาย" },
    { time: "11:30", customer: "จันทร์เพ็ญ", barber: "ช่างหญิง", service: "สระไดร์" },
  ];

  const popularServices = [
    { name: "ตัดผมชาย", count: 45 },
    { name: "สระไดร์", count: 30 },
    { name: "โกนหนวด", count: 12 },
  ];

  const chartData = [
    { date: "15 มิ.ย.", bookings: 12, revenue: 2900, expense: 890 },
    { date: "16 มิ.ย.", bookings: 1, revenue: 2400, expense: 890 },
    { date: "17 มิ.ย.", bookings: 19, revenue: 2500, expense: 890 },
    { date: "18 มิ.ย.", bookings: 3, revenue: 3000, expense: 2890 },
    { date: "19 มิ.ย.", bookings: 10, revenue: 1800, expense: 850 },
    { date: "20 มิ.ย.", bookings: 9, revenue: 1800, expense: 850 },
    { date: "21 มิ.ย.", bookings: 14, revenue: 2800, expense: 1700 },
    { date: "22 มิ.ย.", bookings: 11, revenue: 2200, expense: 1000 },
    { date: "23 มิ.ย.", bookings: 30, revenue: 3000, expense: 1200 },
    { date: "24 มิ.ย.", bookings: 8, revenue: 1600, expense: 1700 },
    { date: "25 มิ.ย.", bookings: 10, revenue: 2000, expense: 100 },
    { date: "26 มิ.ย.", bookings: 15, revenue: 2400, expense: 800 },
    { date: "27 มิ.ย.", bookings: 15, revenue: 2400, expense: 800 },
  ];

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">ยินดีต้อนรับกลับ!</h1>

      {/* ตัวเลขสรุป */}
      <div className="grid grid-cols-2 md:grid-cols-2 gap-4">
 
        <Card title="รายได้วันนี้" value={`฿${summary.todayRevenue}`} variant="success" />
        
        <Card title="รายจ่ายวันนี้" value={`฿${summary.todayExpenses}`} variant="danger" />
        <Card title="ยอดจองวันนี้" value={summary.todayBookings} />
        <Card title="ยอดจองเดือนนี้" value={summary.monthBookings} />
        {/* <Card title="รายได้เดือนนี้" value={`฿${summary.monthRevenue}`} variant="success" /> */}
        {/* <Card title="ลูกค้าใหม่วันนี้" value={summary.newCustomers} /> */}
      </div>

      {/* กราฟยอดจอง */}
      <section>
        <h2 className="text-lg font-semibold mb-2">ยอดจองรายวัน</h2>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={chartData}>
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Bar dataKey="bookings" fill="#3b82f6" name="ยอดจอง" />
          </BarChart>
        </ResponsiveContainer>
      </section>

      {/* กราฟรายได้-รายจ่าย */}
      <section>
        <h2 className="text-lg font-semibold mb-2">รายได้ & รายจ่าย</h2>
        <ResponsiveContainer width="100%" height={300}>
          <LineChart data={chartData}>
            <XAxis dataKey="date" />
            <YAxis />
            <CartesianGrid strokeDasharray="3 3" />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="revenue" stroke="#10b981" name="รายได้" />
            <Line type="monotone" dataKey="expense" stroke="#ef4444" name="รายจ่าย" />
          </LineChart>
        </ResponsiveContainer>
      </section>

      {/* นัดหมายวันนี้ */}
      <section>
        <h2 className="text-lg font-semibold mb-2">นัดหมายวันนี้</h2>
        <ul className="space-y-2">
          {upcomingAppointments.map((a, i) => (
            <li key={i} className="border p-3 rounded shadow-sm">
              <strong>{a.time}</strong> - {a.customer} ({a.service}) โดย {a.barber}
            </li>
          ))}
        </ul>
      </section>

      {/* บริการยอดนิยม */}
      <section>
        <h2 className="text-lg font-semibold mb-2">บริการยอดนิยม</h2>
        <ul className="list-disc list-inside text-gray-700">
          {popularServices.map((s, i) => (
            <li key={i}>{s.name} ({s.count} ครั้ง)</li>
          ))}
        </ul>
      </section>
    </div>
  );
}
function Card({
  title,
  value,
  variant = "default", 
}: {
  title: string;
  value: string | number;
  variant?: "default" | "success" | "danger";
}) {
  const borderColor =
    variant === "success"
      ? "border-green-500"
      : variant === "danger"
      ? "border-red-500"
      : "border-gray-200";

  const bgColor =
    variant === "success"
      ? "bg-green-50"
      : variant === "danger"
      ? "bg-red-50"
      : "bg-white";

  const textColor =
    variant === "success"
      ? "text-green-600"
      : variant === "danger"
      ? "text-red-600"
      : "text-gray-800";

  const titleColor =
    variant === "success"
      ? "text-green-500"
      : variant === "danger"
      ? "text-red-500"
      : "text-gray-500";

  return (
    <div className={`rounded shadow p-4 border ${borderColor} ${bgColor}`}>
      <h3 className={`text-sm font-medium ${titleColor}`}>{title}</h3>
      <p className={`text-2xl font-bold ${textColor}`}>{value}</p>
    </div>
  );
}
