export interface Lead {
  id: string;
  name: string;
  category: string;
  rating: number;
  reviews: number;
  address: string;
  phone: string;
  website: string;
  maps_url: string;
  has_website: boolean;
  has_phone: boolean;
  tier: string;
  reason: string;
}

export interface ScraperStatus {
  running: boolean;
  done: boolean;
  error?: string;
  total: number;
  keyword: string;
}
