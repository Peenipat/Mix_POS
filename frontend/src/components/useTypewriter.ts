import { useEffect, useState } from "react";

export function useTypewriter(messages: string[], typingSpeed = 70, deletingSpeed = 40, pauseTime = 1500) {
  const [displayText, setDisplayText] = useState("");
  const [messageIndex, setMessageIndex] = useState(0);
  const [charIndex, setCharIndex] = useState(0);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const currentMessage = messages[messageIndex];
    const delay = deleting ? deletingSpeed : typingSpeed;

    const timeout = setTimeout(() => {
      if (!deleting) {
        setDisplayText(currentMessage.slice(0, charIndex + 1));
        setCharIndex(charIndex + 1);

        if (charIndex + 1 === currentMessage.length) {
          setTimeout(() => setDeleting(true), pauseTime);
        }
      } else {
        setDisplayText(currentMessage.slice(0, charIndex - 1));
        setCharIndex(charIndex - 1);

        if (charIndex === 0) {
          setDeleting(false);
          setMessageIndex((messageIndex + 1) % messages.length); // วนลูปข้อความ
        }
      }
    }, delay);

    return () => clearTimeout(timeout);
  }, [charIndex, deleting, messageIndex, messages, typingSpeed, deletingSpeed, pauseTime]);

  return displayText;
}
