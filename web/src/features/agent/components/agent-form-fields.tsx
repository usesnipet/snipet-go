import { FormInput } from "@/components/form/input";
import { FormSelect } from "@/components/form/select";
import { FormSlider } from "@/components/form/slider";
import { FormTextarea } from "@/components/form/textarea";
import { FieldGroup, FieldLegend, FieldSeparator, FieldSet } from "@/components/ui/field";

import { useListLLMDrivers } from "../hooks";

export function AgentFormFields() {
  const { data: llmDrivers } = useListLLMDrivers();
  const LLM_OPTIONS = (llmDrivers ?? []).map((driver) => ({
    label: driver.name,
    value: driver.name,
  }));

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

        <FormSelect
          name="llm_key"
          label="Provider"
          options={LLM_OPTIONS}
          placeholder="Select a provider"
        />
        <FormInput
          name="llm_config.api_key"
          label="API key"
          placeholder="AIza..."
          type="password"
          autoComplete="off"
          required
        />
        <FormInput
          name="llm_config.model"
          label="Model"
          placeholder="gemini-2.0-flash"
          required
        />
        <FormSlider
          name="llm_config.temperature"
          label="Temperature"
          min={0}
          max={2}
          step={0.1}
          showValue
        />
        <FormSlider
          name="llm_config.top_p"
          label="Top P"
          min={0}
          max={1}
          step={0.05}
          showValue
        />
      </FieldSet>
    </FieldGroup>
  );
}
