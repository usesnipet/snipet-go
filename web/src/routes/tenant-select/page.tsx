import { TenantSelect } from "@/features/tenant/components/tenant-select";

export const TenantSelectPage = () => {
  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
      <TenantSelect />
    </div>
  );
};
