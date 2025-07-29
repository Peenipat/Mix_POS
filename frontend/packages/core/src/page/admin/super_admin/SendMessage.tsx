import axios from "axios";
import { useState } from "react";

// Card Component
const Card = ({ children }: { children: React.ReactNode }) => (
  <div className="rounded-2xl shadow-md border bg-white p-4 w-full max-w-xl mx-auto my-6">
    {children}
  </div>
);

const CardContent = ({ children }: { children: React.ReactNode }) => (
  <div className="p-2 space-y-4">{children}</div>
);

// Input Component
const Input = (props: React.InputHTMLAttributes<HTMLInputElement>) => (
  <input
    {...props}
    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300"
  />
);

// Textarea Component
const Textarea = (props: React.TextareaHTMLAttributes<HTMLTextAreaElement>) => (
  <textarea
    {...props}
    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300"
  />
);

// Button Component
const Button = ({ children, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) => (
  <button
    {...props}
    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
  >
    {children}
  </button>
);

// Label Component
const Label = ({ children, htmlFor }: { children: React.ReactNode; htmlFor?: string }) => (
  <label htmlFor={htmlFor} className="block text-sm font-medium text-gray-700">
    {children}
  </label>
);

// Checkbox Component
const Checkbox = ({ checked, onChange }: { checked: boolean; onChange: () => void }) => (
  <input
    type="checkbox"
    checked={checked}
    onChange={onChange}
    className="h-4 w-4 text-blue-600 border-gray-300 rounded"
  />
);

// Select Component
const Select = ({ value, onChange, options }: { value: string; onChange: (val: string) => void; options: string[] }) => (
  <select
    value={value}
    onChange={(e) => onChange(e.target.value)}
    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring focus:border-blue-300"
  >
    <option value="">-- ทั้งหมด --</option>
    {options.map((opt) => (
      <option key={opt} value={opt}>{opt}</option>
    ))}
  </select>
);

// Main UI
const SendMessage = () => {
  const [users] = useState([
    { id: "1", name: "Admin A", telegramId: "123456789", branch: "สาขา A", tenant: "Tenant 1" },
    { id: "2", name: "Staff B", telegramId: "987654321", branch: "สาขา A", tenant: "Tenant 1" },
    { id: "3", name: "Manager C", telegramId: "1122334455", branch: "สาขา B", tenant: "Tenant 2" },
  ]);

  const [selectedIds, setSelectedIds] = useState<string[]>([]);
  const [subject, setSubject] = useState("");
  const [message, setMessage] = useState("");
  const [branchFilter, setBranchFilter] = useState("");
  const [tenantFilter, setTenantFilter] = useState("");

  const toggleSelection = (id: string) => {
    setSelectedIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  };

  const filteredUsers = users.filter((u) => {
    return (
      (!branchFilter || u.branch === branchFilter) &&
      (!tenantFilter || u.tenant === tenantFilter)
    );
  });

  const groupedUsers = filteredUsers.reduce((acc, user) => {
    const key = `${user.tenant} / ${user.branch}`;
    if (!acc[key]) acc[key] = [];
    acc[key].push(user);
    return acc;
  }, {} as Record<string, typeof users>);

  const branchOptions = Array.from(new Set(users.map((u) => u.branch)));
  const tenantOptions = Array.from(new Set(users.map((u) => u.tenant)));
  const selectedUsers = users.filter((u) => selectedIds.includes(u.id));
  const mockUser = {name:"nipat",telegramId:8053367943}

  const sendTelegramMessages = async (
    users: { name: string; telegramId: number }[],
    subject: string,
    message: string
  ) => {
    const content = `หัวข้อ: ${subject || "-"}\n\n${message}`;
  
    try {
      await Promise.all(
        users.map((user) =>
          axios.post(" https://a5d1fc82d74a.ngrok-free.app/api/v1/core/telegram/send", {
            chat_id: user.telegramId,
            message: content,
          })
        )
      );
      alert("✅ ส่งข้อความเรียบร้อยแล้ว");
    } catch (err) {
      console.error("❌ ส่งไม่สำเร็จ:", err);
      alert("เกิดข้อผิดพลาดในการส่งข้อความ");
    }
  };
  

  return (
    <Card>
      <CardContent>
        <h2 className="text-xl font-semibold">ส่งข้อความถึงผู้ใช้งานหลายคนผ่าน Telegram</h2>

        <div className="space-y-2">
          <Label>กรองตาม Tenant</Label>
          <Select value={tenantFilter} onChange={setTenantFilter} options={tenantOptions} />
        </div>

        <div className="space-y-2">
          <Label>กรองตามสาขา</Label>
          <Select value={branchFilter} onChange={setBranchFilter} options={branchOptions} />
        </div>

        <div className="space-y-2">
          <Label>เลือกผู้ใช้</Label>
          <div className="space-y-3 border p-3 rounded-md max-h-60 overflow-y-auto">
            {Object.entries(groupedUsers).map(([group, members]) => (
              <div key={group}>
                <div className="font-semibold text-gray-600 mb-1">{group}</div>
                {members.map((user) => (
                  <label key={user.id} className="flex items-center space-x-2 ml-4">
                    <Checkbox
                      checked={selectedIds.includes(user.id)}
                      onChange={() => toggleSelection(user.id)}
                    />
                    <span>{user.name}</span>
                  </label>
                ))}
              </div>
            ))}
          </div>

          <p className="text-sm text-gray-600 pt-1">
            เลือกแล้ว {selectedIds.length} คน
          </p>

          {selectedUsers.length > 0 && (
            <ul className="text-sm text-blue-800 list-disc ml-5">
              {selectedUsers.map((u) => (
                <li key={u.id}>{u.name}</li>
              ))}
            </ul>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="subject">หัวข้อ</Label>
          <Input
            id="subject"
            value={subject}
            onChange={(e) => setSubject(e.target.value)}
            placeholder="เช่น แจ้งเตือนยอดขาย"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="message">ข้อความ</Label>
          <Textarea
            id="message"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            placeholder="พิมพ์ข้อความที่ต้องการส่ง"
            rows={4}
          />
        </div>

        <div className="pt-2">
          <Button
             onClick={() => sendTelegramMessages([mockUser], subject, message)}
             disabled={selectedIds.length === 0 || !message}
          >
            ส่งข้อความผ่าน Telegram
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};

export default SendMessage;
