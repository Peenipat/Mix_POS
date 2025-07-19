import { StructureLayout } from "@object/shared/components/StructureLayout";
import { useState } from "react";

// Mock components ที่อาจใช้ใน Layout ต่าง ๆ
export const Logo = () => (
  <div className="text-2xl font-bold text-blue-600 p-4">Barber Shop</div>
);

export const MainNav = () => (
  <nav className="flex space-x-6 bg-gray-100 p-4">
    <a href="#" className="hover:underline">หน้าหลัก</a>
    <a href="#" className="hover:underline">บริการ</a>
    <a href="#" className="hover:underline">ทีมช่าง</a>
    <a href="#" className="hover:underline">ประวัติการจอง</a>
  </nav>
);

export const Banner = () => (
  <div className="bg-gradient-to-r from-rose-400 to-pink-500 text-white text-center py-6">
    <h1 className="text-3xl font-semibold">จองคิวตัดผมกับช่างมืออาชีพ</h1>
  </div>
);

export const Sidebar = () => (
  <div className="bg-gray-200 h-full p-4 space-y-2">
    <p className="font-medium">ประเภทบริการ</p>
    <ul className="space-y-1">
      <li><a href="#" className="hover:underline">ตัดผม</a></li>
      <li><a href="#" className="hover:underline">สระผม</a></li>
      <li><a href="#" className="hover:underline">กันหนวด</a></li>
    </ul>
  </div>
);

export const MainContent = () => (
  <div className="p-6 space-y-4">
    <h2 className="text-xl font-semibold">บริการยอดนิยม</h2>
    <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
      {["Haircut", "Shampoo", "Beard Trim"].map((item) => (
        <div key={item} className="bg-white shadow p-4 rounded">
          <h3 className="font-medium">{item}</h3>
          <p>ราคา: ฿{item === "Haircut" ? 200 : item === "Shampoo" ? 100 : 150}</p>
        </div>
      ))}
    </div>
  </div>
);

export const Footer = () => (
  <footer className="bg-gray-800 text-white text-center p-4">
    &copy; 2025 Barber Shop. All rights reserved.
  </footer>
);

export const BackgroundDecor = () => (
  <div className="absolute inset-0 bg-gradient-to-tr from-white to-gray-100 -z-10"></div>
);

// Layout options for testing
const layoutOptions = [
  "sidebar-left",
  "sidebar-right",
  "banner-3box",
  "grid-masonry",
  "f-shape",
  "z-shape",
  "card-block",
  "magazine",
  "featured",
  "split-screen",
  "interactive",
  "two-column",
] as const;

export default function TestLayout() {
  const [selectedLayout, setSelectedLayout] = useState<(typeof layoutOptions)[number]>("sidebar-left");

  return (
    <div className="relative min-h-screen">
      {/* Layout selector */}
      <div className="fixed top-0 left-0 z-50 bg-white shadow px-4 py-2 w-full flex items-center space-x-2">
        <label className="font-semibold">เลือก Layout: </label>
        <select
          className="border p-1 rounded"
          value={selectedLayout}
          onChange={(e) => setSelectedLayout(e.target.value as any)}
        >
          {layoutOptions.map((layout) => (
            <option key={layout} value={layout}>
              {layout}
            </option>
          ))}
        </select>
      </div>

      <div className="pt-16">
        <StructureLayout
          layout={selectedLayout}
          logo={<Logo />}
          nav={<MainNav />}
          header={["sidebar-left", "banner-3box", "magazine", "z-shape"].includes(selectedLayout) && <Banner />}
          sidebar={["sidebar-left", "sidebar-right", "featured", "f-shape", "split-screen", "magazine"].includes(selectedLayout) && <Sidebar />}
          body={<MainContent />}
          footer={!["interactive"].includes(selectedLayout) && <Footer />}
          background={["interactive", "banner-3box", "card-block"].includes(selectedLayout) && <BackgroundDecor />}
        />
      </div>
    </div>
  );
}
