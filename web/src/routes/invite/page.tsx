import { useSearchParams } from "react-router";

import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AcceptInvitationCard } from "@/features/invitation/components/accept-invitation-card";

export function AcceptInvitePage() {
  const [searchParams] = useSearchParams();
  const token = searchParams.get("token");

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-6 bg-background p-6 md:p-10">
      {token ? (
        <AcceptInvitationCard token={token} />
      ) : (
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle>Invalid invitation link</CardTitle>
            <CardDescription>This invitation link is missing a token.</CardDescription>
          </CardHeader>
        </Card>
      )}
    </div>
  )
}
