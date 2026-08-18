import { create } from "zustand";

import type { Tenant } from "./schemas";
import type { Member } from "@/models/member";

interface TenantStore {
  tenant: Tenant | null;
  setTenant: (tenant: Tenant | null) => void;
  clearTenant: () => void;
  member: Member | null;
  setMember: (member: Member | null) => void;
  clearMember: () => void;
}

export const useTenantStore = create<TenantStore>((set) => ({
  tenant: null,
  setTenant: (tenant) => set({ tenant }),
  clearTenant: () => set({ tenant: null }),

  member: null,
  setMember: (member) => set({ member }),
  clearMember: () => set({ member: null }),
}));