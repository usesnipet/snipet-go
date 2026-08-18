import { DriverSelect } from "@/components/form/driver-select";
import { FormInput } from "@/components/form/input";
import { FieldGroup } from "@/components/ui/field";
import { useTenantStore } from "@/features/tenant/store";

import { useListLlmDrivers } from "../hooks";

export function LlmFormFields() {
  const tenant = useTenantStore((state) => state.tenant);
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
