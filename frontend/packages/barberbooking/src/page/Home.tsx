// src/page/Home.tsx
import React, { useEffect, useRef } from "react";
import Navbar from "../components/Navbar";
import GridMotion from "../components/GridMotion";
import Stepper, { Step } from "../components/Stepper";
import { useForm } from "react-hook-form";
import { EditBarberFormData } from "../schemas/barberSchema";
import { useState } from 'react'
// @ts-ignore
import { Datepicker } from "flowbite-datepicker";
import TimeSelector from "../components/TimeSelector";

export default function Home() {
  const items: React.ReactNode[] = [
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber3.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber4.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber5.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber6.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber7.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber8.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber3.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber4.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber5.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber6.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber7.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber8.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber10.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber9.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber2.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber8.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber3.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber4.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber5.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber6.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber7.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber8.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber5.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber6.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber7.jpg",
    "https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/imgHome/barber8.jpg",
  ];
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<EditBarberFormData>({
    defaultValues: {
      username: "",
      email: "",
      phone_number: "",
    },
  });
  interface User {
    user_id: number;
    username: string;
    email: string;
  }

  const [users, setUsers] = useState<User[]>([]);

  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!inputRef.current) return;

    // สร้าง instance ของ Datepicker
    const picker = new Datepicker(inputRef.current, {
      autohide: true,
      // คุณสามารถใส่ option เพิ่มได้ เช่น format, minDate, maxDate ฯลฯ
      format: "mm/dd/yyyy",
    });

    // ถ้าต้องการ cleanup เมื่อ component unmount
    return () => {
      picker.hide();
      picker.destroy();
    };
  }, []);

  const [selectedTime, setSelectedTime] = useState("00:00");

  return (
    <div className="relative w-screen h-screen overflow-hidden">
      {/* === GridMotion as animated “background” === */}
      <div className="absolute inset-0 z-0">
        <GridMotion items={items} gradientColor="rgba(0,0,0,0.5)" />
      </div>

      {/* === Foreground content === */}
      <div className="relative z-10 flex flex-col h-full">
        <header>
          <Navbar />
        </header>

        {/* แบ่ง 2 คอลัมน์ ซ้าย–ขวา */}
        <main className="flex-1 flex">
          {/* ซ้าย */}
          <div
            className="
              flex-1 m-4 p-6
              bg-white/20           /* ครึ่งโปร่งแสง */
              backdrop-blur-md      /* เบลอภาพหลัง */
              rounded-lg
              flex items-center justify-center
            "
          >
            <h2 className="text-3xl font-semibold text-white text-center">
              Hello, Left!
            </h2>
          </div>

          {/* ขวา */}
          <div className=" flex-1 m-4 p-6 bg-white/50 rounded-lg">
            <div className="border-2 min-h-screen  rounded-lg ">
              <Stepper
                initialStep={1}
                onStepChange={(step) => console.log(step)}
                onFinalStepCompleted={() => console.log("All steps completed!")}
                backButtonText="Previous"
                nextButtonText="Next"
              >
                <Step>
                  <h2>ขั้นตอนที่ 1 รายละเอียดลูกค้า</h2>
                  <label className="block text-black dark:text-gray-200 mb-1">
                    ชื่อลูกค้า
                  </label>
                  <input
                    type="text"
                    {...register("username")}
                    placeholder="กรอกชื่อลูกค้า"
                    className={`w-full input input-bordered ${errors.username ? "border-red-500" : ""
                      }`}
                  />

                  <label className="block text-black dark:text-gray-200 mb-1">
                    เบอร์โทร
                  </label>
                  <input
                    type="text"
                    placeholder="กรอกเบอร์โทร"
                    {...register("username")}
                    className={`w-full input input-bordered ${errors.username ? "border-red-500" : ""
                      }`}
                  />
                </Step>
                <Step>
                  <h2>ขั้นตอนที่ 2 เลือกบริการ</h2>
                  <label htmlFor="user" className="block  dark:text-gray-200 mb-1">
                    เลือกช่าง
                  </label>
                  <select
                    id="user"
                    // value={selectedUserId || undefined}
                    onChange={(e) => {
                      const v = e.target.value;
                      // setSelectedUserId(v === "" ? "" : Number(v));
                    }}
                    className="w-full select select-bordered"
                  >
                    {users.length != 0 ? (
                      <option value="">-- เลือกผู้ใช้ --</option>
                    ) : <option value="">-- ไม่พบข้อมูล --</option>}

                    {users.map((u) => (
                      <option key={u.user_id} value={u.user_id}>
                        {u.username} ({u.email})
                      </option>
                    ))}
                  </select>

                  <label htmlFor="user" className="block  dark:text-gray-200 mb-1">
                    เลือกบริการ
                  </label>
                  <select
                    id="user"
                    // value={selectedUserId || undefined}
                    onChange={(e) => {
                      const v = e.target.value;
                      // setSelectedUserId(v === "" ? "" : Number(v));
                    }}
                    className="w-full select select-bordered"
                  >
                    {users.length != 0 ? (
                      <option value="">-- เลือกผู้ใช้ --</option>
                    ) : <option value="">-- ไม่พบข้อมูล --</option>}

                    {users.map((u) => (
                      <option key={u.user_id} value={u.user_id}>
                        {u.username} ({u.email})
                      </option>
                    ))}
                  </select>

                </Step>
                <Step>
                  <h2>ขั้นตอนที่ 3 เลือกวันเวลาที่ใช้บริการ</h2>

                  <div className="relative max-w-sm mx-auto">
                    <div className="absolute inset-y-0  flex items-center pl-3 pointer-events-none">
                      <svg
                        className="w-4 h-4 text-gray-500 dark:text-gray-400"
                        aria-hidden="true"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path d="M20 4a2 2 0 0 0-2-2h-2V1a1 1 0 0 0-2 0v1h-3V1a1 1 0 0 0-2 0v1H6V1a1 1 0 0 0-2 0v1H2a2 2 0 0 0-2 2v2h20V4ZM0 18a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V8H0v10Zm5-8h10a1 1 0 0 1 0 2H5a1 1 0 0 1 0-2Z" />
                      </svg>
                    </div>
                    <input
                      ref={inputRef}
                      id="datepicker-autohide"
                      type="text"
                      className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full pl-10 p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                      placeholder="Select date"
                    />
                  </div>
                  <TimeSelector />
                </Step>
                <Step>
                  <h2>ตรวจสอบข้อมูล</h2>
                  <h2>ชื่อลูกค้า : </h2>
                  <h2>เบอร์โทรลูกค้า : </h2>
                  <h2>บริการที่เลือก : </h2>
                  <h2>ช่างที่เลือก : </h2>
                  <h2>วันที่ใช้บริการ : </h2>
                  <h2>เวลาที่ใช้บริการ : </h2>
                </Step>
              </Stepper>
            </div>
          </div>

        </main>
      </div>
    </div>
  );
}
