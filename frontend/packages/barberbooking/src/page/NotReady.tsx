
interface NotReadyProps {
  message?: string;
}

export default function NotReady({ message }: NotReadyProps) {
  return (
    <div className="flex flex-col items-center justify-center text-center min-h-[50vh]">
      <div className="text-[200px] mb-4">🚧</div>
      <h1 className="text-2xl font-bold mb-2">หน้านี้ยังไม่พร้อมใช้งาน</h1>
      {message && <p className="text-gray-600 text-base">{message}</p>}
    </div>
  );
}
