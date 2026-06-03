import React, { useState, useEffect } from "react";
import { QRCodeSVG } from "qrcode.react";

export default function WhatsAppConnect() {
  const [isConnected, setIsConnected] = useState(false);
  const [qrCode, setQrCode] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoggingOut, setIsLoggingOut] = useState(false);

  useEffect(() => {
    checkStatus();
    const interval = setInterval(checkStatus, 3000);
    return () => clearInterval(interval);
  }, []);

  const checkStatus = async () => {
    try {
      const res = await fetch("/api/wa/status");
      if (res.ok) {
        const data = await res.json();
        setIsConnected(data.connected);
        if (data.connected) {
          setQrCode(null); // Clear QR if connected
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoading(false);
    }
  };

  const getQR = async () => {
    try {
      const res = await fetch("/api/wa/qr");
      if (res.ok) {
        const data = await res.json();
        if (data.qr) {
          setQrCode(data.qr);
        }
      }
    } catch (e) {
      console.error(e);
    }
  };

  const handleLogout = async () => {
    if (!confirm("Are you sure you want to logout from WhatsApp?")) return;
    setIsLoggingOut(true);
    try {
      await fetch("/api/wa/logout", { method: "POST" });
      setIsConnected(false);
      setQrCode(null);
    } catch (e) {
      console.error(e);
    } finally {
      setIsLoggingOut(false);
    }
  };

  if (isLoading) {
    return (
      <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 mb-6 flex items-center justify-center">
        <div className="text-zinc-400 text-sm animate-pulse">Loading WhatsApp Status...</div>
      </div>
    );
  }

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 mb-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-bold text-zinc-100 flex items-center gap-2">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="w-5 h-5 text-green-500">
              <path d="M21 11.5a8.38 8.38 0 0 1-.9 3.8 8.5 8.5 0 0 1-7.6 4.7 8.38 8.38 0 0 1-3.8-.9L3 21l1.9-5.7a8.38 8.38 0 0 1-.9-3.8 8.5 8.5 0 0 1 4.7-7.6 8.38 8.38 0 0 1 3.8-.9h.5a8.48 8.48 0 0 1 8 8v.5z" />
            </svg>
            WhatsApp Connection
          </h2>
          <p className="text-sm text-zinc-400 mt-1">
            {isConnected 
              ? "Your WhatsApp account is connected. You can now send messages automatically." 
              : "Connect your WhatsApp to send messages automatically without opening the app."}
          </p>
        </div>
        
        <div>
          {isConnected ? (
            <button
              onClick={handleLogout}
              disabled={isLoggingOut}
              className="px-4 py-2 border border-red-500/30 text-red-400 hover:bg-red-500/10 hover:border-red-500/50 rounded-lg text-[13px] font-bold transition cursor-pointer disabled:opacity-50"
            >
              {isLoggingOut ? "Logging out..." : "Logout"}
            </button>
          ) : (
            <button
              onClick={getQR}
              className="px-4 py-2 bg-green-500 hover:bg-green-600 text-black font-bold rounded-lg text-[13px] transition cursor-pointer"
            >
              Get QR Code
            </button>
          )}
        </div>
      </div>

      {qrCode && !isConnected && (
        <div className="mt-6 p-6 border border-zinc-800 bg-zinc-950 rounded-lg flex flex-col items-center justify-center">
          <div className="bg-white p-4 rounded-xl mb-4">
            <QRCodeSVG value={qrCode} size={256} />
          </div>
          <p className="text-sm text-zinc-400 text-center max-w-md">
            Open WhatsApp on your phone, go to <strong>Linked Devices</strong> and scan this QR code to connect.
          </p>
          <button 
            onClick={getQR}
            className="mt-4 text-[12px] text-green-500 hover:text-green-400 underline cursor-pointer"
          >
            Refresh QR Code
          </button>
        </div>
      )}
    </div>
  );
}
