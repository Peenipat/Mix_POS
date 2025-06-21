
import { CardItem } from "@object/shared/components/CardItem";
export default function AdminDashboard() {

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
    </>
  )
}