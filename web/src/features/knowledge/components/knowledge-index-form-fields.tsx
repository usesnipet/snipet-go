import { DriverSelect } from "@/components/form/driver-select";
import { FormInput } from "@/components/form/input";
import { FieldGroup } from "@/components/ui/field";

import { useListKnowledgeIndexDrivers } from "../hooks";

export function KnowledgeIndexFormFields() {
  const { data: drivers } = useListKnowledgeIndexDrivers();

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
