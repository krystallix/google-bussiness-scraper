"use client";

import React, { useState, useEffect, useCallback, useMemo } from "react";
import toast from "react-hot-toast";
import Sidebar from "@/components/Sidebar";
import Topbar from "@/components/Topbar";
import Statusbar from "@/components/Statusbar";
import FilterBar from "@/components/FilterBar";
import LeadsTable from "@/components/LeadsTable";
import LeadModal from "@/components/LeadModal";
import WhatsAppModal from "@/components/WhatsAppModal";
import WhatsAppConnect from "@/components/WhatsAppConnect";
import { Lead, ScraperStatus } from "@/types";

export default function Dashboard() {
  const [currentSection, setCurrentSection] = useState("scraper");
  const [allLeads, setAllLeads] = useState<Lead[]>([]);
  const [activeFilter, setActiveFilter] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [sortField, setSortField] = useState("tier");
  const [sortAsc, setSortAsc] = useState(true);
  const [selectedIds, setSelectedIds] = useState<string[]>([]);

  // Scraper control state
  const [keyword, setKeyword] = useState("");
  const [limit, setLimit] = useState(50);
  const [headful, setHeadful] = useState(false);
  const [status, setStatus] = useState<"running" | "done" | "error" | "ready">("ready");
  const [statusMessage, setStatusMessage] = useState("Ready. Enter a keyword and start scraping.");
  const [isScraping, setIsScraping] = useState(false);

  // Modals state
  const [isLeadModalOpen, setIsLeadModalOpen] = useState(false);
  const [activeLead, setActiveLead] = useState<Lead | null>(null);
  const [isWaModalOpen, setIsWaModalOpen] = useState(false);
  const [waName, setWaName] = useState("");
  const [waPhone, setWaPhone] = useState("");
  const [waLeadId, setWaLeadId] = useState("");

  const loadResults = useCallback(async () => {
    try {
      const res = await fetch("/api/results");
      if (!res.ok) throw new Error("Failed to load results");
      const data: Lead[] = await res.json();
      setAllLeads(data);

      if (data.length > 0) {
        setStatus("done");
        setStatusMessage(`Loaded ${data.length} leads from database.`);
      } else {
        setStatus("done");
        setStatusMessage("Ready");
      }
    } catch (e) {
      console.error(e);
    }
  }, []);

  const pollStatus = useCallback(() => {
    const timer = setInterval(async () => {
      try {
        const res = await fetch("/api/status");
        if (!res.ok) throw new Error("Failed to check status");
        const data: ScraperStatus = await res.json();

        if (data.error) {
          setStatus("error");
          setStatusMessage("Error: " + data.error);
          setIsScraping(false);
          clearInterval(timer);
          return;
        }

        if (data.running) {
          setStatus("running");
          setStatusMessage(`Scraping "${data.keyword}"... (browser working)`);
          return;
        }

        if (data.done) {
          await loadResults();
          setStatus("done");
          setStatusMessage(`Done — ${data.total} results for "${data.keyword}"`);
          setIsScraping(false);
          clearInterval(timer);
        }
      } catch (err) {
        console.error(err);
      }
    }, 2000);

    return () => clearInterval(timer);
  }, [loadResults]);

  // Initial load
  useEffect(() => {
    let cleanupPoll: (() => void) | undefined;

    const checkInitialStatus = async () => {
      try {
        const res = await fetch("/api/status");
        if (!res.ok) throw new Error("Status check failed");
        const statusData: ScraperStatus = await res.json();

        if (statusData.running) {
          setStatus("running");
          setStatusMessage(`Scraping "${statusData.keyword}"... (browser working)`);
          setIsScraping(true);
          cleanupPoll = pollStatus();
        } else {
          await loadResults();
        }
      } catch (e) {
        await loadResults();
      }
    };

    checkInitialStatus();

    return () => {
      if (cleanupPoll) cleanupPoll();
    };
  }, [loadResults, pollStatus]);

  const startScrape = async () => {
    if (!keyword.trim()) return;

    setIsScraping(true);
    setStatus("running");
    setStatusMessage("Starting scraper for: " + keyword);
    setAllLeads([]);

    try {
      const res = await fetch("/api/scrape", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ keyword, limit, headful, delay: 1.5 }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Failed to start");
      pollStatus();
    } catch (e) {
      setStatus("error");
      setStatusMessage("Error: " + (e as Error).message);
      setIsScraping(false);
    }
  };

  const handleSaveLead = async (payload: any) => {
    const endpoint = payload.id ? "/api/leads/update" : "/api/leads/add";
    const res = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    const data = await res.json();
    if (!res.ok) throw new Error(data.error || "Failed to save lead");

    await loadResults();
  };

  const handleDeleteLead = async (id: string) => {
    const lead = allLeads.find((l) => l.id === id);
    if (!lead) return;

    if (!confirm(`Are you sure you want to delete "${lead.name}"?`)) return;

    try {
      const res = await fetch("/api/leads/delete", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Failed to delete");

      await loadResults();
    } catch (e) {
      alert("Error deleting lead: " + (e as Error).message);
    }
  };

  const handleWhatsApp = (id: string, name: string, phone: string) => {
    setWaLeadId(id);
    setWaName(name);
    setWaPhone(phone);
    setIsWaModalOpen(true);
  };
  const getGreeting = (): string => {
    const hour = new Date().getHours();
    if (hour >= 4 && hour < 11) return "Good morning";
    if (hour >= 11 && hour < 15) return "Good afternoon";
    if (hour >= 15 && hour < 19) return "Good afternoon";
    return "Good evening";
  };

  const getWATemplate = (bizName: string): string =>
    `${getGreeting()},\n\n` +
    `I checked ${bizName} on Google Maps — great rating and plenty of reviews. ` +
    `It looks like you don't have a website yet.\n\n` +
    `Does ${bizName} have a website?\n\n` +
    `Most people search Google first before visiting a business. ` +
    `A website makes it easy for potential customers to find you. ` +
    `You can also add features like online ordering, booking, or product catalogs.\n\n` +
    `Would you like a quick audit of what a website for ${bizName} could look like?`;

  const handleBulkWhatsApp = async () => {
    if (selectedIds.length === 0) return;

    if (selectedIds.length === 1) {
      const lead = allLeads.find((l) => l.id === selectedIds[0]);
      if (lead && lead.phone) {
        handleWhatsApp(lead.id, lead.name, lead.phone);
      }
      return;
    }

    const selectedLeads = allLeads.filter(
      (l) => selectedIds.includes(l.id) && l.phone
    );

    if (selectedLeads.length === 0) {
      toast.error("No selected leads have phone numbers");
      return;
    }

    if (!confirm(`Send WhatsApp to ${selectedLeads.length} selected lead(s)?`)) return;

    toast.success(`Processing ${selectedLeads.length} lead(s)...`);

    for (const lead of selectedLeads) {
      const rawPhone = lead.phone.replace(/[^0-9]/g, "");
      const formattedPhone = rawPhone.startsWith("0")
        ? "62" + rawPhone.slice(1)
        : rawPhone;

      const message = getWATemplate(lead.name);

      let apiOk = false;
      try {
        const res = await fetch("/api/wa/send", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ phone: formattedPhone, message }),
        });
        apiOk = res.ok;
      } catch { }

      if (!apiOk) {
        const url = `https://wa.me/${formattedPhone}?text=${encodeURIComponent(message)}`;
        window.open(url, "_blank");
        await new Promise((r) => setTimeout(r, 800));
      }
    }

    for (const id of selectedIds) {
      try {
        await fetch("/api/leads/complete", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ id }),
        });
      } catch { }
    }

    setSelectedIds([]);
    toast.success(`Done! Processed ${selectedLeads.length} lead(s)`);
    await loadResults();
  };

  const handleBulkDelete = async () => {
    if (selectedIds.length === 0) return;

    if (!confirm(`Delete ${selectedIds.length} selected lead(s)?`)) return;

    for (const id of selectedIds) {
      try {
        await fetch("/api/leads/delete", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ id }),
        });
      } catch {}
    }

    setSelectedIds([]);
    toast.success(`Deleted ${selectedIds.length} lead(s)`);
    await loadResults();
  };

  // Stats summary for Sidebar
  const summaryText = useMemo(() => {
    if (allLeads.length === 0) return "No data yet";
    const high = allLeads.filter((l) => l.tier === "HIGH POTENTIAL").length;
    const std = allLeads.filter((l) => l.tier === "STANDARD").length;
    const skip = allLeads.filter((l) => l.tier === "SKIP").length;
    const comp = allLeads.filter((l) => l.tier === "COMPLETED").length;
    return `HIGH ${high} / STD ${std} / SKIP ${skip} / COMP ${comp}`;
  }, [allLeads]);

  const duplicatePhones = useMemo(() => {
    const freq = new Map<string, number>();
    for (const l of allLeads) {
      if (l.phone) {
        freq.set(l.phone, (freq.get(l.phone) || 0) + 1);
      }
    }
    return new Set(Array.from(freq.entries()).filter(([_, c]) => c > 1).map(([p]) => p));
  }, [allLeads]);

  // Filter and sort computation
  const filteredLeads = useMemo(() => {
    let result = allLeads;

    // Filter by tier
    if (activeFilter !== "all") {
      result = result.filter((l) => l.tier === activeFilter);
    }

    // Filter by sidebar section
    if (currentSection === "whatsapp") {
      // In WhatsApp section, show only high potential or completed with phone numbers
      result = result.filter((l) => (l.tier === "HIGH POTENTIAL" || l.tier === "STANDARD" || l.tier === "COMPLETED") && l.phone);
    } else if (currentSection === "leads") {
      // In Leads section, show only high potential, standard, and completed
      result = result.filter((l) => l.tier === "HIGH POTENTIAL" || l.tier === "STANDARD" || l.tier === "COMPLETED");
    }

    // Filter by search query
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      result = result.filter(
        (l) =>
          l.name.toLowerCase().includes(query) ||
          l.phone.toLowerCase().includes(query) ||
          l.address.toLowerCase().includes(query)
      );
    }

    // Sort results
    const tierOrder: Record<string, number> = {
      "HIGH POTENTIAL": 0,
      STANDARD: 1,
      "COMPLETED": 2,
      SKIP: 3,
    };

    result = [...result].sort((a, b) => {
      let va: any, vb: any;
      if (sortField === "name") {
        va = a.name;
        vb = b.name;
      } else if (sortField === "rating") {
        va = a.rating;
        vb = b.rating;
      } else if (sortField === "reviews") {
        va = a.reviews;
        vb = b.reviews;
      } else if (sortField === "tier") {
        va = tierOrder[a.tier] ?? 9;
        vb = tierOrder[b.tier] ?? 9;
      } else {
        return 0;
      }

      if (va < vb) return sortAsc ? -1 : 1;
      if (va > vb) return sortAsc ? 1 : -1;
      return 0;
    });

    return result;
  }, [allLeads, activeFilter, currentSection, searchQuery, sortField, sortAsc]);

  const handleFilterChange = (filter: string) => {
    setActiveFilter(filter);
    setSelectedIds([]);
  };

  const handleSectionChange = (section: string) => {
    setCurrentSection(section);
    setSelectedIds([]);
  };

  const handleSort = (field: string) => {
    if (sortField === field) {
      setSortAsc(!sortAsc);
    } else {
      setSortField(field);
      setSortAsc(true);
    }
  };

  return (
    <div className="flex h-screen overflow-hidden bg-zinc-950 font-sans">
      <Sidebar
        currentSection={currentSection}
        setCurrentSection={handleSectionChange}
        summary={summaryText}
      />

      <div className="flex-1 flex flex-col overflow-hidden">
        <Topbar
          keyword={keyword}
          setKeyword={setKeyword}
          limit={limit}
          setLimit={setLimit}
          headful={headful}
          setHeadful={setHeadful}
          onStartScrape={startScrape}
          onAddManualLead={() => {
            setActiveLead(null);
            setIsLeadModalOpen(true);
          }}
          disabledScrape={isScraping}
        />

        <Statusbar status={status} message={statusMessage} />

        <main className="flex-1 overflow-auto p-6">
          {isScraping && (
            <div className="h-0.5 bg-zinc-800 relative overflow-hidden rounded-full mb-4">
              <div className="absolute top-0 left-[-60%] w-3/5 h-full bg-green-500 animate-[slide_1.2s_ease-in-out_infinite]" />
            </div>
          )}

          <FilterBar
            activeFilter={activeFilter}
            setActiveFilter={handleFilterChange}
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            count={filteredLeads.length}
          />

          {selectedIds.length > 0 && (
            <div className="flex items-center gap-3 px-4 py-2 mb-4 bg-zinc-800/60 border border-zinc-700 rounded-lg">
              <span className="text-[13px] text-zinc-300 font-mono">
                {selectedIds.length} selected
              </span>
              <button
                onClick={handleBulkWhatsApp}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-[12px] font-semibold bg-green-500/10 hover:bg-green-500/20 text-green-400 border border-green-500/30 transition cursor-pointer font-mono"
              >
                <svg viewBox="0 0 24 24" fill="currentColor" className="w-3.5 h-3.5">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z" />
                  <path d="M11.998 0C5.372 0 0 5.373 0 12c0 2.117.554 4.107 1.523 5.832L0 24l6.335-1.652A11.954 11.954 0 0011.998 24C18.626 24 24 18.627 24 12S18.626 0 11.998 0zm0 21.818a9.818 9.818 0 01-5.01-1.374l-.36-.213-3.76.984.998-3.65-.233-.375A9.82 9.82 0 012.18 12c0-5.42 4.4-9.818 9.818-9.818 5.42 0 9.82 4.398 9.82 9.818 0 5.42-4.4 9.818-9.82 9.818z" />
                </svg>
                Send WA to Selected
              </button>
              <button
                onClick={handleBulkDelete}
                className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-[12px] font-semibold bg-red-500/10 hover:bg-red-500/20 text-red-400 border border-red-500/30 transition cursor-pointer font-mono"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="w-3.5 h-3.5">
                  <polyline points="3 6 5 6 21 6" />
                  <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                </svg>
                Delete Selected
              </button>
              <button
                onClick={() => setSelectedIds([])}
                className="text-[12px] text-zinc-500 hover:text-zinc-300 transition cursor-pointer font-mono"
              >
                Clear
              </button>
            </div>
          )}

          {currentSection === "whatsapp" && <WhatsAppConnect />}

          <LeadsTable
            leads={filteredLeads}
            onEdit={(lead) => {
              setActiveLead(lead);
              setIsLeadModalOpen(true);
            }}
            onDelete={handleDeleteLead}
            onWhatsApp={handleWhatsApp}
            sortField={sortField}
            sortAsc={sortAsc}
            onSort={handleSort}
            selectedIds={selectedIds}
            onSelectionChange={setSelectedIds}
            duplicatePhones={duplicatePhones}
          />
        </main>
      </div>

      <LeadModal
        isOpen={isLeadModalOpen}
        onClose={() => setIsLeadModalOpen(false)}
        lead={activeLead}
        onSave={handleSaveLead}
      />

      <WhatsAppModal
        isOpen={isWaModalOpen}
        onClose={() => setIsWaModalOpen(false)}
        name={waName}
        phone={waPhone}
        leadId={waLeadId}
        onSuccess={loadResults}
      />
    </div>
  );
}
