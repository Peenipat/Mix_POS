
import { CardItem } from "@object/shared/components/CardItem";
export default function AdminDashboard() {
  const mockEmployees = [
    {
      id: 1,
      logoSrc: "/Starbucks-Logo.jpg",
      avatarSrc: "https://randomuser.me/api/portraits/men/18.jpg",
      avatarAlt: "Portrait of Somchai",
      name: "สมชาย ใจดี",
      jobTitle: "วิศวกรซอฟต์แวร์",
      idNo: "001-1234",
      dob: "15-03-1990",
      email: "somchai.jaidee@example.com",
      phone: "+66-8123-4567",
      sex: "Male",
    },
    {
      id: 2,
      logoSrc: "/Starbucks-Logo.jpg",
      avatarSrc: "https://randomuser.me/api/portraits/women/65.jpg",
      avatarAlt: "Portrait of Sarisa",
      name: "ศิริสา พงษ์วัฒนา",
      jobTitle: "นักการตลาด",
      idNo: "002-5678",
      dob: "22-09-1993",
      email: "sarisa.pong@example.com",
      phone: "+66-9123-4567",
      sex: "Female",
    },
    {
      id: 3,
      logoSrc: "/Starbucks-Logo.jpg", // ไม่มีโลโก้
      avatarSrc: "https://randomuser.me/api/portraits/women/65.jpg",
      avatarAlt: "Portrait of Kittipong",
      name: "กิตติพงษ์ ศรีสมชัย",
      jobTitle: "ผู้ดูแลระบบ (System Admin)",
      idNo: "003-9012",
      dob: "05-12-1988",
      email: "kittipong.sri@example.com",
      phone: "+66-8222-3344",
      sex: "Male",
    },
    {
      id: 4,
      logoSrc: "/Starbucks-Logo.jpg",
      avatarSrc: "https://randomuser.me/api/portraits/women/29.jpg",
      avatarAlt: "Portrait of Napaporn",
      name: "นภาพร วิทยากร",
      jobTitle: "นักวิเคราะห์ข้อมูล",
      idNo: "004-3456",
      dob: "11-07-1995",
      email: "napaporn.wit@example.com",
      phone: "+66-8555-6677",
      sex: "Female",
    },
  ];

  const handleView = (id: number) => {
    alert(`ดูรายละเอียดพนักงาน ID: ${id}`);
  };
  const handleEdit = (id: number) => {
    alert(`แก้ไขข้อมูลพนักงาน ID: ${id}`);
  };
  const handleDelete = (id: number) => {
    if (window.confirm("ต้องการลบข้อมูลพนักงานนี้หรือไม่?")) {
      alert(`ลบพนักงาน ID: ${id} เรียบร้อยแล้ว`);
      // ตรงนี้สามารถเรียก API หรือลบจาก state ได้
    }
  };
  return (

    <>
      <h1 className="">ยินดีต้อนกลับ!</h1>
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {mockEmployees.map((emp) => (
            <CardItem
              key={emp.id}
              logoSrc={emp.logoSrc}
              avatarSrc={emp.avatarSrc}
              avatarAlt={emp.avatarAlt}
              name={emp.name}
              jobTitle={emp.jobTitle}
              // idNo={emp.idNo}
              // dob={emp.dob}
              email={emp.email}
              phone={emp.phone}
              sex={emp.sex}
              onView={() => handleView(emp.id)}
              onEdit={() => handleEdit(emp.id)}
              onDelete={() => handleDelete(emp.id)}
            />
          ))}
        </div>
      </div>
    </>
  )
}