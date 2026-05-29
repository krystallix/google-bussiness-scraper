import React from "react";

interface TopbarProps {
  keyword: string;
  setKeyword: (kw: string) => void;
  limit: number;
  setLimit: (limit: number) => void;
  headful: boolean;
  setHeadful: (headful: boolean) => void;
  onStartScrape: () => void;
  onAddManualLead: () => void;
  disabledScrape: boolean;
}

export default function Topbar({
  keyword,
  setKeyword,
  limit,
  setLimit,
  headful,
  setHeadful,
  onStartScrape,
  onAddManualLead,
  disabledScrape,
}: TopbarProps) {
  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && !disabledScrape && keyword.trim()) {
      onStartScrape();
    }
  };

  return (
    <div className="px-6 py-4 border-b border-zinc-800 bg-zinc-900 flex items-center gap-3 flex-wrap">
      <div className="flex items-center gap-2.5 flex-1 max-w-lg">
        <div className="relative flex-1">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-400 pointer-events-none"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
          <input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder="e.g. coffee shop Bandung, toko aki Semarang..."
            className="w-full bg-zinc-950 border border-zinc-700 hover:border-zinc-600 focus:border-green-500 rounded-lg py-2 pl-9 pr-4 text-[14px] text-zinc-100 placeholder-zinc-500 outline-none transition"
            autoComplete="off"
          />
        </div>
      </div>

      <div className="flex items-center gap-2">
        <label htmlFor="limit" className="text-[12px] text-zinc-400 font-mono">
          Limit
        </label>
        <input
          id="limit"
          type="number"
          value={limit}
          onChange={(e) => setLimit(parseInt(e.target.value) || 0)}
          min="1"
          max="500"
          className="w-20 bg-zinc-950 border border-zinc-700 hover:border-zinc-600 focus:border-green-500 rounded-lg py-2 px-3 text-[14px] text-center font-mono text-zinc-100 outline-none transition"
        />
      </div>

      <div className="flex items-center gap-2">
        <span className="text-[12px] text-zinc-400">Show browser</span>
        <label className="relative inline-flex items-center cursor-pointer select-none">
          <input
            type="checkbox"
            checked={headful}
            onChange={(e) => setHeadful(e.target.checked)}
            className="sr-only peer"
          />
          <div className="w-9 h-5 bg-zinc-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-zinc-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-green-500"></div>
        </label>
      </div>

      <button
        onClick={onStartScrape}
        disabled={disabledScrape || !keyword.trim()}
        className="px-5 py-2 bg-green-500 hover:bg-green-600 text-black font-semibold rounded-lg text-[13px] transition active:scale-[0.98] disabled:opacity-50 disabled:cursor-not-allowed disabled:active:scale-100 cursor-pointer"
      >
        Start Scraping
      </button>

      <button
        onClick={onAddManualLead}
        className="px-5 py-2 border border-zinc-700 hover:border-zinc-500 hover:bg-zinc-800 text-zinc-100 font-semibold rounded-lg text-[13px] transition active:scale-[0.98] cursor-pointer"
      >
        Add Manual Lead
      </button>
    </div>
  );
}
