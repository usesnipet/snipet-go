import { ErrorFallback } from "@/components/error-fallback";
import { LoadingFallback } from "@/components/loading-fallback";

import { useListSession } from "../hooks";

export const SessionList = ({ clientCode }: { clientCode: string }) => {
  const { data, isLoading, error } = useListSession(clientCode);
  if (isLoading) return <LoadingFallback />;
  if (error) return <ErrorFallback error={error} />;

  return (
    <div className="flex flex-col gap-2">
      {data.data.map((session) => (
        <div key={session.id}>
          <h3>{session.id}</h3>
        </div>
      ))}
    </div>
  )
}