# `hooks.ts`

Wraps `service.ts` with TanStack Query's `useQuery`/`useMutation`. This is
the layer that owns caching, loading/error state, and the auth mode for
each call — components only ever import from here, never from `service.ts`
directly.

```typescript
const BASE_QUERY_KEY = "llm";

export const listLlmQueryKey = () => [BASE_QUERY_KEY] as const;
export const useListLlm = (
  opts?: ServiceGetOptions<PaginatedLlm, ListLlmSearchParams>,
): UseQueryResult<PaginatedLlm, Error> => {
  return useQuery({
    queryKey: [...listLlmQueryKey(), opts?.searchParams],
    queryFn: () => llmService.list(opts),
  });
};

export const createLlmQueryKey = () => [BASE_QUERY_KEY, "create"] as const;
export const useCreateLlm = (
  opts?: ServicePostOptions<CreateLlm, Llm>,
): UseMutationResult<Llm, Error, CreateLlm> => {
  return useMutation({
    mutationKey: createLlmQueryKey(),
    mutationFn: (data: CreateLlm) => llmService.create(data, opts),
    onSuccess: () => {
      toast({ title: "LLM created successfully", description: "..." });
      queryClient.invalidateQueries({ queryKey: listLlmQueryKey() });
    },
    onError: () => {
      toast({ title: "Failed to create LLM", description: "...", variant: "destructive" });
    },
  });
};
```

## Conventions

- **Export a query-key factory next to every hook**: `listLlmQueryKey()`,
  `createLlmQueryKey()`, etc. Anything that needs to target this query's
  cache entry (another hook's `invalidateQueries`, a manual prefetch)
  imports the factory instead of hand-writing the key array. All keys for
  a feature start with the same `BASE_QUERY_KEY`.
- **One hook per service method**, named `use<Verb><Feature>`
  (`useListLlm`, `useCreateLlm`, `useUpdateLlm`, `useDeleteLlm`).
- **Auth mode is set here, not in `service.ts`**: `{ ...opts, auth:
  "api-key" }` (or `"jwt"`, depending on which surface of the app the
  feature belongs to). This keeps `service.ts` agnostic of *who* is
  calling it — a service function can be reused by hooks with different
  auth requirements.
- **Mutations toast and invalidate**: on `onSuccess`, show a success toast
  (`@/hooks/use-toast`) and `queryClient.invalidateQueries` for whichever
  query keys the mutation affects (usually the feature's `list` key). On
  `onError`, show a `variant: "destructive"` toast. Queries (`useQuery`)
  don't toast — let the caller render `isLoading`/`error` from the result.
- **Forward, don't swallow, caller options**: every hook takes an optional
  `opts` of the matching `Service*Options` type and spreads it into the
  `service` call, so a caller can still override `searchParams`,
  `headers`, etc.

## What doesn't belong here

No JSX. Hooks return query/mutation results for a component to consume —
they don't render anything themselves.
