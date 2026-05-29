import React from "react";

interface StatusbarProps {
  status: "running" | "done" | "error" | "ready";
  message: string;
}

export default function Statusbar({ status, message }: StatusbarProps) {
  const getStatusColor = () => {
    switch (status) {
      case "running":
        return "bg-amber-500 animate-pulse";
      case "done":
        return "bg-green-500";
      case "error":
        return "bg-red-500";
      default:
        return "bg-zinc-700";
    }
  };

  return (
    <div className="px-6 py-2 bg-zinc-950 border-b border-zinc-800 flex items-center gap-4 text-[12px] font-mono text-zinc-400 min-h-[36px]">
      <div className={`w-2 h-2 rounded-full ${getStatusColor()}`} />
      <span>{message}</span>
    </div>
  );
}
