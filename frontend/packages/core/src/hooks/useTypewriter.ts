import { useEffect, useState } from "react";

/**
 * Custom hook สำหรับพิมพ์ข้อความแบบ typewriter effect
 * @param messages – อาร์เรย์ของข้อความที่จะวนพิมพ์
 * @param typingSpeed – ความเร็วในการพิมพ์ 
 * @param deletingSpeed – ความเร็วในการลบ 
 * @param pauseTime – เวลาหยุดพักก่อนเริ่มลบ 
 * @returns displayText – ข้อความที่ควรแสดงในปัจจุบัน
 */
export function useTypewriter(messages: string[], typingSpeed = 70, deletingSpeed = 40, pauseTime = 1500) {
  const [displayText, setDisplayText] = useState(""); //ข้อความที่แสดง
  const [messageIndex, setMessageIndex] = useState(0); //กำลังพิมพ์ข้อความไหน
  const [charIndex, setCharIndex] = useState(0); // กำลังพิมพ์ หรือ ลบ ที่ละตัว
  const [deleting, setDeleting] = useState(false); // state กำลังลบ หรือ state กำลังพิมพ์

  useEffect(() => {
    //เลือกตาม index และใส่ delay
    const currentMessage = messages[messageIndex];
    const delay = deleting ? deletingSpeed : typingSpeed;


    const timeout = setTimeout(() => {
      if (!deleting) {
        //กำลังพิมพ์
        setDisplayText(currentMessage.slice(0, charIndex + 1));
        setCharIndex(charIndex + 1); 

        // stop ก่อนจะลบ
        if (charIndex + 1 === currentMessage.length) {
          setTimeout(() => setDeleting(true), pauseTime);
        }
      } else {
        //กำลังลบ
        setDisplayText(currentMessage.slice(0, charIndex - 1));
        setCharIndex(charIndex - 1);

        if (charIndex === 0) {
          setDeleting(false); // ลบเสร็จ
          setMessageIndex((messageIndex + 1) % messages.length); // วนลูปข้อความ (กลับไปพิมพ์ใหม่)
        } 
      }
    }, delay);

    return () => clearTimeout(timeout);
  }, [charIndex, deleting, messageIndex, messages, typingSpeed, deletingSpeed, pauseTime]);

  return displayText;
}
