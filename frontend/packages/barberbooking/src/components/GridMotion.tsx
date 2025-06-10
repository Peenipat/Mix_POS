import React, { useEffect, useRef, FC } from "react";
import { gsap } from "gsap";
import "./GridMotion.css";

interface GridMotionProps {
  items?: React.ReactNode[];
  gradientColor?: string;
}

const GridMotion: FC<GridMotionProps> = ({
  items = [],
  gradientColor = "black",
}) => {
  const gridRef = useRef<HTMLDivElement>(null);
  const rowRefs = useRef<(HTMLDivElement | null)[]>([]);
  const mouseXRef = useRef<number>(window.innerWidth / 2);

  // Ensure the grid has 28 items
  const totalItems = 28;
  const defaultItems = Array.from(
    { length: totalItems },
    (_, index) => `Item ${index + 1}`
  );

  const combinedItems: React.ReactNode[] =
    items.length > 0 ? items.slice(0, totalItems) : defaultItems;

  useEffect(() => {
    // Disable lag smoothing for consistent ticker
    gsap.ticker.lagSmoothing(0);

    // 1) Mousemove listener to update mouseXRef
    const handleMouseMove = (e: MouseEvent) => {
      mouseXRef.current = e.clientX;
    };
    window.addEventListener("mousemove", handleMouseMove);

    // 2) Parallax loop driven by GSAP ticker
    const parallaxLoop = gsap.ticker.add(() => {
      const maxMove = 300;
      const baseDur = 0.8;
      const inertia = [0.6, 0.4, 0.3, 0.2];

      rowRefs.current.forEach((row, idx) => {
        if (!row) return;
        const dir = idx % 2 === 0 ? 1 : -1;
        const targetX =
          ((mouseXRef.current / window.innerWidth) * maxMove -
            maxMove / 2) *
          dir;

        gsap.to(row, {
          x: targetX,
          duration: baseDur + inertia[idx % inertia.length],
          ease: "power3.out",
          overwrite: "auto",
        });
      });
    });

    // 3) Continuous drift tween on each row
    rowRefs.current.forEach((row, idx) => {
      if (!row) return;
      const drift = 100 * (idx % 2 === 0 ? 1 : -1);
      gsap.to(row, {
        x: `+=${drift}`,
        duration: 6,
        ease: "sine.inOut",
        repeat: -1,
        yoyo: true,
        overwrite: "auto",
      });
    });

    // Cleanup on unmount
    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
      parallaxLoop(); // remove ticker callback
      rowRefs.current.forEach((row) => {
        if (row) gsap.killTweensOf(row);
      });
    };
  }, []);

  return (
    <div className="noscroll loading" ref={gridRef}>
      <section
        className="intro"
        style={{
          background: `radial-gradient(circle, ${gradientColor} 0%, transparent 100%)`,
        }}
      >
        <div className="gridMotion-container">
          {Array.from({ length: 4 }, (_, rowIndex) => (
            <div
              key={rowIndex}
              className="row"
              ref={(el) => (rowRefs.current[rowIndex] = el)}
            >
              {Array.from({ length: 7 }, (_, itemIndex) => {
                const content = combinedItems[rowIndex * 7 + itemIndex];
                return (
                  <div key={itemIndex} className="row__item">
                    <div className="row__item-inner" style={{ backgroundColor: "#111" }}>
                      {typeof content === "string" && content.startsWith("http") ? (
                        <div
                          className="row__item-img"
                          style={{ backgroundImage: `url(${content})` }}
                        />
                      ) : (
                        <div className="row__item-content">{content}</div>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          ))}
        </div>
        <div className="fullview"></div>
      </section>
    </div>
  );
};

export default GridMotion;
