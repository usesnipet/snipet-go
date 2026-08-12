import { FormInput } from "@/components/form/input";
import { Button } from "@/components/ui/button";
import { Field, FieldDescription, FieldGroup } from "@/components/ui/field";
import { Form } from "@/components/ui/form";
import { Link } from "@/components/ui/link";
import { cn } from "@/lib/utils";
import { ROUTES } from "@/routes";
import { zodResolver } from "@hookform/resolvers/zod";
import { GalleryVerticalEnd } from "lucide-react";
import { useForm } from "react-hook-form";

import { useRegister } from "../hooks";
import { registerSchema } from "../schemas";

import type { Register } from "../schemas";

export function AuthRegisterForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const form = useForm<Register>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: "", email: "", password: "" },
  });

  const { mutate, isPending } = useRegister();
  const onSubmit = form.handleSubmit((data) => {
    mutate(data);
  });

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
              <h1 className="text-xl font-bold">Create your account</h1>
              <FieldDescription>
                Register with email and password. You&apos;ll need to activate
                via email before signing in.
              </FieldDescription>
            </div>

            <FormInput
              type="text"
              name="name"
              label="Name"
              placeholder="Jane Doe"
              autoComplete="name"
              required
            />
            <FormInput
              type="email"
              name="email"
              label="Email"
              placeholder="you@example.com"
              autoComplete="email"
              required
            />
            <FormInput
              type="password"
              name="password"
              label="Password"
              placeholder="At least 8 characters"
              autoComplete="new-password"
              description="Minimum 8 characters"
              required
            />
            <Field>
              <Button type="submit" disabled={isPending}>
                {isPending ? "Creating account…" : "Register"}
              </Button>
            </Field>
            <FieldDescription className="text-center">
              Already have an account?{" "}
              <Link href={ROUTES.authLogin} className="underline underline-offset-2">
                Sign in
              </Link>
            </FieldDescription>
          </FieldGroup>
        </form>
      </Form>
    </div>
  );
}
