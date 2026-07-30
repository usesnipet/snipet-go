import { FormInput } from "@/components/form/input";
import { FormTextarea } from "@/components/form/textarea";
import { FieldGroup, FieldLegend, FieldSeparator, FieldSet } from "@/components/ui/field";

import { AgentLlmList } from "./agent-llm-list";

export function AgentFormFields() {
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
        <FieldLegend variant="label">Language models</FieldLegend>
        <AgentLlmList />
      </FieldSet>
    </FieldGroup>
  );
}
