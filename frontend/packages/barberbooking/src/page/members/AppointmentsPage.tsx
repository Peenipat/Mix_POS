import Modal from "@object/shared/components/Modal";
import { FiInfo } from "react-icons/fi";
import { useState, useCallback, useEffect } from "react";
import axios from "../../lib/axios";
import { Toast } from "@object/shared/components/Toast";
import { Service } from "../ServicePage";
import { Barber } from "../../types/barber";
import { TotalBarberSchedule } from "../BarberPage";
import "intro.js/introjs.css";
import introJs from "intro.js";
import { Divide } from "lucide-react";
export interface Appointment {
    barberId: number;
    start: string; // "HH:mm"
    end: string;   // "HH:mm"
    date: string;  // "yyyy-MM-dd"
}

export const mockAppointments: Appointment[] = [
    {
        barberId: 1,
        start: "09:15",
        end: "10:30",
        date: "2025-07-12",
    },
    {
        barberId: 2,
        start: "09:30",
        end: "10:45",
        date: "2025-07-12",
    },
];

export default function AppointmentsPage() {
    const [isModalOpen, setIsModalOpen] = useState(false);

    const [barbers, setBarbers] = useState<Barber[]>([]);
    const [loadingBarbers, setLoadingBarbers] = useState<boolean>(false);
    const [errorBarbers, setErrorBarbers] = useState<string | null>(null);
    const [selectedServiceId, setSelectedServiceId] = useState<Number | null>(null)

    const [services, setServices] = useState<Service[]>([]);
    const [loadingServices, setLoadingServices] = useState<boolean>(false);
    const [errorServices, setErrorServices] = useState<string | null>(null);

    const loadServices = async () => {
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
            setErrorServices(
                err.response?.data?.message || err.message || "Failed to load services"
            );
        } finally {
            setLoadingServices(false);
        }
    };

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
    useEffect(() => {
        loadServices();
    }, []);



    const [selectedBooking, setSelectedBooking] = useState<{
        date: string;
        time: string;
        barberId: number;
        barberName: string;
    } | null>(null);

    const [countdown, setCountdown] = useState(0);

    const handleOpenModal = (barberId: number, time: string) => {
        const barber = barbers.find((b) => b.id === barberId);
        if (!barber) return;

        setSelectedBooking({
            date: new Date().toISOString().slice(0, 10),
            time,
            barberId,
            barberName: barber.username,
        });
        setCountdown(420);
        setIsModalOpen(true);
    };

    const handleClose = () => {
        setIsModalOpen(false);
        setSelectedBooking(null);
        setCountdown(0);
    };

    useEffect(() => {
        if (!isModalOpen || countdown <= 0) return;

        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    handleClose();
                    setToastInfo({
                        message: "คุณทำรายการเกินระยะเวลาที่กำหนด ",
                        variant: "error",
                    });
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => clearInterval(timer);
    }, [isModalOpen, countdown]);


    const formatTime = (seconds: number) => {
        const min = Math.floor(seconds / 60)
            .toString()
            .padStart(2, "0");
        const sec = (seconds % 60).toString().padStart(2, "0");
        return `${min}:${sec}`;
    };

    const [toastInfo, setToastInfo] = useState<{
        message: string;
        variant: "success" | "error";
    } | null>(null);

    const handleBookingSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        handleClose();
        setToastInfo({ message: "จองคิวสำเร็จเรียบร้อย!", variant: "success" });
    };

    const bookingStepsMock = [
        {
            step: 1,
            title: "เลือกวันที่และช่วงเวลาที่ว่าง",
            description:
                "ใช้ปฏิทินเพื่อเลือกวันที่ต้องการจอง และดูเวลาที่ว่างในแต่ละช่วง เพื่อเลือกเวลาที่สะดวกสำหรับคุณ",
        },
        {
            step: 2,
            title: "เลือกช่างที่ต้องการ",
            description:
                "ดูชื่อช่างและช่วงเวลาว่าง แล้วเลือกช่างที่คุณต้องการให้บริการในเวลาที่สะดวก",
        },
        {
            step: 3,
            title: "กรอกข้อมูลผู้จอง",
            description:
                "กรอกชื่อ-นามสกุล และเบอร์โทรศัพท์สำหรับติดต่อ เพื่อใช้ยืนยันการจองคิว",
        },
        {
            step: 4,
            title: "ยืนยันการจองคิว",
            description:
                "ตรวจสอบรายละเอียดอีกครั้งก่อนยืนยัน ระบบจะทำการล็อกคิวไว้ชั่วคราวให้คุณ 7 นาที",
        },
        {
            step: 5,
            title: "ชำระเงินผ่าน QR Code",
            description:
                "หลังจากยืนยันแล้ว ระบบจะแสดงใบจองพร้อม QR Code สำหรับชำระเงิน กรุณาชำระภายในเวลาที่กำหนดเพื่อยืนยันการจอง",
        },
        {
            step: 6,
            title: "รอรับ SMS หรือแจ้งเตือน",
            description:
                "เมื่อการชำระเงินเสร็จสิ้น คุณจะได้รับ SMS หรือการแจ้งเตือนยืนยันการจองทันที",
        },
    ];
    const [showSteps, setShowSteps] = useState(false);

    const [hideForm, setHideForm] = useState(false)
    const [hideConfirm, setHideConfirm] = useState(false)
    const [step, setStep] = useState(1)

    const [showModalTour, setShowModalTour] = useState(false)
    const [isHideMock, setIsHideMock] = useState(false)
    const handleCloseModalTour = () => {
        setShowModalTour(false);
    };
    const handleStartIntro = async () => {
        const bookingEl = document.getElementById("booking-mock");
        if (bookingEl) {
            bookingEl.click();
            setShowModalTour(true)
            setIsHideMock(true)
        }

        const waitForModal = () =>
            new Promise<void>((resolve) => {
                const checkExist = () => {
                    const el = document.getElementById("modal-booking-form");
                    if (el && el.offsetParent !== null) {
                        resolve();
                    } else {
                        setTimeout(checkExist, 100);
                    }
                };
                checkExist();
            });

        await waitForModal();

        const intro = introJs();

        intro.setOptions({
            steps: [
                { intro: "ขอต้อนรับเข้าสู่ระบบแนะนำการใช้งานสำหรับการจองคิว" },
                { element: "#step-booking-button", intro: "เริ่มต้นโดยคลิกที่ปุ่มนี้เพื่อดูขั้นตอนการจอง" },
                { element: "#step-booking-detail", intro: "คุณสามารถดูรายละเอียดของขั้นตอนการจองได้ที่นี่" },
                { intro: "ถัดไป เราจะพาคุณไปรู้จักส่วนต่าง ๆ ที่ใช้ในการจองคิวในระบบ" },
                { element: "#date-picker", intro: "เริ่มจากการเลือกวันที่คุณต้องการจองคิว" },
                { element: "#filter-button", intro: "คลิกปุ่มนี้เพื่อระบุช่วงเวลาที่คุณสะดวก" },
                { element: "#all-option", intro: "คุณสามารถเลือกตัวกรองเวลาได้ตามต้องการ" },
                { element: "#select-option", intro: "ระบบมีตัวเลือกให้เลือกช่วงสัปดาห์หรือเดือนที่ต้องการจอง" },
                { element: "#time-option", intro: "จากนั้นเลือกช่วงเวลาที่คุณว่าง เพื่อกรองผลลัพธ์ให้ตรงกับเวลาที่ต้องการ" },
                { element: "#search-button", intro: "คลิกปุ่มค้นหาเพื่อให้ระบบแสดงผลตามเงื่อนไขที่คุณกำหนด" },
                { element: "#data-box", intro: "ผลลัพธ์จะปรากฏในรูปแบบตารางดังนี้" },
                { element: "#barber-name", intro: "คอลัมน์แรกแสดงชื่อของช่างแต่ละคน" },
                { element: "#barber-time", intro: "คอลัมน์ถัดมาจะแสดงเวลาที่ร้านเปิดทำการ พร้อมแสดงเวลาที่ช่างแต่ละคนว่าง" },
                { element: "#booking", intro: "คุณสามารถคลิกที่ช่อง 'ว่าง' เพื่อเริ่มการจองคิวได้ทันที" },
                { element: "#modal-booking-form", intro: "เมื่อคลิกแล้ว ระบบจะแสดงหน้าต่างสำหรับกรอกรายละเอียดการจองเพิ่มเติม" },
                { element: "#modal-fix", intro: "ระบบจะระบุวัน เวลา และชื่อช่างให้อัตโนมัติ ไม่สามารถแก้ไขได้ หากต้องการเปลี่ยน กรุณาเริ่มเลือกใหม่" },
                { element: "#modal-time", intro: "ระบบจะล็อกคิวให้คุณเป็นเวลา 7 นาที กรุณาดำเนินการให้เสร็จภายในเวลานี้ มิฉะนั้นคิวจะถูกยกเลิก" },
                { element: "#modal-form", intro: "กรอกข้อมูลของคุณ เช่น ชื่อ และเบอร์โทรศัพท์ เพื่อดำเนินการต่อ" },
                { element: "#comfirm-booking", intro: "เมื่อกรอกข้อมูลครบถ้วน คลิกปุ่มนี้เพื่อดำเนินการต่อ" },
                { element: "#modal-comfirm", intro: "ตรวจสอบข้อมูลของคุณ" },
                { element: "#pay-button", intro: "กดเพื่อชำระเงิน ตอนนี้ระบบรองรับแค่ promt pay" },
                { element: "#pay-qr", intro: "สามารถ scan จ่ายได้ทุกธนาคาร" },
                { element: "#last-button", intro: "กดยืนยัน ก็เรียบร้อยสำหรับการจอง" },
                { element: "#nav-history", intro: "คุณสามารถตรวจสอบรายเอียดประวัติการจองของคุณได้ที่เมนู ประวัติการจอง" },

            ]

        });

        intro.onchange((targetElement) => {
            if (targetElement?.id === "step-booking-button") {
                setTimeout(() => targetElement.click(), 300);
            }

            if (targetElement?.id === "filter-button") {
                setTimeout(() => targetElement.click(), 300);
            }
            if (targetElement?.id === "modal-booking-form") {
                setIsHideMock(false)
                setHideConfirm(true)
            }
            if (targetElement?.id === "modal-comfirm") {
                setHideForm(true)
                setHideConfirm(false)
                setStep(2)
            }
            if (targetElement?.id === "pay-qr") {
                setTimeout(() => targetElement.click(), 300);
            }
            if (targetElement?.id === "nav-history") {
                setTimeout(() => targetElement.click(), 300);
            }

        });

        intro.start();
    };



    return (
        <div className="h-full">
            <div className="h-full w-full">
                <h1 className="text-2xl font-bold text-center my-3">จองคิวรายบุคคล | จองคิวแบบกลุ่ม
                    <span className="ml-2 text-xs bg-red-200 p-[3.5px] rounded-md text-red-600">ยังไม่พร้อมใช้งาน</span>
                </h1>
                <div className="grid grid-cols-3">
                    <div className="col-span-2 ">
                        <TotalBarberSchedule barbers={barbers} onClick={handleOpenModal} appointments={mockAppointments} />
                    </div>
                    <div className="max-w-xl">
                        <div className="border rounded-md border-gray-900 p-2 space-y-4">
                            <div className="flex gap-3 mt-3 ml-3">
                                <button
                                    className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1.5 rounded-md flex items-center transition"
                                    onClick={handleStartIntro}
                                >
                                    ตัวช่วยสอนจอง
                                </button>

                                <button
                                    id="step-booking-button"
                                    className="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1.5 rounded-md flex items-center transition"
                                    onClick={() => setShowSteps((prev) => !prev)}
                                >
                                    ขั้นตอนการจอง
                                </button>
                            </div>

                            <div
                                id="step-booking-detail"
                                className={`pt-3 pl-3 border-t text-sm space-y-4 text-gray-700 transition-all duration-300 ${showSteps ? "block" : "hidden"
                                    }`}
                            >
                                <h3 className="text-base font-semibold">🪪 ขั้นตอนการจองคิว</h3>
                                <ol className="list-decimal list-inside space-y-2">
                                    {bookingStepsMock.map((step) => (
                                        <li key={step.step}>
                                            <span className="font-semibold">{step.title}</span>
                                            <p className="text-gray-600">{step.description}</p>
                                        </li>
                                    ))}
                                </ol>
                            </div>
                            <h1 className="text-xl ml-3">บริการยอดฮิต</h1>

                            <div className={` flex flex-col  gap-6 p-3 pt-0`}>

                                {loadingServices && <p>Loading barbers…</p>}
                                {errorServices && <p className="text-red-500">Error loading barbers: {errorServices}</p>}
                                {services.map((service) => {
                                    return (
                                        <div
                                            key={service.id}
                                            className={`bg-gray-100 rounded-lg shadow-lg transition flex flex-row w-full`}
                                        >
                                            <img
                                                src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${service.Img_path}/${service.Img_name}`}
                                                alt={service.name}
                                                className={`object-cover w-52 h-full rounded-l-lg `}
                                            />

                                            <div
                                                className={`flex-1 p-2 flex flex-col justify-between pl-6`}
                                            >
                                                <div>
                                                    <h2 className="text-2xl font-semibold mb-2">{service.name}</h2>
                                                    <p className="text-gray-400 text-sm mb-4">{service.description}</p>
                                                </div>

                                                <div className="flex items-center justify-between text-gray-900">
                                                    <span className="font-bold text-lg">฿{service.price}</span>
                                                    <span className="text-sm bg-gray-400 px-2 py-1 rounded">
                                                        {service.duration} นาที
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                    );
                                })}


                            </div>
                        </div>
                    </div>
                </div>
                <button id="booking-mock" className="invisible">mock</button>
                <button id="confirm-mock" className="invisible">mock</button>
                <BookingModalExample isOpen={showModalTour} hide={isHideMock} hideForm={hideForm} hideConfirm={hideConfirm} onClose={handleCloseModalTour} step={step} />

                <Modal
                    isOpen={isModalOpen}
                    onClose={handleClose}
                    title="ยืนยันการจองคิว"
                    blurBackground
                >
                    <div className="text-right text-red-600 font-semibold ">
                        ระบบจะล็อกคิวนี้ไว้ {formatTime(countdown)} นาที
                    </div>

                    <form className="" onSubmit={handleBookingSubmit}>
                        <div
                            id="modal-fix"
                            className={selectedBooking ? "block" : "hidden"}
                        >
                            <div className="mb-3">
                                <label className="block text-sm font-medium">วันที่จอง</label>
                                <input
                                    type="text"
                                    className="input input-bordered w-full bg-gray-200 rounded-md"
                                    value={selectedBooking?.date || ""}
                                    readOnly
                                />
                            </div>
                            <div className="mb-3">
                                <label className="block text-sm font-medium">เวลา</label>
                                <input
                                    type="text"
                                    className="input input-bordered w-full bg-gray-200 rounded-md"
                                    value={selectedBooking?.time || ""}
                                    readOnly
                                />
                            </div>
                            <div className="mb-3">
                                <label className="block text-sm font-medium ">ช่าง</label>
                                <input
                                    type="text"
                                    className="input input-bordered w-full bg-gray-200 rounded-md"
                                    value={selectedBooking?.barberName || ""}
                                    readOnly
                                />
                            </div>
                        </div>

                        <div id="modal-form">
                            <div className="mb-3">
                                <label className="block text-sm font-medium">ชื่อลูกค้า</label>
                                <input
                                    type="text"
                                    className="input input-bordered w-full rounded-md"
                                    placeholder="กรุณากรอกชื่อลูกค้า"
                                    data-step="modal-name"
                                />
                            </div>

                            <div className="mb-3">
                                <label className="block text-sm font-medium">เบอร์โทรศัพท์</label>
                                <input
                                    type="text"
                                    className="input input-bordered w-full rounded-md"
                                    placeholder="กรุณากรอกเบอร์โทร"
                                    data-step="modal-phone"
                                />
                            </div>
                        </div>
                        <h3 className="mb-1">เลือกบริการ</h3>
                        <div className="grid grid-cols-4 w-full gap-3 ">
                            {loadingServices && <p>Loading services…</p>}
                            {errorServices && (
                                <p className="text-red-500">Error loading services: {errorServices}</p>
                            )}

                            {services.map((service) => {
                                const isSelected = service.id === selectedServiceId;
                                if (selectedServiceId && !isSelected) return null;

                                return (
                                    <div
                                        key={service.id}
                                        className="bg-gray-100 rounded-xl shadow-lg flex flex-col p-1"
                                    >
                                        <div className="flex items-center p-1">
                                            {/* รูปภาพวงกลมเล็ก */}
                                            {/* <div className="flex-shrink-0">
                                                <img
                                                    src={`https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/${service.Img_path}/${service.Img_name}`}
                                                    alt={service.name}
                                                    className="w-8 h-8 rounded-full object-cover"
                                                />
                                            </div> */}

                                            {/* ข้อความ */}
                                            <div className="flex flex-col justify-center flex-grow overflow-hidden">
                                                <div className="text-lg font-semibold truncate">{service.name}</div>
                                                <div className="text-gray-400 text-sm truncate">{service.description}</div>
                                            </div>

                                            {/* ราคา + นาที */}
                                            <div className="bg-gray-400 text-white text-right rounded-md px-1 py-1 flex flex-col items-end justify-center min-w-[60px]">
                                                <span className="text-md font-bold">฿{service.price}</span>
                                                <span className="text-sm">{service.duration} นาที</span>
                                            </div>
                                        </div>

                                        {/* ปุ่ม เลือก / ยกเลิก */}
                                        {!isSelected ? (
                                            <button
                                                type="button"
                                                className="w-full bg-green-400 hover:bg-green-700 text-white  rounded"
                                                onClick={() => setSelectedServiceId(service.id)}
                                            >
                                                เลือก
                                            </button>
                                        ) : (
                                            <button
                                                type="button"
                                                className="w-full bg-gray-300 text-gray-700  rounded"
                                                onClick={() => setSelectedServiceId(null)}
                                            >
                                                ยกเลิก
                                            </button>
                                        )}

                                    </div>
                                );
                            })}
                        </div>


                        <div className="flex gap-3 pt-4">
                            <button
                                className="w-full bg-gray-600 hover:bg-green-700 text-white py-1 rounded"
                                onClick={handleClose}

                            >
                                ยกเลิก
                            </button>
                            <button
                                type="submit"
                                className="w-full bg-green-600 hover:bg-green-700 text-white py-1 rounded"
                                id="comfirm-booking"
                            >
                                ยืนยันการจอง
                            </button>

                        </div>
                    </form>
                </Modal>



                {/* Toast แสดงผลลัพธ์ */}
                {toastInfo && (
                    <div className="fixed bottom-4 right-4 z-50">
                        <Toast
                            message={toastInfo.message}
                            variant={toastInfo.variant}
                            onClose={() => setToastInfo(null)}
                            position="top-right"
                            duration={6000}
                        />
                    </div>
                )}

            </div>

        </div >
    );
}
interface BookingModalProps {
    isOpen: boolean;
    hide: boolean
    onClose: () => void;
    step: Number
    hideForm: boolean
    hideConfirm: boolean
}

const BookingModalExample = ({ isOpen, hide, onClose, step, hideForm, hideConfirm }: BookingModalProps) => {
    return (
        <div className={`${hide ? "invisible" : ""}`}>
            <Modal
                isOpen={isOpen}
                onClose={onClose}
                title={step === 1 ? "กรอกรายละเอียดจองคิว" : "ชำระเงินและตรวจสอบข้อมูล"}
                blurBackground
                modalName={"booking-form"}
            >
                <div>
                    <BookingFormExample isOpen={isOpen} hide={hideForm} />
                    <ConfirmExample hide={hideConfirm} />
                </div>

            </Modal>
        </div>
    );
};

interface BookingFormProps {
    isOpen: boolean;
    hide: boolean;
}

const BookingFormExample = ({ isOpen, hide }: BookingFormProps) => {
    const [countdown, setCountdown] = useState(7 * 60);

    const [selectedBooking] = useState({
        date: "2025-07-15",
        time: "10:30",
        barberName: "ช่างแจ็ค",
        cusName: "ขนมต้ม",
        cusPhone: "012 345 6789"
    });

    const handleBookingSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        alert("Mock: บันทึกข้อมูลการจองเรียบร้อยแล้ว");
    };

    useEffect(() => {

        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    clearInterval(timer);
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);

        return () => clearInterval(timer);
    }, [isOpen]);

    const formatTime = (seconds: number) => {
        const m = Math.floor(seconds / 60);
        const s = seconds % 60;
        return `${m}:${s.toString().padStart(2, "0")}`;
    };

    return (
        <div className={`${hide ? "hidden" : ""}`}>
            <div className="text-right text-red-600 font-semibold mb-3 " >
                <p id="modal-time" className="inline">ระบบจะล็อกคิวนี้ไว้ {formatTime(countdown)} นาที</p>
            </div>

            <form className="space-y-5" onSubmit={handleBookingSubmit}>
                <div id="modal-fix" className={selectedBooking ? "block" : "hidden"}>
                    <div className="mb-3">
                        <label className="block text-sm font-medium">วันที่จอง</label>
                        <input
                            type="text"
                            className="input input-bordered w-full bg-gray-200 rounded-md"
                            value={selectedBooking.date}
                            readOnly
                        />
                    </div>
                    <div className="mb-3">
                        <label className="block text-sm font-medium">เวลา</label>
                        <input
                            type="text"
                            className="input input-bordered w-full bg-gray-200 rounded-md"
                            value={selectedBooking.time}
                            readOnly
                        />
                    </div>
                    <div className="mb-3">
                        <label className="block text-sm font-medium ">ช่าง</label>
                        <input
                            type="text"
                            className="input input-bordered w-full bg-gray-200 rounded-md"
                            value={selectedBooking.barberName}
                            readOnly
                        />
                    </div>
                </div>

                <div id="modal-form">
                    <div className="mb-3">
                        <label className="block text-sm font-medium">ชื่อลูกค้า</label>
                        <input
                            type="text"
                            className="input input-bordered w-full rounded-md"
                            placeholder="กรุณากรอกชื่อลูกค้า"
                            data-step="modal-name"
                            value={selectedBooking.cusName}
                        />
                    </div>

                    <div className="mb-3">
                        <label className="block text-sm font-medium">เบอร์โทรศัพท์</label>
                        <input
                            type="text"
                            className="input input-bordered w-full rounded-md"
                            placeholder="กรุณากรอกเบอร์โทร"
                            data-step="modal-phone"
                            value={selectedBooking.cusPhone}
                        />
                    </div>
                </div>

                <div className="pt-4">
                    <button
                        type="submit"
                        className="w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded"
                        id="comfirm-booking"
                    >
                        ดำเนินการต่อ
                    </button>
                </div>
            </form>
        </div>
    )

}
interface ConfirmModalProps {
    hide: boolean
}
const ConfirmExample = ({ hide }: ConfirmModalProps) => {
    return (
        <div className={`${hide ? "hidden" : ""} w-full flex`}>
            <div className="w-1/2">
                <PaymentOptions />
            </div>
            <div className="w-1/2 flex flex-col gap-3" id="modal-comfirm">
                <div className="flex flex-col items-center">

                    <h2 className="text-2xl font-semibold mt-2">ใบจอง</h2>
                </div>

                <div className="flex justify-between text-base px-2">
                    <p>วันที่จอง: <span className="font-medium">2025-07-15</span></p>
                    <p>เวลาที่จอง: <span className="font-medium">10:30</span></p>
                </div>

                <hr className="border-t border-dashed border-gray-900 " />
                <div className="grid grid-cols-2 gap-y-1 text-base px-2 mt-0">
                    <p>ชื่อลูกค้า:</p>
                    <p className="text-right">คุณ ขนมต้ม</p>

                    <p>สมาชิก:</p>
                    <p className="text-right">ไม่ได้เป็น</p>

                    <p>ช่างที่เลือก:</p>
                    <p className="text-right">ช่างแจ็ค</p>

                    <p>บริการที่เลือก:</p>
                    <p className="text-right">ตัดผมชาย</p>

                    <p>ราคา:</p>
                    <p className="text-right">440 บาท</p>

                    <p>เวลาโดยประมาณ:</p>
                    <p className="text-right">50 นาที</p>
                </div>
                <hr className="border-t border-gray-900 px-2 border-dashed border-4" />
                <div className="mt-0 px-2">
                    <h5 className="font-semibold">ข้อความฝากถึงช่าง:</h5>
                    <p className="text-gray-700 mt-1">
                        -
                    </p>
                </div>
                <hr className="border-t border-dashed border-gray-900 " />
                <div className="text-center">
                    <p>ขอบคุณที่ใช้บริการ</p>
                </div>
                <div className="pt-4">
                    <button
                        type="submit"
                        className="w-full bg-green-600 hover:bg-green-700 text-white py-2 rounded"
                        id="last-button"
                    >
                        ยืนยันการจอง
                    </button>
                </div>
            </div>


        </div>
    )
}





{/* <button onClick={handleOpen} className="bg-blue-600 text-white px-4 py-2 rounded">
                    เปิด Modal
                </button> */}


{/* <Modal isOpen={isModalOpen} onClose={handleClose} title={isRegisterForm ? "ลงทะเบียน" : "เข้าสู่ระบบ"} blurBackground>
                    {isRegisterForm ? (
                        <form className="space-y-5">
                            <div>
                                <label className="flex items-end text-sm font-medium text-gray-700 mb-1">
                                    ชื่อลูกค้า
                                    <CustomTooltip
                                        id="tooltip-cusname"
                                        content="แนะนำเป็นภาษาไทย"
                                        trigger="hover"
                                        placement="top"
                                        bgColor="bg-gray-200"
                                        textColor="text-gray-900"
                                        textSize="text-sm"
                                        className="ml-1"
                                    >
                                        <span><FiInfo /></span>
                                    </CustomTooltip>
                                </label>
                                <input
                                    type="text"
                                    placeholder="กรุณากรอกชื่อ"
                                    className="input input-bordered w-full"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">เบอร์โทร</label>
                                <input
                                    type="text"
                                    placeholder="กรุณากรอกเบอร์โทรศัพท์"
                                    className="input input-bordered w-full"
                                />
                            </div>

                            <div>
                                <label className="flex items-end text-sm font-medium text-gray-700 mb-1 ">
                                    รหัสผ่าน
                                    <CustomTooltip
                                        id="tooltip-password"
                                        content="รหัสผ่านสำหรับเข้าใช้งานครั้งถัดไป"
                                        trigger="hover"
                                        placement="top"
                                        bgColor="bg-gray-200"
                                        textColor="text-gray-900"
                                        textSize="text-sm"
                                        className="ml-1"
                                    >
                                        <span><FiInfo /></span>
                                    </CustomTooltip>
                                </label>
                                <input
                                    type="password"
                                    placeholder="ตั้งรหัสผ่าน"
                                    className="input input-bordered w-full"
                                />
                            </div>

                            <div className="flex gap-4 pt-4">
                                <button type="submit" className="w-1/2 bg-green-600 hover:bg-green-700 text-white py-2 rounded">
                                    ลงทะเบียน
                                </button>
                                <button
                                    type="button"
                                    onClick={() => setIsRegisterForm(false)} // สลับไป login
                                    className="w-1/2 bg-gray-400 hover:bg-gray-500 text-white py-2 rounded"
                                >
                                    มีบัญชีอยู่แล้ว?
                                </button>
                            </div>
                        </form>
                    ) : (
                        <form className="space-y-5">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">เบอร์โทร</label>
                                <input
                                    type="text"
                                    placeholder="กรุณากรอกเบอร์โทรศัพท์"
                                    className="input input-bordered w-full"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">รหัสผ่าน</label>
                                <input
                                    type="password"
                                    placeholder="กรุณากรอกรหัสผ่าน"
                                    className="input input-bordered w-full"
                                />
                            </div>

                            <div className="flex gap-4 pt-4">
                                <button type="submit" className="w-1/2 bg-blue-600 hover:bg-blue-700 text-white py-2 rounded">
                                    เข้าสู่ระบบ
                                </button>
                                <button
                                    type="button"
                                    onClick={() => setIsRegisterForm(true)} // สลับกลับไป register
                                    className="w-1/2 bg-gray-400 hover:bg-gray-500 text-white py-2 rounded"
                                >
                                    สมัครสมาชิก
                                </button>
                            </div>
                        </form>
                    )}
                </Modal> */}


function PaymentOptions() {
    const [paymentMethod, setPaymentMethod] = useState("");

    const handlePaymentClick = (method: string) => {
        setPaymentMethod(method);
    };

    return (
        <div>
            <div className="flex justify-start gap-4">
                <button
                    className="bg-blue-600 text-white px-2 py-2 rounded"
                    onClick={() => handlePaymentClick("qrcode")}
                    id="pay-button"
                >
                    จ่ายผ่าน prompt pay
                </button>
            </div>

            <div className="flex justify-center mt-4">
                <img src="/QR.jpeg" alt="qrcode" className={`${paymentMethod === "qrcode" ? "invisible" : ""} w-72`} id="pay-qr" />
            </div>
        </div>
    );
}
