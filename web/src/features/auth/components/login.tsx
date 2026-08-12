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
import { useSearchParams } from "react-router";

import { useLogin } from "../hooks";
import { loginSchema } from "../schemas";

import type { Login } from "../schemas";
import type { RoutePath } from "@/routes";

export function AuthLoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [searchParams] = useSearchParams();
  const redirect = searchParams.get("redirect") as RoutePath;
  const form = useForm<Login>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: "", password: "" },
  });

  const { mutate, isPending } = useLogin(redirect);
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
              <h1 className="text-xl font-bold">Welcome to Snipet.</h1>
              <FieldDescription>
                Sign in with your email and password
              </FieldDescription>
            </div>

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
              placeholder="Enter your password"
              autoComplete="current-password"
              required
            />
            <Field>
              <Button type="submit" disabled={isPending}>
                {isPending ? "Signing in…" : "Sign in"}
              </Button>
            </Field>
            <FieldDescription className="text-center">
              Don&apos;t have an account?{" "}
              <Link href={ROUTES.authRegister} className="underline underline-offset-2">
                Register
              </Link>
              {" · "}
              <Link
                href={ROUTES.authForgotPassword}
                className="underline underline-offset-2"
              >
                Forgot password?
              </Link>
            </FieldDescription>
          </FieldGroup>
        </form>
      </Form>
    </div>
  );
}
