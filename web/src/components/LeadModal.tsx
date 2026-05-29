import React, { useState, useEffect } from "react";
import { Lead } from "../types";

interface LeadModalProps {
  isOpen: boolean;
  onClose: () => void;
  lead: Lead | null;
  onSave: (leadData: any) => Promise<void>;
}

export default function LeadModal({
  isOpen,
  onClose,
  lead,
  onSave,
}: LeadModalProps) {
  const [name, setName] = useState("");
  const [category, setCategory] = useState("");
  const [rating, setRating] = useState(0);
  const [reviews, setReviews] = useState(0);
  const [address, setAddress] = useState("");
  const [phone, setPhone] = useState("");
  const [website, setWebsite] = useState("");
  const [mapsUrl, setMapsUrl] = useState("");
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (lead) {
      setName(lead.name || "");
      setCategory(lead.category || "");
      setRating(lead.rating || 0);
      setReviews(lead.reviews || 0);
      setAddress(lead.address || "");
      setPhone(lead.phone || "");
      setWebsite(lead.website || "");
      setMapsUrl(lead.maps_url || "");
    } else {
      setName("");
      setCategory("");
      setRating(0);
      setReviews(0);
      setAddress("");
      setPhone("");
      setWebsite("");
      setMapsUrl("");
    }
  }, [lead, isOpen]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;

    setIsSaving(true);
    try {
      const payload: any = {
        name: name.trim(),
        category: category.trim(),
        rating: parseFloat(rating as any) || 0,
        reviews: parseInt(reviews as any) || 0,
        address: address.trim(),
        phone: phone.trim(),
        website: website.trim(),
        maps_url: mapsUrl.trim(),
      };
      if (lead) {
        payload.id = lead.id;
      }
      await onSave(payload);
      onClose();
    } catch (err) {
      alert("Failed to save: " + (err as Error).message);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 backdrop-blur-sm p-4" onClick={onClose}>
      <div
        className="bg-zinc-900 border border-zinc-800 rounded-xl p-6 w-full max-w-lg flex flex-col gap-4 shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-zinc-800 pb-3">
          <span className="text-[15px] font-semibold text-zinc-100">
            {lead ? `Edit Business: ${lead.name}` : "Add Manual Lead"}
          </span>
          <button
            onClick={onClose}
            className="text-zinc-400 hover:text-zinc-100 text-[18px] leading-none cursor-pointer"
          >
            &times;
          </button>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col gap-3">
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
                Business Name *
              </label>
              <input
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 placeholder-zinc-650 outline-none transition"
              />
            </div>
            <div>
              <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
                Category
              </label>
              <input
                type="text"
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 placeholder-zinc-650 outline-none transition"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
                Rating (0.0 - 5.0)
              </label>
              <input
                type="number"
                min="0"
                max="5"
                step="0.1"
                value={rating}
                onChange={(e) => setRating(parseFloat(e.target.value) || 0)}
                className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 outline-none transition"
              />
            </div>
            <div>
              <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
                Reviews Count
              </label>
              <input
                type="number"
                min="0"
                value={reviews}
                onChange={(e) => setReviews(parseInt(e.target.value) || 0)}
                className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 outline-none transition"
              />
            </div>
          </div>

          <div>
            <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
              Address
            </label>
            <input
              type="text"
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 outline-none transition"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
                Phone Number
              </label>
              <input
                type="text"
                placeholder="e.g. 08123456789"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 placeholder-zinc-650 outline-none transition"
              />
            </div>
            <div>
              <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
                Website URL
              </label>
              <input
                type="text"
                placeholder="e.g. https://..."
                value={website}
                onChange={(e) => setWebsite(e.target.value)}
                className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 placeholder-zinc-650 outline-none transition"
              />
            </div>
          </div>

          <div>
            <label className="text-[12px] text-zinc-400 block mb-1 font-medium">
              Google Maps URL
            </label>
            <input
              type="text"
              placeholder="e.g. https://maps.google.com/..."
              value={mapsUrl}
              onChange={(e) => setMapsUrl(e.target.value)}
              className="w-full bg-zinc-950 border border-zinc-800 hover:border-zinc-700 focus:border-green-500 rounded-lg py-2 px-3 text-[13px] text-zinc-100 placeholder-zinc-650 outline-none transition"
            />
          </div>

          <div className="flex gap-2 justify-end mt-4 pt-3 border-t border-zinc-800">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 border border-zinc-700 hover:border-zinc-500 rounded-lg text-zinc-400 hover:text-zinc-100 text-[13px] transition cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSaving || !name.trim()}
              className="px-5 py-2 bg-green-500 hover:bg-green-600 disabled:opacity-50 disabled:cursor-not-allowed text-black font-semibold rounded-lg text-[13px] transition cursor-pointer"
            >
              {isSaving ? "Saving..." : "Save Lead"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
