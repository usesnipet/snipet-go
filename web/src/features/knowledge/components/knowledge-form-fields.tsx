import { DriverSelect } from "@/components/form/driver-select";
import { FormInput } from "@/components/form/input";
import { FormTextarea } from "@/components/form/textarea";
import { FieldGroup } from "@/components/ui/field";
import { useTenantStore } from "@/features/tenant/store";

import { useListKnowledgeDrivers } from "../hooks";

export function KnowledgeFormFields() {
  const tenant = useTenantStore((state) => state.tenant);
  const { data: drivers } = useListKnowledgeDrivers(tenant?.id ?? "");

  return (
    <FieldGroup>
      <FormInput name="name" label="Name" required />
      <FormTextarea
        name="description"
        label="Description"
        placeholder="What this knowledge source contains"
        rows={2}
      />
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
