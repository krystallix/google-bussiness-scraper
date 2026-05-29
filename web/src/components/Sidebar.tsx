import React from "react";

interface SidebarProps {
  currentSection: string;
  setCurrentSection: (section: string) => void;
  summary: string;
}

export default function Sidebar({
  currentSection,
  setCurrentSection,
  summary,
}: SidebarProps) {
  return (
    <aside className="w-56 shrink-0 bg-zinc-900 border-r border-zinc-800 flex flex-col py-5 select-none">
      <div className="flex items-center gap-2.5 px-5 pb-6 border-b border-zinc-800 mb-3">
        <div className="w-8 h-8 bg-green-500 rounded-md flex items-center justify-center font-semibold text-black">
          <svg
            viewBox="0 0 20 20"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="w-[18px] h-[18px]"
          >
            <circle cx="10" cy="10" r="3" />
            <path d="M10 2v2M10 16v2M2 10h2M16 10h2" />
            <path d="M4.93 4.93l1.41 1.41M13.66 13.66l1.41 1.41M4.93 15.07l1.41-1.41M13.66 6.34l1.41-1.41" />
          </svg>
        </div>
        <div className="font-bold text-[15px] tracking-tight leading-none text-zinc-100">
          G-Lead
          <span className="block text-[10px] text-zinc-400 font-normal tracking-[0.06em] uppercase font-mono mt-1">
            Scraper v2
          </span>
        </div>
      </div>

      <nav className="flex flex-col gap-0.5 px-2.5">
        <button
          onClick={() => setCurrentSection("scraper")}
          className={`flex items-center gap-2.5 py-2 px-2.5 rounded-lg text-[13px] font-medium transition cursor-pointer hover:bg-zinc-800 hover:text-zinc-100 ${
            currentSection === "scraper"
              ? "bg-green-500/10 text-green-500 hover:bg-green-500/10 hover:text-green-500"
              : "text-zinc-400"
          }`}
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="w-4 h-4 shrink-0"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
          Scraper
        </button>

        <button
          onClick={() => setCurrentSection("leads")}
          className={`flex items-center gap-2.5 py-2 px-2.5 rounded-lg text-[13px] font-medium transition cursor-pointer hover:bg-zinc-800 hover:text-zinc-100 ${
            currentSection === "leads"
              ? "bg-green-500/10 text-green-500 hover:bg-green-500/10 hover:text-green-500"
              : "text-zinc-400"
          }`}
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="w-4 h-4 shrink-0"
          >
            <path d="M9 11l3 3L22 4" />
            <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11" />
          </svg>
          Leads
        </button>

        <button
          onClick={() => setCurrentSection("whatsapp")}
          className={`flex items-center gap-2.5 py-2 px-2.5 rounded-lg text-[13px] font-medium transition cursor-pointer hover:bg-zinc-800 hover:text-zinc-100 ${
            currentSection === "whatsapp"
              ? "bg-green-500/10 text-green-500 hover:bg-green-500/10 hover:text-green-500"
              : "text-zinc-400"
          }`}
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="w-4 h-4 shrink-0"
          >
            <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
          </svg>
          WhatsApp
        </button>
      </nav>

      <div className="mt-auto px-5 pt-4 border-t border-zinc-800">
        <p className="text-[11px] text-zinc-400 font-mono tracking-wider">
          {summary || "No data yet"}
        </p>
      </div>
    </aside>
  );
}
