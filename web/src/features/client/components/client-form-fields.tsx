import { FormInput } from "@/components/form/input";
import { FormSwitch } from "@/components/form/switch";
import { FieldGroup, FieldLegend, FieldSeparator, FieldSet } from "@/components/ui/field";
import { useFormContext } from "react-hook-form";

import type { ClientFormValues } from "./client-form";

export function ClientFormFields() {
  const form = useFormContext<ClientFormValues>();
  const oidcEnabled = form.watch("config.oidc.enabled");
  const webhookEnabled = form.watch("config.webhook.enabled");

  return (
    <FieldGroup>
      <FormInput name="name" label="Name" placeholder="Acme Inc." required />

      <FieldSet>
        <FieldSeparator />
        <FieldLegend variant="label">Authentication</FieldLegend>

        <FormSwitch
          name="config.anonymous.enabled"
          label="Allow anonymous access"
          fieldclassname="justify-between"
        />

        <FormSwitch
          name="config.oidc.enabled"
          label="OIDC"
          fieldclassname="justify-between"
        />
        {oidcEnabled && (
          <>
            <FormInput
              name="config.oidc.issuer"
              label="Issuer"
              placeholder="https://issuer.example.com"
            />
            <FormInput
              name="config.oidc.audience"
              label="Audience"
              placeholder="https://api.example.com"
            />
          </>
        )}

        <FormSwitch
          name="config.webhook.enabled"
          label="Webhook"
          fieldclassname="justify-between"
        />
        {webhookEnabled && (
          <FormInput
            name="config.webhook.url"
            label="Webhook URL"
            placeholder="https://example.com/webhook"
          />
        )}
      </FieldSet>
    </FieldGroup>
  );
}
