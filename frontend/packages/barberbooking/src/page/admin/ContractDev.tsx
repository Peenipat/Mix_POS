import { FaFacebookF } from "react-icons/fa";
import { HiOutlineMail } from "react-icons/hi";

interface NotReadyProps {
    message?: string;
}

export default function ContractDev({ message }: NotReadyProps) {
    return (
        <div className="flex flex-col items-center justify-center text-center min-h-[50vh] space-y-4">
            <img src={"https://test-img-upload-xs-peenipat.s3.ap-southeast-1.amazonaws.com/avatars/Me.png"} className="w-86 h-86"/>
            <h1 className="text-2xl font-bold">นิพัทธ์ (พี)</h1>

            <div className="flex items-center space-x-2 text-gray-700 text-base">
                <FaFacebookF className="w-5 h-5 text-blue-600" />
                <span>: Nipat Chapakdee</span>
            </div>

            <div className="flex items-center space-x-2 text-gray-700 text-base">
                <HiOutlineMail className="w-5 h-5 text-red-500" />
                <span>: nipatchapakdee@gmail.com</span>
            </div>

            <p className="text-sm text-gray-500">แนะนำติดต่อผ่าน Facebook จะเร็วที่สุด</p>
        </div>
    );
}
