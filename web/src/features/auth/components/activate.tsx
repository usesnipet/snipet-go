import { CheckCircle2, GalleryVerticalEnd, XCircle } from "lucide-react";
import { useEffect, useRef } from "react";
import { useSearchParams } from "react-router";

import { Button } from "@/components/ui/button";
import { Field, FieldDescription, FieldGroup } from "@/components/ui/field";
import { Link } from "@/components/ui/link";
import { Loading } from "@/components/ui/loading";
import { cn } from "@/lib/utils";
import { ROUTES } from "@/routes";

import { useActivate } from "../hooks";

export function AuthActivateForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token") ?? "";

  const { mutate, isPending, isSuccess, isError } = useActivate();
  const triggered = useRef(false);

  useEffect(() => {
    if (!token || triggered.current) return;
    triggered.current = true;
    mutate({ token });
  }, [token, mutate]);

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <FieldGroup>
        <div className="flex flex-col items-center gap-2 text-center">
          <div className="flex size-8 items-center justify-center rounded-md">
            <GalleryVerticalEnd className="size-6" />
          </div>
          <h1 className="text-xl font-bold">Activate your account</h1>
        </div>

        {!token ? (
          <FieldDescription className="text-center">
            This activation link is missing a token. Check the link in your email or request a
            new one.
            <div className="mt-2">
              <Link
                href={ROUTES.authResendActivation}
                className="underline underline-offset-2"
              >
                Resend activation email
              </Link>
            </div>
          </FieldDescription>
        ) : isPending ? (
          <Loading />
        ) : isSuccess ? (
          <>
            <FieldDescription className="flex items-center justify-center gap-2 text-center">
              <CheckCircle2 className="size-4" />
              Your account has been activated. You can now sign in.
            </FieldDescription>
            <Field>
              <Button asChild>
                <Link href={ROUTES.authLogin}>Sign in</Link>
              </Button>
            </Field>
          </>
        ) : isError ? (
          <>
            <FieldDescription className="flex items-center justify-center gap-2 text-center text-destructive">
              <XCircle className="size-4" />
              This activation link is invalid or has expired.
            </FieldDescription>
            <FieldDescription className="text-center">
              <Link
                href={ROUTES.authResendActivation}
                className="underline underline-offset-2"
              >
                Resend activation email
              </Link>
            </FieldDescription>
          </>
        ) : null}
      </FieldGroup>
    </div>
  );
}
