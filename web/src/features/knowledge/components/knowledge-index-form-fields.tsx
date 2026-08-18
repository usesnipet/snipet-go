import { DriverSelect } from "@/components/form/driver-select";
import { FormInput } from "@/components/form/input";
import { FieldGroup } from "@/components/ui/field";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useParams } from "react-router";

import { useListKnowledgeIndexDrivers } from "../hooks";

export function KnowledgeIndexFormFields() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { data: drivers } = useListKnowledgeIndexDrivers(tenant?.id ?? "");

  return (
    <FieldGroup>
      <FormInput name="name" label="Name" required />
      <DriverSelect
        name="driver"
        configName="configuration"
        drivers={drivers ?? []}
        label="Driver"
        placeholder="Select a driver"
      />
    </FieldGroup>
  );
}
