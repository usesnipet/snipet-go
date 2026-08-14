import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useLogout } from "@/features/auth/hooks";
import { LogOutIcon, UserIcon } from "lucide-react";

import { useMeUser } from "../hooks";

export function UserCard() {
  const { data: user } = useMeUser();
  const { mutate: logout } = useLogout();

  const handleLogout = () => logout({});

  return (
    <Card className="border-none bg-muted/50 shadow-none group-data-[collapsible=icon]:bg-transparent">
      <CardContent className="flex items-center gap-2 p-2 group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:p-0">
        <Avatar size="sm">
          {user?.picture && <AvatarImage src={user.picture} alt={user.name} />}
          <AvatarFallback>
            <UserIcon className="size-4" />
          </AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1 group-data-[collapsible=icon]:hidden">
          <p className="truncate text-sm font-semibold leading-tight">{user?.name}</p>
          <p className="truncate text-xs text-muted-foreground">{user?.email}</p>
        </div>
        <Button
          variant="ghost"
          size="icon-sm"
          className="shrink-0 text-muted-foreground hover:text-destructive group-data-[collapsible=icon]:hidden"
          onClick={handleLogout}
          aria-label="Logout"
        >
          <LogOutIcon className="size-4" />
        </Button>
      </CardContent>
    </Card>
  );
}
