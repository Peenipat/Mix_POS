import { FC, ReactNode } from "react";

type LayoutType =
  | "sidebar-left"
  | "sidebar-right"
  | "banner-3box"
  | "grid-masonry"
  | "f-shape"
  | "z-shape"
  | "card-block"
  | "magazine"
  | "featured"
  | "split-screen"
  | "interactive"
  | "two-column";


interface StructureLayoutProps {
  layout: LayoutType;
  logo: ReactNode;
  nav: ReactNode;
  header?: ReactNode;
  sidebar?: ReactNode;
  body: ReactNode;
  footer?: ReactNode;
  background?: ReactNode;
}
export const StructureLayout: FC<StructureLayoutProps> = ({
  layout,
  logo,
  nav,
  header,
  sidebar,
  body,
  footer,
  background,
}) => {
  const HeaderBlock = (
    <>
      {logo && <div className="p-2">{logo}</div>}
      {nav && <nav>{nav}</nav>}
      {header && <div>{header}</div>}
    </>
  );

  const FooterBlock = footer && <footer>{footer}</footer>;

  switch (layout) {
    case "sidebar-left":
      return (
        <div className="min-h-screen flex flex-col">
          {background}
          {HeaderBlock}
          <main className="flex flex-1">
            {sidebar && <aside className="w-1/4">{sidebar}</aside>}
            <section className="flex-1">{body}</section>
          </main>
          {FooterBlock}
        </div>
      );

    case "sidebar-right":
      return (
        <div className="min-h-screen flex flex-col">
          {HeaderBlock}
          <main className="flex flex-1">
            <section className="flex-1">{body}</section>
            {sidebar && <aside className="w-1/4">{sidebar}</aside>}
          </main>
          {FooterBlock}
        </div>
      );

    case "banner-3box":
      return (
        <div>
          {HeaderBlock}
          <section className="grid grid-cols-1 md:grid-cols-3 gap-4 p-4">{body}</section>
          {FooterBlock}
        </div>
      );

    case "grid-masonry":
      return (
        <div>
          {HeaderBlock}
          <div className="columns-1 md:columns-2 lg:columns-3 gap-4 p-4">{body}</div>
          {FooterBlock}
        </div>
      );

    case "f-shape":
      return (
        <div className="flex flex-col space-y-4">
          {HeaderBlock}
          <div className="flex">
            <aside className="w-1/4">{sidebar}</aside>
            <div className="flex-1">{body}</div>
          </div>
          {FooterBlock}
        </div>
      );

    case "z-shape":
      return (
        <div className="p-4 space-y-4">
          {HeaderBlock}
          <div className="grid grid-cols-2 gap-4">{body}</div>
          {FooterBlock}
        </div>
      );

    case "card-block":
      return (
        <div>
          {HeaderBlock}
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 p-4">
            {body}
          </div>
          {FooterBlock}
        </div>
      );

    case "magazine":
      return (
        <div>
          {HeaderBlock}
          <div className="grid grid-cols-3 gap-4 p-4">
            <div className="col-span-2">{body}</div>
            <div>{sidebar}</div>
          </div>
          {FooterBlock}
        </div>
      );

    case "featured":
      return (
        <div>
          {HeaderBlock}
          <div className="p-4 flex flex-col lg:flex-row gap-4">
            <div className="flex-1">{body}</div>
            <div className="w-full lg:w-1/3">{sidebar}</div>
          </div>
          {FooterBlock}
        </div>
      );

    case "split-screen":
      return (
        <div className="grid grid-cols-1 md:grid-cols-2 h-screen">
          <div className="p-4 bg-gray-100">{sidebar}</div>
          <div className="p-4">{body}</div>
        </div>
      );

    case "interactive":
      return (
        <div className="flex flex-col h-screen">
          {HeaderBlock}
          <div className="flex-1 flex items-center justify-center p-8">
            {body}
          </div>
          {FooterBlock}
        </div>
      );

    case "two-column":
      return (
        <div>
          {HeaderBlock}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4">
            {body}
          </div>
          {FooterBlock}
        </div>
      );

    default:
      return <>{body}</>;
  }
};
