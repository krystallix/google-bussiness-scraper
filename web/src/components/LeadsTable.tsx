import React, { useMemo } from "react";
import { Lead } from "../types";

const TIER_ORDER = ["HIGH POTENTIAL", "STANDARD", "COMPLETED", "SKIP"];
const TIER_LABELS: Record<string, string> = {
  "HIGH POTENTIAL": "High Potential",
  STANDARD: "Standard",
  COMPLETED: "Completed",
  SKIP: "Skip",
};

interface LeadsTableProps {
  leads: Lead[];
  onEdit: (lead: Lead) => void;
  onDelete: (id: string) => void;
  onWhatsApp: (id: string, name: string, phone: string) => void;
  sortField: string;
  sortAsc: boolean;
  onSort: (field: string) => void;
  selectedIds: string[];
  onSelectionChange: (ids: string[]) => void;
  duplicatePhones: Set<string>;
}

export default function LeadsTable({
  leads,
  onEdit,
  onDelete,
  onWhatsApp,
  sortField: _sortField,
  sortAsc: _sortAsc,
  selectedIds,
  onSelectionChange,
  duplicatePhones,
}: LeadsTableProps) {
  const renderStars = (rating: number) => {
    if (!rating) return "—";
    const full = Math.floor(rating);
    const half = rating - full >= 0.25 && rating - full < 0.75 ? 1 : 0;
    const empty = 5 - full - half;
    return (
      <span className="text-amber-500 text-[12px] tracking-tighter">
        {"★".repeat(full)}
        {half ? "½" : ""}
        {"☆".repeat(Math.max(0, empty))}
      </span>
    );
  };

  const getTierColor = (tier: string) => {
    switch (tier) {
      case "HIGH POTENTIAL":
        return "bg-green-500/10 text-green-500 border border-green-500/20";
      case "STANDARD":
        return "bg-blue-500/10 text-blue-400 border border-blue-500/20";
      case "SKIP":
        return "bg-red-500/10 text-red-400 border border-red-500/20 opacity-70";
      case "COMPLETED":
        return "bg-purple-500/10 text-purple-400 border border-purple-500/20";
      default:
        return "bg-zinc-800 text-zinc-400 border border-zinc-700";
    }
  };

  const groupedLeads = useMemo(() => {
    const map = new Map<string, Lead[]>();
    for (const lead of leads) {
      const arr = map.get(lead.tier);
      if (arr) arr.push(lead);
      else map.set(lead.tier, [lead]);
    }
    const groups: { tier: string; leads: Lead[] }[] = [];
    for (const tier of TIER_ORDER) {
      if (map.has(tier)) groups.push({ tier, leads: map.get(tier)! });
    }
    return groups;
  }, [leads]);

  const isGroupSelected = (groupLeads: Lead[]) =>
    groupLeads.length > 0 && groupLeads.every((l) => selectedIds.includes(l.id));

  const toggleGroup = (groupLeads: Lead[]) => {
    const ids = groupLeads.map((l) => l.id);
    const allSelected = isGroupSelected(groupLeads);
    if (allSelected) {
      onSelectionChange(selectedIds.filter((id) => !ids.includes(id)));
    } else {
      const set = new Set(selectedIds);
      for (const id of ids) set.add(id);
      onSelectionChange(Array.from(set));
    }
  };

  const toggleOne = (id: string) => {
    if (selectedIds.includes(id)) {
      onSelectionChange(selectedIds.filter((sid) => sid !== id));
    } else {
      onSelectionChange([...selectedIds, id]);
    }
  };

  if (leads.length === 0) {
    return (
      <div className="border border-zinc-800 rounded-lg overflow-hidden bg-zinc-900/50">
        <div className="flex flex-col items-center justify-center py-20 px-4 text-zinc-400 gap-3 text-center">
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="1.5"
            className="w-10 h-10 opacity-30"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
          <p className="text-[14px]">No results match this filter</p>
          <small className="text-[12px] opacity-60 font-mono">
            Try a different filter or search term
          </small>
        </div>
      </div>
    );
  }

  const checkboxClass =
    "appearance-none w-4 h-4 rounded bg-zinc-900 border border-zinc-700 checked:bg-green-500 checked:border-green-500 focus:ring-1 focus:ring-green-500 focus:outline-none cursor-pointer transition";

  return (
    <div className="border border-zinc-800 rounded-lg overflow-hidden bg-zinc-900">
      <div className="overflow-x-auto w-full">
        <table className="w-full border-collapse">
          <thead className="bg-zinc-800/80 border-b border-zinc-800 text-[11px] font-mono tracking-wider text-zinc-400 uppercase select-none">
            <tr>
              <th className="py-3 px-4 w-10"></th>
              <th className="py-3 px-4 text-left font-semibold">Business Name</th>
              <th className="py-3 px-4 text-left font-semibold">Rating</th>
              <th className="py-3 px-4 text-left font-semibold">Reviews</th>
              <th className="py-3 px-4 text-left font-semibold">Category</th>
              <th className="py-3 px-4 text-left font-semibold">Phone</th>
              <th className="py-3 px-4 text-left font-semibold">Website</th>
              <th className="py-3 px-4 text-left font-semibold">Tier</th>
              <th className="py-3 px-4 text-left font-semibold">Action</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-800">
            {groupedLeads.map((group) => (
              <React.Fragment key={group.tier}>
                <tr className="bg-zinc-800/40">
                  <td className="py-2 px-4">
                    <input
                      type="checkbox"
                      checked={isGroupSelected(group.leads)}
                      onChange={() => toggleGroup(group.leads)}
                      className={checkboxClass}
                    />
                  </td>
                  <td colSpan={8} className="py-2 px-0">
                    <div className="flex items-center gap-2">
                      <span className="text-[11px] font-mono font-bold tracking-wider uppercase text-zinc-300">
                        {TIER_LABELS[group.tier] || group.tier}
                      </span>
                      <span className="text-[11px] font-mono text-zinc-500">
                        {group.leads.length}
                      </span>
                    </div>
                  </td>
                </tr>
                {group.leads.map((l) => (
                  <tr
                    key={l.id}
                    className="hover:bg-zinc-800/30 transition duration-75"
                  >
                    <td className="py-3 px-4">
                      <input
                        type="checkbox"
                        checked={selectedIds.includes(l.id)}
                        onChange={() => toggleOne(l.id)}
                        className={checkboxClass}
                      />
                    </td>
                    <td className="py-3 px-4 max-w-[220px]">
                      <a
                        href={l.maps_url}
                        target="_blank"
                        rel="noreferrer"
                        className="font-medium text-zinc-100 hover:text-green-500 transition text-[13.5px] block truncate"
                      >
                        {l.name}
                      </a>
                      {l.address && (
                        <div
                          className="text-[11px] text-zinc-400 mt-0.5 truncate font-mono"
                          title={l.address}
                        >
                          {l.address}
                        </div>
                      )}
                    </td>
                    <td className="py-3 px-4 whitespace-nowrap">
                      <div className="flex items-center gap-1.5">
                        {renderStars(l.rating)}
                        <span className="text-[12px] font-mono text-zinc-400">
                          {l.rating > 0 ? l.rating.toFixed(1) : "—"}
                        </span>
                      </div>
                    </td>
                    <td className="py-3 px-4 text-[13px] font-mono text-zinc-300">
                      {l.reviews > 0 ? l.reviews.toLocaleString() : "—"}
                    </td>
                    <td className="py-3 px-4 text-[13px] text-zinc-300 truncate max-w-[120px]">
                      {l.category || "—"}
                    </td>
                    <td className="py-3 px-4 text-[12px] font-mono text-zinc-300 whitespace-nowrap">
                      {l.phone ? (
                        <div className="flex items-center gap-1.5">
                          <span>{l.phone}</span>
                          {duplicatePhones.has(l.phone) && (
                            <span className="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-bold bg-amber-500/15 text-amber-400 border border-amber-500/25 font-mono leading-none">
                              dup
                            </span>
                          )}
                        </div>
                      ) : "—"}
                    </td>
                    <td className="py-3 px-4 whitespace-nowrap">
                      {l.has_website ? (
                        <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-semibold bg-zinc-800 text-zinc-400 border border-zinc-700 font-mono">
                          has site
                        </span>
                      ) : (
                        <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-semibold bg-amber-500/10 text-amber-500 border border-amber-500/20 font-mono">
                          no site
                        </span>
                      )}
                    </td>
                    <td className="py-3 px-4 whitespace-nowrap">
                      <span
                        className={`inline-block px-2.5 py-0.5 rounded text-[11px] font-bold tracking-wider font-mono ${getTierColor(
                          l.tier
                        )}`}
                        title={l.reason}
                      >
                        {l.tier.split(" ")[0]}
                      </span>
                    </td>
                    <td className="py-3 px-4 whitespace-nowrap">
                      <div className="flex items-center gap-2">
                        {(l.tier === "HIGH POTENTIAL" ||
                          l.tier === "STANDARD" ||
                          l.tier === "COMPLETED") &&
                          l.phone && (
                            <button
                              onClick={() =>
                                onWhatsApp(l.id, l.name, l.phone)
                              }
                              className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-[11px] font-semibold transition cursor-pointer font-mono ${
                                l.tier === "COMPLETED"
                                  ? "bg-purple-500/10 hover:bg-purple-500/20 text-purple-400 border border-purple-500/30"
                                  : "bg-green-500/10 hover:bg-green-500/20 text-green-400 border border-green-500/30"
                              }`}
                            >
                              <svg
                                viewBox="0 0 24 24"
                                fill="currentColor"
                                className="w-3 h-3"
                              >
                                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z" />
                                <path d="M11.998 0C5.372 0 0 5.373 0 12c0 2.117.554 4.107 1.523 5.832L0 24l6.335-1.652A11.954 11.954 0 0011.998 24C18.626 24 24 18.627 24 12S18.626 0 11.998 0zm0 21.818a9.818 9.818 0 01-5.01-1.374l-.36-.213-3.76.984.998-3.65-.233-.375A9.82 9.82 0 012.18 12c0-5.42 4.4-9.818 9.818-9.818 5.42 0 9.82 4.398 9.82 9.818 0 5.42-4.4 9.818-9.82 9.818z" />
                              </svg>
                              {l.tier === "COMPLETED"
                                ? "Resend WA"
                                : "WhatsApp"}
                            </button>
                          )}
                        <button
                          onClick={() => onEdit(l)}
                          className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-[11px] font-semibold bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 border border-blue-500/30 transition cursor-pointer font-mono"
                        >
                          <svg
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                            className="w-3 h-3"
                          >
                            <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                            <path d="M18.5 2.5a2.121 2.121 0 1 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                          </svg>
                          Edit
                        </button>
                        <button
                          onClick={() => onDelete(l.id)}
                          className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md text-[11px] font-semibold bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/30 transition cursor-pointer font-mono"
                        >
                          <svg
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            strokeWidth="2"
                            className="w-3 h-3"
                          >
                            <polyline points="3 6 5 6 21 6" />
                            <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                          </svg>
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </React.Fragment>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
