// BarberDetailAdmin.tsx
import { useAppSelector } from "@object/shared/store/hook";
import { getBarberById, updateBarber } from "../../../api/barber";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import type { EditBarberFormData } from "../../../schemas/barberSchema";
import { editBarberSchema } from "../../../schemas/barberSchema";
import type { BarberDetail } from "../../../api/barber"
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

const BarberDetail = () => {
  const { id } = useParams<{ id: string }>();
  const [barber, setBarber] = useState<BarberDetail | null>(null);
  useEffect(() => {
    async function fetchBarber() {
      try {
        const barber = await getBarberById(Number(id));
        setBarber(barber)
      } catch (err) {
        console.error("โหลดข้อมูล Barber ล้มเหลว", err);
      }
    }

    fetchBarber();
  }, []);

  const {
    register,
    handleSubmit,
    reset,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<EditBarberFormData>({
    resolver: zodResolver(editBarberSchema),
    defaultValues: {
      username: "",
      description: "",
      email: "",
      phone_number: "",
      img_path: "",
      img_name: "",
      roleUser: "",
      profilePicture: "",
    },
  });
  console.log(barber)

  useEffect(() => {
    if (barber) {
      reset({
        username: barber.user.username,
        email: barber.user.email,
        description: barber.description,
        phone_number: barber.user.phone_number,
        img_path: barber.user.Img_path,
        img_name: barber.user.Img_name,
        roleUser: barber.role_user,
        branch_id: barber.branch_id,
      });
    }
  }, [barber, reset]);
  const [previewUrl, setPreviewUrl] = useState<string | null>(`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${barber?.user.Img_path}/${barber?.user.Img_name}`);

  useEffect(() => {
    if (barber?.user?.Img_path && barber?.user?.Img_name) {
      setPreviewUrl(
        `https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${barber.user.Img_path}/${barber.user.Img_name}`
      );
    }
  }, [barber]);

  const onSubmit = async (data: EditBarberFormData) => {
    console.log("data : ",data)
    try {
      const formData = new FormData();

      if (data.profilePicture instanceof File || typeof data.profilePicture === "object") {
        formData.append("profilePicture", data.profilePicture[0]);
      }
      formData.append("branch_id", data.branch_id.toString());
      formData.append("username", data.username);
      formData.append("description", data.description ?? "");
      formData.append("email", data.email);
      formData.append("phone_number", data.phone_number);
      formData.append("img_path", data.img_path ?? "");
      formData.append("img_name", data.img_name ?? "");
      formData.append("role_user", data.roleUser ?? "");

      if (barber?.tenant_id) {
        await updateBarber(Number(barber?.tenant_id), Number(id), formData);
      }


    } catch (err) {
      console.error("อัปเดตล้มเหลว", err);
    }
  };



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
  const barberMock = mockBarberWithAdminData;

  return (
    <div className="">
      <h1 className="text-xl font-bold mb-3">แก้ไขข้อมูลของช่าง</h1>
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="mb-4">
          <label className="block text-gray-700 dark:text-gray-200 mb-1">
            รูปภาพโปรไฟล์
          </label>
          <input
            type="file"
            accept="image/*"
            {...register("profilePicture")}
            className="file-input file-input-bordered w-full"
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              const file = e.target.files?.[0];
              if (file) {
                const url = URL.createObjectURL(file);
                setPreviewUrl(url);
              }
            }}
          />
          {/* {errors.profilePicture && (
            <p className="text-red-600 text-sm mt-1">
              {errors.profilePicture.message}
            </p>
          )} */}
          {previewUrl && previewUrl !== "" && (
            <div className="mt-3">
              <img
                src={previewUrl}
                alt="Preview"
                className="h-20 w-24 object-cover rounded-md border object-top"
              />
            </div>
          )}
        </div>
        <div className="flex w-full gap-3">
          <input type="hidden" {...register("branch_id")} />
          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              ชื่อช่าง
            </label>
            <input
              type="text"
              {...register("username")}
              className={`w-full input input-bordered ${errors.username ? "border-red-500" : ""
                }`}
            />
            {errors.username && (
              <p className="text-red-600 text-sm mt-1">
                {errors.username.message}
              </p>
            )}
          </div>

          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              คำอธิบาย
            </label>
            <input
              type="text"
              placeholder="กรอกคำอธิบายช่าง"
              {...register("description")}
              className={`w-full input input-bordered ${errors.description ? "border-red-500" : ""
                }`}
            />
            {errors.description && (
              <p className="text-red-600 text-sm mt-1">
                {errors.description.message}
              </p>
            )}
          </div>
        </div>


        <h3 className="text-lg text-red-500">** ข้อมูลภายใน **</h3>
        <div className="flex gap-3">
          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              อีเมล์
            </label>
            <input
              type="email"
              {...register("email")}
              readOnly
              className={`w-full input input-bordered bg-gray-200 ${errors.email ? "border-red-500" : ""}`}
            />
            {errors.email && (
              <p className="text-red-600 text-sm mt-1">
                {errors.email.message}
              </p>
            )}
          </div>

          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              ชื่อตำแหน่ง
            </label>
            <input
              type="text"
              placeholder="กรอกชื่อตำแหน่ง"
              {...register("roleUser")}
              className={`w-full input input-bordered ${errors.roleUser ? "border-red-500" : ""
                }`}
            />
            {errors.roleUser && (
              <p className="text-red-600 text-sm mt-1">
                {errors.roleUser.message}
              </p>
            )}
          </div>

          {/* Phone Number */}
          <div className="mb-4 w-full">
            <label className="block text-gray-700 dark:text-gray-200 mb-1">
              เบอร์โทร
            </label>
            <input
              type="text"
              {...register("phone_number")}
              placeholder="0812345678"
              className={`w-full input input-bordered ${errors.phone_number ? "border-red-500" : ""
                }`}
            />
            {errors.phone_number && (
              <p className="text-red-600 text-sm mt-1">
                {errors.phone_number.message}
              </p>
            )}
          </div>
        </div>
        <div className="flex justify-end mb-4">
          <button className="bg-green-500 text-white p-1.5 rounded-md justify-end">บันทึก</button>
        </div>

      </form>

      <h3 className="text-red-500">** ส่วนนี้ยังไม่พร้อมใช้งาน **</h3>
      <div className="flex justify-between gap-6">

        {/* กล่องที่ 1: ข้อมูลประจำตัว */}
        <div className="flex flex-col justify-between h-full min-h-[200px] w-1/3 border rounded-lg p-3 opacity-50 pointer-events-none cursor-not-allowed">
          <div>
            <h2 className="text-xl font-semibold mb-2 text-center">ข้อมูลประจำตัว</h2>
            <p>ประเภทเอกสาร: {barberMock.identityDocument.type}</p>
            <p>หมายเลข: {barberMock.identityDocument.number}</p>
            <a
              href={barberMock.identityDocument.fileUrl}
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

        {/* กล่องที่ 2: ข้อมูลเงินเดือน */}
        <div className="flex flex-col justify-between h-full min-h-[200px] w-1/3 border rounded-lg p-3 opacity-50 pointer-events-none cursor-not-allowed">
          <div>
            <h2 className="text-xl font-semibold mb-2 text-center">ข้อมูลเงินเดือน</h2>
            <p>เงินเดือนปัจจุบัน: {barberMock.salary.toLocaleString()} บาท</p>
          </div>
          <div className="mt-2 flex justify-center">
            <button className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
              แก้ไขเงินเดือน
            </button>
          </div>
        </div>

        {/* กล่องที่ 3: ประวัติสลิปเงินเดือน */}
        <div className="flex flex-col justify-between h-full min-h-[200px] w-1/3 border rounded-lg p-3 opacity-50 pointer-events-none cursor-not-allowed">
          <div>
            <h2 className="text-xl font-semibold mb-2 text-center">ประวัติสลิปเงินเดือน</h2>
            <ul className="list-disc ml-6">
              {barberMock.payslips.map((slip) => (
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
