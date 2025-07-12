import { useState, useCallback, useEffect, useRef, useMemo } from "react";
import { format } from "date-fns";
import Stepper, { Step } from "../components/Stepper";
import { useForm, } from "react-hook-form";
import { appointmentForm } from "../schemas/appointmentSchema";
// @ts-ignore
import { Datepicker } from "flowbite-datepicker";
import TimeSelector from "../components/TimeSelector";
import axios from "../lib/axios";
import type { Barber } from "../types/barber";
import { appointmentSchema } from "../schemas/appointmentSchema";
import { zodResolver } from "@hookform/resolvers/zod";

interface Service {
  id: number;
  name: string;
  description: string;
  price: number;
  duration: string;
  Img_path: string;
  Img_name: string
}

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
    setValue,
    watch,
    formState: { errors },
    trigger,
  } = useForm<appointmentForm>({
    resolver: zodResolver(appointmentSchema),
    defaultValues: {
      cusName: "",
      phoneNumber: "",
      barberId: 0,
      serviceId: 0,
      date: "",
      time: "",
      note: ""
    },
  });

  const onSubmit = (data: appointmentForm) => {
    console.log(data);
  };


  const stepFields: (keyof appointmentForm)[][] = [
    ["cusName", "phoneNumber"],
    ["barberId", "serviceId"],
    ["date", "time", "note"],
  ];
  const [currentStep, setCurrentStep] = useState(0);

  const onStepChange = async (nextStep: number) => {
    if (nextStep <= currentStep) {
      console.log('next : ', nextStep, "cur : ", currentStep)
      setCurrentStep(nextStep);
      return;
    }

    const fieldsToValidate = stepFields[currentStep];

    if (!fieldsToValidate) {
      console.warn("No fields to validate for step:", currentStep - 1);
      setCurrentStep(nextStep);
      return;
    }

    const isValid = await trigger(fieldsToValidate);

    if (isValid) {
      setCurrentStep(nextStep);
    }
    if (!isValid) {
      console.log("Validation failed at step", currentStep - 1);
      console.log(errors);
    }
  };


  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!inputRef.current) return;
    const picker = new Datepicker(inputRef.current, {
      autohide: true,
      format: "mm/dd/yyyy",
    });

    return () => {
      picker.hide();
      picker.destroy();
    };
  }, []);

  const [barbers, setBarbers] = useState<Barber[]>([]);
  const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
  const [errorBarbers, setErrorBarbers] = useState<string | null>(null);

  const loadBarbers = useCallback(async () => {
    setLoadingBarbers(true);
    setErrorBarbers(null);
    try {
      const res = await axios.get<{ status: string; data: Barber[] }>(
        `/barberbooking/branches/1/barbers`
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      setBarbers(res.data.data);
    } catch (err: any) {
      setErrorBarbers(err.response?.data?.message || err.message || "Failed to load barbers");
    } finally {
      setLoadingBarbers(false);
    }
  }, []);

  useEffect(() => {
    loadBarbers();
  }, []);

  const [services, setServices] = useState<Service[]>([]);
  const [loadingServices, setLoadingServices] = useState<boolean>(false);
  const [errorServices, setErrorServices] = useState<string | null>(null);

  const loadServices = useCallback(async () => {
    setLoadingServices(true);
    setErrorServices(null);
    try {
      const res = await axios.get<{ status: string; data: Service[] }>(
        `/barberbooking/tenants/1/branch/1/services`
      );
      if (res.data.status !== "success") {
        throw new Error(res.data.status);
      }
      setServices(res.data.data);
    } catch (err: any) {
      setErrorServices(err.response?.data?.message || err.message || "Failed to load barbers");
    } finally {
      setLoadingServices(false);
    }
  }, []);

  useEffect(() => {
    loadServices();
  }, []);

  const serviceMap = useMemo(() => {
    const map: Record<number, typeof services[0]> = {};
    services.forEach((s) => {
      map[s.id] = s;
    });
    return map;
  }, [services]);

  const barberMap = useMemo(() => {
    const map: Record<number, typeof barbers[0]> = {};
    barbers.forEach((b) => {
      map[b.id] = b;
    });
    return map;
  }, [barbers]);

  const selectedService = serviceMap[watch("serviceId")];
  const selectedBarber = barberMap[watch("barberId")];
  const isCompleted = currentStep === 4;
  const today = format(new Date(), "yyyy-MM-dd");

  return (
    <div className="h-full p-6">
      <main className="flex flex-row gap-4 h-full">
        {/* ซ้าย */}
        <div
          className="border-2 h-full w-1/2 rounded-lg overflow-auto"
        >
          <iframe
            title="ร้านของเรา"
            className="w-full h-[400px] rounded-lg shadow-lg"
            src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3875.736229316383!2d100.52973631543195!3d13.736717401097444!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x30e29edc115c79e7%3A0x1b5bdee6214e67e1!2sSiam%20Paragon!5e0!3m2!1sth!2sth!4v1688888888888!5m2!1sth!2sth"
            allowFullScreen={true}
            loading="lazy"
            referrerPolicy="no-referrer-when-downgrade"
          ></iframe>
        </div>

        {/* ขวา */}
        <div className="border-2 h-full w-1/2 rounded-lg overflow-auto">
          <div className="">
            <Stepper
              step={currentStep}
              onStepChange={onStepChange}
              nextButtonText="ถัดไป"
              backButtonText="ย้อนกลับ"
              onFinalStepCompleted={() => {
                handleSubmit(onSubmit)
                setCurrentStep(4);
              }}
            >
              <Step>
                <h2>ขั้นตอนที่ {currentStep + 1} รายละเอียดลูกค้า</h2>
                <label className="block">ชื่อลูกค้า</label>
                <input
                  type="text"
                  {...register("cusName")}
                  placeholder="กรุณากรอกชื่อลูกค้า"
                  className={`input input-bordered w-full ${errors.cusName ? "border-red-500" : ""}`}
                />
                {errors.cusName && <p className="text-red-500">{errors.cusName.message}</p>}

                <label className="block mt-4">เบอร์โทร</label>
                <input
                  type="text"
                  {...register("phoneNumber")}
                  placeholder="กรุณากรอกเบอร์โทร"
                  className={`input input-bordered w-full ${errors.phoneNumber ? "border-red-500" : ""}`}
                />
                {errors.phoneNumber && <p className="text-red-500">{errors.phoneNumber.message}</p>}
              </Step>
              <Step>
                <h2>ขั้นตอนที่ {currentStep + 1} เลือกบริการ</h2>
                <label>เลือกช่าง</label>
                {loadingBarbers && <p>Loading barbers…</p>}
                {errorBarbers && <p className="text-red-500">Error loading barbers: {errorBarbers}</p>}
                <select
                  {...register("barberId", { valueAsNumber: true })}
                  className={`select select-bordered w-full ${errors.barberId ? "border-red-500" : ""}`}
                >
                  <option value={0}>-- เลือกช่าง --</option>
                  {barbers.map((barber) => (
                    <option key={barber.id} value={barber.id}>{barber.username}</option>
                  ))}
                </select>
                {errors.barberId && <p className="text-red-500">{errors.barberId.message}</p>}

                <label className="mt-4">เลือกบริการ</label>
                {loadingServices && <p>Loading barbers…</p>}
                {errorServices && <p className="text-red-500">Error loading barbers: {errorServices}</p>}
                <select
                  {...register("serviceId", { valueAsNumber: true })}
                  className={`select select-bordered w-full ${errors.serviceId ? "border-red-500" : ""}`}
                >
                  <option value={0}>-- เลือกบริการ --</option>
                  {services.map((service) => (
                    <option key={service.id} value={service.id}>
                      {service.name} {service.price} บาท {service.duration} นาที
                    </option>
                  ))}
                </select>
                {errors.serviceId && <p className="text-red-500">{errors.serviceId.message}</p>}
              </Step>

              <Step>
                <h2>ขั้นตอนที่ {currentStep + 1} เลือกวันเวลาที่ใช้บริการ</h2>

                <label>เลือกวันที่</label>
                <input
                  type="date"
                  {...register("date")}
                  min={today}
                  className={`input input-bordered w-full ${errors.date ? "border-red-500" : ""}`}
                />
                {errors.date && <p className="text-red-500">{errors.date.message}</p>}

                <label className="">เลือกเวลา</label>
                <TimeSelector setValue={setValue} date={watch('date')} disabled={!watch("date")} />
                {errors.time && <p className="text-red-500">{errors.time.message}</p>}
                <label className="block ">ข้อความถึงช่าง</label>
                <input
                  type="text"
                  {...register("note")}
                  placeholder="ฝากข้อความถึงช่าง"
                  className={`input input-bordered w-full ${errors.cusName ? "border-red-500" : ""}`}
                />
              </Step>

              <Step>
                <h2>ตรวจสอบข้อมูล</h2>
                <p>ชื่อลูกค้า: {watch("cusName")}</p>
                <p>เบอร์โทร: {watch("phoneNumber")}</p>

                <p>บริการที่เลือก:</p>
                {selectedService ? (
                  <ul className="list-disc list-inside ">
                    <li>ชื่อบริการ: {selectedService.name}</li>
                    <li>ราคา: {selectedService.price} บาท</li>
                    <li>ระยะเวลา: {selectedService.duration} นาที</li>
                  </ul>
                ) : (
                  <p className="text-red-500">ไม่พบข้อมูล</p>
                )}

                <p>
                  ช่างที่เลือก: {selectedBarber?.username || "ไม่พบข้อมูล"}
                </p>
                <p>วันที่: {watch("date")}</p>
                <p>เวลา: {watch("time")}</p>
                <p>ฝากข้อความถึงช่าง: {watch("note")}</p>
              </Step>

              {isCompleted && (
                <div className="text-center py-10">
                  <h2 className="text-2xl font-bold text-green-600">🎉 การจองเสร็จสมบูรณ์</h2>
                  <p className="text-gray-500 mt-2">เราจะติดต่อคุณเร็ว ๆ นี้</p>
                </div>
              )}

            </Stepper>
          </div>
        </div>

      </main>
    </div>
  );
}
