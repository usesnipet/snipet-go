import { FormInput } from "@/components/form/input";
import { FormSwitch } from "@/components/form/switch";
import { FormTextarea } from "@/components/form/textarea";
import {
  FieldDescription, FieldGroup, FieldLegend, FieldSeparator, FieldSet
} from "@/components/ui/field";

import { AppAgentList } from "./app-agent-list";

export function AppFormFields() {
  return (
    <FieldGroup>
      <FormInput name="name" label="Name" placeholder="Acme Inc." required />
      <FormTextarea name="description" label="Description" placeholder="What this app is for" />
      <FormSwitch name="public" label="Public" fieldclassname="justify-between" />
      <FieldDescription>
        Public apps expose their name, code, and description on an unauthenticated endpoint, and can be manually activated without a key ping.
      </FieldDescription>

      <FieldSet>
        <FieldSeparator />
        <FieldLegend variant="label">Agents</FieldLegend>
        <AppAgentList />
        <FieldDescription>
          Linked agents are returned alongside the app on its public endpoint.
        </FieldDescription>
      </FieldSet>
    </FieldGroup>
  );
}
