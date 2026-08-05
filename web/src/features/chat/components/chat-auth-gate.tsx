import { Button } from "@/components/ui/button";
import {
  Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle
} from "@/components/ui/card";
import { Link } from "@/components/ui/link";
import { applyPathParams, applySearchParams } from "@/lib/http";
import { ROUTES } from "@/routes";
import { LogInIcon, ShieldOffIcon } from "lucide-react";

type ChatAuthGateProps = {
  clientCode: string;
  allowAnonymous: boolean;
  redirectPath: string;
};

export function ChatAuthGate({
  clientCode,
  allowAnonymous,
  redirectPath,
}: ChatAuthGateProps) {
  if (!allowAnonymous) {
    return (
      <div className="flex min-h-svh items-center justify-center bg-background p-6">
        <Card className="w-full max-w-md">
          <CardHeader>
            <div className="mb-2 flex size-10 items-center justify-center rounded-md bg-muted">
              <ShieldOffIcon className="size-5 text-muted-foreground" />
            </div>
            <CardTitle className="text-xl">Anonymous login disabled</CardTitle>
            <CardDescription>
              This client does not allow anonymous login, so the embedded chat cannot be used.
            </CardDescription>
          </CardHeader>
        </Card>
      </div>
    );
  }

  const loginHref = applySearchParams(
    applyPathParams(ROUTES.clientChatLoginAnonymous, { clientCode }),
    { redirect: redirectPath },
  );

  return (
    <div className="flex min-h-svh items-center justify-center bg-background p-6">
      <Card className="w-full max-w-md">
        <CardHeader>
          <div className="mb-2 flex size-10 items-center justify-center rounded-md bg-muted">
            <LogInIcon className="size-5 text-muted-foreground" />
          </div>
          <CardTitle className="text-xl">Sign in required</CardTitle>
          <CardDescription>
            You need to sign in to access the chat.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Continue as an anonymous user to start chatting with this client.
          </p>
        </CardContent>
        <CardFooter>
          <Button asChild className="w-full">
            <Link href={loginHref}>Continue anonymously</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
