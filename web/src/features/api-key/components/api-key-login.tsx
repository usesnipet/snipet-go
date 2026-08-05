import { FormInput } from "@/components/form/input";
import { Button } from "@/components/ui/button";
import { Field, FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { cn } from "@/lib/utils";
import { zodResolver } from "@hookform/resolvers/zod";
import { GalleryVerticalEnd } from "lucide-react";
import { useForm } from "react-hook-form";
import { useSearchParams } from "react-router";
import z from "zod";

import { useApiKeyLogin } from "../hooks";
import { apiKeyKeySchema } from "../schemas";

import type { RoutePath } from "@/routes";

export function ApiKeyLoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [searchParams] = useSearchParams();
  const redirect = searchParams.get("redirect") as RoutePath;
  const form = useForm({
    resolver: zodResolver(z.object({ apiKey: apiKeyKeySchema })),
    defaultValues: { apiKey: "" },
  })

  const { mutate } = useApiKeyLogin(redirect);
  const onSubmit = form.handleSubmit((data) => {
    mutate(data.apiKey);
  })

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Form {...form}>
        <form onSubmit={onSubmit}>
          <FieldGroup>
            <div className="flex flex-col items-center gap-2 text-center">
              <a
                href="#"
                className="flex flex-col items-center gap-2 font-medium"
              >
                <div className="flex size-8 items-center justify-center rounded-md">
                  <GalleryVerticalEnd className="size-6" />
                </div>
                <span className="sr-only">Snipet.</span>
              </a>
              <h1 className="text-xl font-bold">Welcome to Snipet.</h1>
            </div>

            <FormInput
              type="password"
              name="apiKey"
              label="API Key"
              placeholder="Enter your API Key"
              description="The API key is used to authenticate your requests to the API."
              required
            />
            <Field>
              <Button type="submit">Login</Button>
            </Field>
          </FieldGroup>
        </form>
      </Form>
    </div>
  )
}
