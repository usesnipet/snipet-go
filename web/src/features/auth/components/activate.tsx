import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Link } from "@/components/ui/link";
import { Loading } from "@/components/ui/loading";
import { cn } from "@/lib/utils";
import { ROUTES } from "@/routes";
import { CheckCircle2, XCircle } from "lucide-react";
import { useSearchParams } from "react-router";

import { useActivate } from "../hooks";

export function AuthActivateForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token") ?? "";

  const { isLoading, isSuccess, isError } = useActivate(token);

  return (
    <Card className={cn("flex flex-col gap-6", className)} {...props}>
      <CardHeader>
        <CardTitle className="text-center">Activate your account</CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading && <Loading />}
        {isSuccess && (
          <>
            <p className="flex items-center justify-center gap-2 text-center">
              <CheckCircle2 className="size-4" />
              Your account has been activated. You can now sign in.
            </p>
            <Button asChild>
              <Link href={ROUTES.authLogin}>Sign in</Link>
            </Button>
          </>
        )}
        {(isError || !token) && (
          <div className="flex flex-col items-center justify-center gap-6 text-center">
            <XCircle className="size-12 text-destructive" />
            <p className="text-sm text-muted-foreground">
              This activation link is invalid or has expired. Try to login with your email and password to resend the activation email.
            </p>
            <Button asChild>
              <Link href={ROUTES.authLogin}>Sign in</Link>
            </Button>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
