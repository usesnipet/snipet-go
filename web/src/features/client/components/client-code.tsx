import { Button } from "@/components/ui/button";
import { useClipboard } from "@/hooks/use-clipboard";
import { CopyIcon } from "lucide-react";

export function ClientCode({ code }: { code: string }) {
  const { copy } = useClipboard();

  const handleCopy = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    event.stopPropagation();
    copy(code, {
      successTitle: `Client code copied to clipboard`,
      successDescription: "The code has been copied to your clipboard.",
      errorTitle: "Failed to copy to clipboard",
      errorDescription: "Please try again.",
    });
  }

  return (
    <div className="flex flex-row items-center gap-2">
      <span className="text-sm text-muted-foreground">{code}</span>
      <Button variant="outline" size="icon-xs" aria-label="Copy code" onClick={handleCopy}>
        <CopyIcon className="size-3" />
      </Button>
    </div>
  )
}