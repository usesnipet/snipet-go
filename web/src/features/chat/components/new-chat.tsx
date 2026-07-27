import { ChatInput } from "./chat-input";

export function NewChat() {
  return (
    <div className="flex flex-col gap-4 p-4 h-full items-center justify-center">
      <h1 className="text-2xl font-bold">New Chat</h1>
      <p className="text-sm text-muted-foreground">Start a new chat with me</p>
      <ChatInput containerclassname="w-full max-w-2xl" placeholder="Ask me anything..." />
    </div>
  )
}