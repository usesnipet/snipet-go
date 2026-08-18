import { DriverSelect } from "@/components/form/driver-select";
import { FormInput } from "@/components/form/input";
import { FieldGroup } from "@/components/ui/field";
import { useFindBySlugTenant } from "@/features/tenant/hooks";
import { useParams } from "react-router";

import { useListLlmDrivers } from "../hooks";

export function LlmFormFields() {
  const { tenantSlug = "" } = useParams<{ tenantSlug: string }>();
  const { data: tenant } = useFindBySlugTenant(tenantSlug);
  const { data: drivers } = useListLlmDrivers(tenant?.id ?? "");

  return (
    <FieldGroup>
      <FormInput name="name" label="Name" required />
      <DriverSelect
        name="provider"
        configName="configuration"
        drivers={drivers ?? []}
        label="Provider"
        placeholder="Select a provider"
      />
    </FieldGroup>
  );
}
