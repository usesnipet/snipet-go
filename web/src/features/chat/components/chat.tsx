import { useParams } from "react-router";

import { ChatInput } from "./chat-input";

export function Chat() {
  const { clientCode: clientCodeParam } = useParams<{ clientCode: string }>();
  const clientCode = clientCodeParam ?? "";

  return (
    <div className="flex flex-col gap-4 p-4 h-full items-center">
      <div className="flex-1">
        messages
      </div>
      <ChatInput
        clientCode={clientCode}
        containerclassname="w-full max-w-2xl"
        placeholder="Ask me anything..."
      />
    </div>
  );
}
