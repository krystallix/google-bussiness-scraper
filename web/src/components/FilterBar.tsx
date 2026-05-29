import React from "react";

interface FilterBarProps {
  activeFilter: string;
  setActiveFilter: (filter: string) => void;
  searchQuery: string;
  setSearchQuery: (q: string) => void;
  count: number;
}

export default function FilterBar({
  activeFilter,
  setActiveFilter,
  searchQuery,
  setSearchQuery,
  count,
}: FilterBarProps) {
  const filters = [
    { label: "All", value: "all", className: "border-zinc-700 hover:border-zinc-500" },
    {
      label: "High Potential",
      value: "HIGH POTENTIAL",
      className: "border-green-800/40 hover:border-green-500",
      activeClass: "bg-green-500 text-black border-green-500",
    },
    {
      label: "Standard",
      value: "STANDARD",
      className: "border-blue-800/40 hover:border-blue-500",
      activeClass: "bg-blue-500 text-white border-blue-500",
    },
    {
      label: "Skip",
      value: "SKIP",
      className: "border-red-800/40 hover:border-red-500",
      activeClass: "bg-red-500 text-white border-red-500",
    },
  ];

  return (
    <div className="flex items-center gap-2 mb-4 flex-wrap select-none">
      {filters.map((f) => {
        const isActive = activeFilter === f.value;
        const customActive = f.activeClass || "bg-zinc-100 text-black border-zinc-100";
        return (
          <button
            key={f.value}
            onClick={() => setActiveFilter(f.value)}
            className={`px-3.5 py-1.5 rounded-full text-[12px] font-medium border font-mono transition cursor-pointer ${
              isActive ? customActive : `bg-transparent text-zinc-400 ${f.className}`
            }`}
          >
            {f.label}
          </button>
        );
      })}

      <div className="ml-auto flex items-center gap-2">
        <input
          type="text"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          placeholder="Filter by name, phone, or address..."
          className="px-3 py-1.5 bg-zinc-900 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg text-[13px] text-zinc-100 placeholder-zinc-500 outline-none transition w-52 focus:w-64"
        />
        <span className="text-[11px] font-mono text-zinc-400 bg-zinc-900 border border-zinc-800 px-3 py-1 rounded-full">
          {count} results
        </span>
      </div>
    </div>
  );
}
