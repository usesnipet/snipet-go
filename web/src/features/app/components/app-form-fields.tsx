import { FormInput } from "@/components/form/input";
import { FormTextarea } from "@/components/form/textarea";
import { FieldGroup } from "@/components/ui/field";

export function AppFormFields() {
  return (
    <FieldGroup>
      <FormInput name="name" label="Name" placeholder="Acme Inc." required />
      <FormTextarea name="description" label="Description" placeholder="What this app is for" />
    </FieldGroup>
  );
}
