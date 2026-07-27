import { FormInput } from "@/components/form/input";
import { DriverSelect } from "@/components/form/driver-select";
import { FormTextarea } from "@/components/form/textarea";
import { FieldGroup, FieldLegend, FieldSeparator, FieldSet } from "@/components/ui/field";

import { useListLLMDrivers } from "../hooks";

export function AgentFormFields() {
  const { data: llmDrivers } = useListLLMDrivers();

  return (
    <FieldGroup>
      <FormInput name="name" label="Name" placeholder="Support assistant" required />
      <FormTextarea
        name="description"
        label="Description"
        placeholder="What this agent is for"
        rows={2}
      />
      <FormTextarea
        name="instructions"
        label="Instructions"
        placeholder="You are a helpful assistant..."
        rows={4}
      />

      <FieldSet>
        <FieldSeparator />
        <FieldLegend variant="label">Language model</FieldLegend>

        <DriverSelect
          name="llms.0.key"
          configName="llms.0.config"
          drivers={llmDrivers ?? []}
          label="Provider"
          placeholder="Select a provider"
        />
      </FieldSet>
    </FieldGroup>
  );
}
