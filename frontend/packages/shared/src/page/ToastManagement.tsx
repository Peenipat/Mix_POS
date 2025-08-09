import { Toast } from "../components/Toast";
import { useState } from "react";
export default function ToastManagement() {

    const [message, setMessage] = useState("ตัวอย่างแจ้งเตือน");
    const [variant, setVariant] = useState<"success" | "error" | "warning">("success");
    const [showIcon, setShowIcon] = useState(true);
    const [useFixed, setUseFixed] = useState(false);
    const [disableClose, setDisableClose] = useState(true);
    const [duration, setDuration] = useState(0);
    const [position, setPosition] = useState<
        "top-left" | "top-right" | "bottom-left" | "bottom-right"
    >("bottom-right");

    return (
        <div className="min-h-screen  items-center justify-center">
            <h1 className="text-xl font-bold my-3 text-gray-700">ตัวอย่างแบบ Toast ที่ใช้งานอยู่ </h1>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-3 w-full max-w-6xl">
                {/* Success Preview */}
                <div className="border border-gray-300 rounded-lg  bg-white shadow-md flex items-center justify-center min-h-[150px]">
                    <Toast
                        message={"ข้อความแจ้งเตือน"}
                        variant="success"
                        useFixed={false}
                        showIcon
                        duration={0} // ค้างไว้
                        disableClose={true}
                    />
                </div>

                {/* Error Preview */}
                <div className="border border-gray-300 rounded-lg p-6 bg-white shadow-md flex items-center justify-center min-h-[150px]">
                    <Toast
                        message={"ข้อความแจ้งเตือน"}
                        variant="error"
                        useFixed={false}
                        showIcon
                        duration={0} 
                        disableClose={true}
                    />
                </div>

                {/* Warning Preview */}
                <div className="border border-gray-300 rounded-lg p-6 bg-white shadow-md flex items-center justify-center min-h-[150px]">
                    <Toast
                        message={"ข้อความแจ้งเตือน"}
                        variant="warning"
                        useFixed={false}
                        showIcon
                        duration={0}
                        disableClose={true}
                    />
                </div>
            </div>

            <h1 className="text-xl my-3 font-bold text-gray-700">ส่วนที่แก้ไขได้</h1>

            {/* Form */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl">
                <div className="flex flex-col gap-2">
                    <label className="font-medium">ข้อความแจ้งเตือน</label>
                    <input
                        className="border rounded px-3 py-2"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                    />

                    <label className="font-medium mt-4">ประเภท</label>
                    <select
                        className="border rounded px-3 py-2"
                        value={variant}
                        onChange={(e) => setVariant(e.target.value as any)}
                    >
                        <option value="success">Success</option>
                        <option value="error">Error</option>
                        <option value="warning">Warning</option>
                    </select>

                    <label className="font-medium mt-4">ตำแหน่ง</label>
                    <select
                        className="border rounded px-3 py-2"
                        value={position}
                        onChange={(e) => setPosition(e.target.value as any)}
                    >
                        <option value="top-left">Top Left</option>
                        <option value="top-right">Top Right</option>
                        <option value="bottom-left">Bottom Left</option>
                        <option value="bottom-right">Bottom Right</option>
                    </select>
                </div>

                <div className="flex flex-col gap-2">
                    <label className="font-medium">⏱ Duration (ms)</label>
                    <input
                        type="number"
                        className="border rounded px-3 py-2"
                        value={duration}
                        onChange={(e) => setDuration(Number(e.target.value))}
                    />

                    <label className="font-medium flex items-center gap-2 mt-4">
                        <input
                            type="checkbox"
                            checked={showIcon}
                            onChange={() => setShowIcon(!showIcon)}
                        />
                        แสดง Icon
                    </label>

                    <label className="font-medium flex items-center gap-2">
                        <input
                            type="checkbox"
                            checked={disableClose}
                            onChange={() => setDisableClose(!disableClose)}
                        />
                        ปุ่ม X กดไม่ได้
                    </label>
                </div>
            </div>

            {/* 🔲 Toast Preview */}
            <div className="rounded-lg p-10 min-h-[250px] w-full max-w-2xl mx-auto bg-white">
                <Toast
                    message={message}
                    variant={variant}
                    position={position}
                    showIcon={showIcon}
                    useFixed={false}         
                    containerMode={true}     
                    disableClose={disableClose}
                    closable={true}
                    duration={duration}
                />
            </div>
        </div>
    );
}
