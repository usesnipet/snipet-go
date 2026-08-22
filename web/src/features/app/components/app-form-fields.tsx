import { FormInput } from "@/components/form/input";
import { FormSwitch } from "@/components/form/switch";
import { FormTextarea } from "@/components/form/textarea";
import { FieldDescription, FieldGroup } from "@/components/ui/field";

export function AppFormFields() {
  return (
    <FieldGroup>
      <FormInput name="name" label="Name" placeholder="Acme Inc." required />
      <FormTextarea name="description" label="Description" placeholder="What this app is for" />
      <FormSwitch name="public" label="Public" fieldclassname="justify-between" />
      <FieldDescription>
        Public apps expose their name, code, and description on an unauthenticated endpoint, and can be manually activated without a key ping.
      </FieldDescription>
    </FieldGroup>
  );
}
