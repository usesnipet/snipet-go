import { Building2Icon, CheckCircle2, MailQuestion, XCircle } from "lucide-react";
import { useState } from "react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle
} from "@/components/ui/card";
import { Loading } from "@/components/ui/loading";
import { Spinner } from "@/components/ui/spinner";
import { useNavigate } from "@/hooks/use-navigate";
import { ROUTES } from "@/routes";

import { useAcceptInvitation, useDeclineInvitation, useGetInvitationByToken } from "../hooks";

type AcceptInvitationCardProps = {
  token: string
}

export function AcceptInvitationCard({ token }: AcceptInvitationCardProps) {
  const navigate = useNavigate();
  const [declined, setDeclined] = useState(false);

  const { data, isLoading, isError } = useGetInvitationByToken(token);
  const { mutateAsync: accept, isPending: accepting } = useAcceptInvitation();
  const { mutateAsync: decline, isPending: declining } = useDeclineInvitation();

  const isActionPending = accepting || declining;

  const handleAccept = async () => {
    await accept({ token });
    navigate(ROUTES.tenantHome, { params: { tenantSlug: data?.tenant.slug ?? "" } });
  };

  const handleDecline = async () => {
    await decline({ token });
    setDeclined(true);
  };

  if (isLoading) {
    return (
      <Card className="w-full max-w-md">
        <CardContent className="py-6">
          <Loading />
        </CardContent>
      </Card>
    )
  }

  if (isError || !data) {
    return (
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <XCircle className="size-5 text-muted-foreground" />
            Invitation not found
          </CardTitle>
          <CardDescription>
            This invitation link is invalid or no longer exists.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  const { invite, tenant } = data;
  const isExpired = invite.status === "pending" && invite.expires_at < new Date();

  if (declined || invite.status === "declined") {
    return (
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <XCircle className="size-5 text-muted-foreground" />
            Invitation declined
          </CardTitle>
          <CardDescription>
            You have declined the invitation to join{" "}
            <span className="font-medium text-foreground">{tenant.name}</span>.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  if (invite.status === "accepted") {
    return (
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <CheckCircle2 className="size-5 text-muted-foreground" />
            Invitation already accepted
          </CardTitle>
          <CardDescription>
            This invitation to join{" "}
            <span className="font-medium text-foreground">{tenant.name}</span>{" "}
            has already been accepted.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  if (isExpired) {
    return (
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <XCircle className="size-5 text-muted-foreground" />
            Invitation expired
          </CardTitle>
          <CardDescription>
            This invitation to join{" "}
            <span className="font-medium text-foreground">{tenant.name}</span>{" "}
            has expired. Ask an admin to send you a new one.
          </CardDescription>
        </CardHeader>
      </Card>
    )
  }

  return (
    <Card className="w-full max-w-md">
      <CardHeader>
        <div className="flex items-center gap-3">
          <Avatar size="lg">
            {tenant.icon && <AvatarImage src={tenant.icon} alt={tenant.name} />}
            <AvatarFallback>
              <Building2Icon className="size-4" />
            </AvatarFallback>
          </Avatar>
          <div className="min-w-0 flex-1">
            <CardTitle className="flex items-center gap-2 truncate">
              <MailQuestion className="size-4 shrink-0 text-muted-foreground" />
              {tenant.name}
            </CardTitle>
            <CardDescription>You have been invited to join this tenant.</CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-muted-foreground">
          Invited as{" "}
          <Badge variant={invite.role === "admin" ? "default" : "secondary"}>
            {invite.role === "admin" ? "Admin" : "Member"}
          </Badge>
        </p>
      </CardContent>
      <CardFooter className="flex justify-end gap-2">
        <Button type="button" variant="outline" disabled={isActionPending} onClick={handleDecline}>
          {declining ? <Spinner size="sm" /> : <XCircle />}
          Decline
        </Button>
        <Button type="button" disabled={isActionPending} onClick={handleAccept}>
          {accepting ? <Spinner size="sm" /> : <CheckCircle2 />}
          Accept
        </Button>
      </CardFooter>
    </Card>
  )
}
