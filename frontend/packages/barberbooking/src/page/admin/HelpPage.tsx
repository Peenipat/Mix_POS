import React from "react";

export default function HelpPage() {
  return (
    <div className="max-w-4xl mx-auto py-10 px-4">
      <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white">
        ความช่วยเหลือในการใช้งานระบบ
      </h1>

      <div className="space-y-6">
        <div>
          <h2 className="text-xl font-semibold text-gray-700 dark:text-white">วิธีการสร้างการนัดหมาย</h2>
          <p className="text-gray-600 dark:text-gray-300 mt-1">
            ไปที่เมนู <strong>การนัดหมาย</strong> จากนั้นคลิกปุ่ม “เพิ่มการนัดหมาย” แล้วเลือกบริการ, วันที่, และช่างที่ต้องการ
          </p>
        </div>

        <div>
          <h2 className="text-xl font-semibold text-gray-700 dark:text-white">การจัดการข้อมูลช่าง</h2>
          <p className="text-gray-600 dark:text-gray-300 mt-1">
            เข้าสู่เมนู <strong>ข้อมูลช่าง</strong> เพื่อเพิ่ม / แก้ไข / ลบช่าง โดยระบุชื่อ, รูปภาพ, ความสามารถ และเวลาทำงาน
          </p>
        </div>

        <div>
          <h2 className="text-xl font-semibold text-gray-700 dark:text-white">การตั้งเวลาเปิด - ปิดร้าน</h2>
          <p className="text-gray-600 dark:text-gray-300 mt-1">
            ในเมนู <strong>เวลาทำการ</strong> คุณสามารถกำหนดวัน-เวลาที่ร้านเปิด และตั้งค่าช่างแต่ละคนให้ไม่ซ้ำกันได้
          </p>
        </div>

        <div>
          <h2 className="text-xl font-semibold text-gray-700 dark:text-white">เคล็ดลับ</h2>
          <ul className="list-disc list-inside text-gray-600 dark:text-gray-300 mt-1 space-y-1">
            <li>คุณสามารถเปิด/ปิดสถานะของบริการได้โดยไม่ต้องลบ</li>
            <li>ระบบจะเตือนหากเวลานัดทับซ้อนกับการนัดเดิม</li>
            <li>แนะนำให้ใส่รูปโปรไฟล์ช่าง เพื่อให้ลูกค้าเลือกได้ง่ายขึ้น</li>
          </ul>
        </div>

        <div className="mt-8 text-sm text-gray-500 dark:text-gray-400">
          หากต้องการติดต่อผู้พัฒนา กรุณาไปที่หน้า <strong>“ติดต่อผู้พัฒนา”</strong> จากเมนูด้านซ้าย
        </div>
      </div>
    </div>
  );
}
