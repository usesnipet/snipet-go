import { FormInput } from "@/components/form/input";
import { FormSwitch } from "@/components/form/switch";
import { FieldLegend, FieldSeparator, FieldSet } from "@/components/ui/field";
import { useFormContext } from "react-hook-form";

import type { UpdateAppAuthConfig } from "../schemas";

export function AppAuthConfigFields() {
  const form = useFormContext<UpdateAppAuthConfig>();
  const oidcEnabled = form.watch("auth_config.oidc.enabled");
  const webhookEnabled = form.watch("auth_config.webhook.enabled");

  return (
    <FieldSet>
      <FieldSeparator />
      <FieldLegend variant="label">Authentication</FieldLegend>

      <FormSwitch
        name="auth_config.oidc.enabled"
        label="OIDC"
        fieldclassname="justify-between"
      />
      {oidcEnabled && (
        <>
          <FormInput
            name="auth_config.oidc.issuer"
            label="Issuer"
            placeholder="https://issuer.example.com"
          />
          <FormInput
            name="auth_config.oidc.audience"
            label="Audience"
            placeholder="https://api.example.com"
          />
        </>
      )}

      <FormSwitch
        name="auth_config.webhook.enabled"
        label="Webhook"
        fieldclassname="justify-between"
      />
      {webhookEnabled && (
        <FormInput
          name="auth_config.webhook.url"
          label="Webhook URL"
          placeholder="https://example.com/webhook"
        />
      )}
    </FieldSet>
  );
}
