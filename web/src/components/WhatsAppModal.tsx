import React, { useState, useEffect } from "react";

interface WhatsAppModalProps {
  isOpen: boolean;
  onClose: () => void;
  name: string;
  phone: string;
}

export default function WhatsAppModal({
  isOpen,
  onClose,
  name,
  phone,
}: WhatsAppModalProps) {
  const [message, setMessage] = useState("");

  const getTemplate = (bizName: string) =>
    `Halo ${bizName}, perkenalkan saya dari tim kami.\n\n` +
    `Kami melihat bisnis Anda di Google Maps dan ingin menawarkan:\n` +
    `- Website profesional dengan landing page\n` +
    `- Sistem manajemen stok digital\n` +
    `- Invoice digital & laporan penjualan\n\n` +
    `Apakah Anda tertarik untuk mengetahui lebih lanjut? ` +
    `Kami siap berdiskusi kapan saja.`;

  useEffect(() => {
    if (isOpen && name) {
      setMessage(getTemplate(name));
    }
  }, [isOpen, name]);

  if (!isOpen) return null;

  const handleSend = () => {
    const rawPhone = phone.replace(/[^0-9]/g, "");
    const formattedPhone = rawPhone.startsWith("0")
      ? "62" + rawPhone.slice(1)
      : rawPhone;
    const url = `https://wa.me/${formattedPhone}?text=${encodeURIComponent(
      message.trim()
    )}`;
    window.open(url, "_blank");
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 backdrop-blur-sm p-4" onClick={onClose}>
      <div
        className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 w-full max-w-lg flex flex-col gap-4 shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-zinc-850 pb-3">
          <span className="text-[15px] font-semibold text-zinc-100">
            Compose WhatsApp Message
          </span>
          <button
            onClick={onClose}
            className="text-zinc-400 hover:text-zinc-100 text-[18px] leading-none cursor-pointer"
          >
            &times;
          </button>
        </div>

        <div className="flex flex-col gap-3">
          <div>
            <label className="text-[12px] text-zinc-400 block mb-1 font-mono">
              Recipient
            </label>
            <div className="text-[13.5px] font-semibold text-zinc-100 bg-zinc-950 border border-zinc-800 rounded-lg py-2 px-3">
              {name} ({phone})
            </div>
          </div>

          <div>
            <label className="text-[12px] text-zinc-400 block mb-1 font-mono">
              Message
            </label>
            <textarea
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Type your message..."
              className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2.5 px-3 text-[13px] text-zinc-100 outline-none transition min-h-[140px] resize-y leading-relaxed"
            />
          </div>

          <div className="flex gap-2 justify-end mt-4 pt-3 border-t border-zinc-800">
            <button
              onClick={onClose}
              className="px-4 py-2 border border-zinc-700 hover:border-zinc-500 rounded-lg text-zinc-400 hover:text-zinc-100 text-[13px] transition cursor-pointer"
            >
              Cancel
            </button>
            <button
              onClick={handleSend}
              disabled={!message.trim()}
              className="px-5 py-2 bg-green-500 hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed text-black font-bold rounded-lg text-[13px] transition cursor-pointer inline-flex items-center gap-1.5"
            >
              Open WhatsApp
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
