import React, { useState, useEffect } from "react";
import toast from "react-hot-toast";

interface WhatsAppModalProps {
  isOpen: boolean;
  onClose: () => void;
  name: string;
  phone: string;
  leadId: string;
  onSuccess: () => void;
}

export default function WhatsAppModal({
  isOpen,
  onClose,
  name,
  phone,
  leadId,
  onSuccess,
}: WhatsAppModalProps) {
  const [message, setMessage] = useState("");
  const [isSending, setIsSending] = useState(false);

  const getGreeting = (): string => {
    const hour = new Date().getHours();
    if (hour >= 4 && hour < 11) return "Selamat pagi";
    if (hour >= 11 && hour < 15) return "Selamat siang";
    if (hour >= 15 && hour < 19) return "Selamat sore";
    return "Selamat malam";
  };

  const getTemplate = (bizName: string): string =>
    `${getGreeting()}, Bpk/Ibu owner ${bizName} 👋\n\n` +
    `Perkenalkan, saya Aji Nursafiki — developer dari Risewise, sekaligus owner dari Siswanto Aki Jogja (akimobiljogja.com).\n\n` +
    `Dalam kesempatan ini, saya ingin menawarkan website sistem manajemen toko aki yang saya bangun langsung dari pengalaman menjalankan toko aki sendiri — jadi fitur-fiturnya dibuat untuk kebutuhan nyata di lapangan, bukan sekadar template.\n\n` +
    `Website ini memiliki beberapa fitur utama:\n` +
    `• Halaman utama sebagai company profile yang bisa disearch di pencarian google.\n` +
    `• Halaman Katalog aki\n` +
    `• Rekomendasi aki berdasarkan jenis kendaraan pelanggan\n` +
    `• Laporan transaksi harian dan bulanan\n` +
    `• Manajemen stok otomatis\n` +
    `• Invoice digital yang bisa langsung dikirim ke pelanggan\n` +
    `• Pelacakan garansi aki per pelanggan\n\n` +
    `Proposal lengkap bisa dilihat di sini:\n` +
    `- https://proposal.arkane.my.id\n\n` +
    `Referensi tampilan yang sudah berjalan:\n` +
    `- https://akimobiljogja.com\n\n` +
    `Jika Bpk/Ibu tertarik atau ingin mengetahui lebih lanjut, saya siap menjawab pertanyaan atau menjadwalkan demo gratis kapan saja.\n\n` +
    `Terima kasih atas waktunya, semoga bisa menjadi solusi yang bermanfaat 🙏\n\n` +
    `— Aji Nursafiki\n` +
    `Developer @ Risewise | Owner Siswanto Aki Jogja`;

  useEffect(() => {
    if (isOpen && name) {
      setMessage(getTemplate(name));
    }
  }, [isOpen, name]);

  if (!isOpen) return null;

  const handleSend = async () => {
    setIsSending(true);
    const rawPhone = phone.replace(/[^0-9]/g, "");
    const formattedPhone = rawPhone.startsWith("0")
      ? "62" + rawPhone.slice(1)
      : rawPhone;
      
    try {
      const res = await fetch("/api/wa/send", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ phone: formattedPhone, message: message.trim() })
      });
      
      if (res.ok) {
        toast.success("Message sent successfully via WhatsApp API!");
        
        // Mark lead as completed
        if (leadId) {
          try {
            await fetch("/api/leads/complete", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({ id: leadId })
            });
            onSuccess();
          } catch (e) {
            console.error("Failed to mark lead as completed:", e);
          }
        }

        onClose();
        setIsSending(false);
        return;
      }
    } catch (e) {
      console.error("Failed to send via API:", e);
    }
    
    // Fallback to wa.me if not connected or error
    setIsSending(false);
    toast.error("Failed to send via API. Opening WhatsApp Web...");
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
              disabled={!message.trim() || isSending}
              className="px-5 py-2 bg-green-500 hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed text-black font-bold rounded-lg text-[13px] transition cursor-pointer inline-flex items-center gap-1.5"
            >
              {isSending ? "Sending..." : "Send Message"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
