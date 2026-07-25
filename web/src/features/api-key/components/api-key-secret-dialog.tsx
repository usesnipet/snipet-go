import { toast } from "@/hooks/use-toast";
import { Check, Copy } from "lucide-react";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";

type ApiKeySecretDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  secret: string | null
  title?: string
  description?: string
}

export function ApiKeySecretDialog({
  open,
  onOpenChange,
  secret,
  title = "API Key created",
  description = "Copy this key now. You will not be able to see it again.",
}: ApiKeySecretDialogProps) {
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    if (!secret) return;
    await navigator.clipboard.writeText(secret);
    setCopied(true);
    toast({ title: "Copied to clipboard" });
    window.setTimeout(() => setCopied(false), 2000);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md" showCloseButton={false}>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>
        <div className="flex gap-2">
          <Input readOnly value={secret ?? ""} className="font-mono text-xs" />
          <Button type="button" variant="outline" size="icon" onClick={handleCopy}>
            {copied ? <Check /> : <Copy />}
            <span className="sr-only">Copy</span>
          </Button>
        </div>
        <DialogFooter>
          <Button type="button" onClick={() => onOpenChange(false)}>
            Done
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
