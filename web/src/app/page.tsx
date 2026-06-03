"use client";

import React, { useState, useEffect, useCallback, useMemo } from "react";
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

  // Stats summary for Sidebar
  const summaryText = useMemo(() => {
    if (allLeads.length === 0) return "No data yet";
    const high = allLeads.filter((l) => l.tier === "HIGH POTENTIAL").length;
    const std = allLeads.filter((l) => l.tier === "STANDARD").length;
    const skip = allLeads.filter((l) => l.tier === "SKIP").length;
    const comp = allLeads.filter((l) => l.tier === "COMPLETED").length;
    return `HIGH ${high} / STD ${std} / SKIP ${skip} / COMP ${comp}`;
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
      result = result.filter((l) => (l.tier === "HIGH POTENTIAL" || l.tier === "COMPLETED") && l.phone);
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
        setCurrentSection={setCurrentSection}
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
            setActiveFilter={setActiveFilter}
            searchQuery={searchQuery}
            setSearchQuery={setSearchQuery}
            count={filteredLeads.length}
          />

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
