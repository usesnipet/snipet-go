# Routes

`src/routes.ts`, `src/router.tsx` and `src/routes/` together define every
page in the app. Each has a distinct job:

- **`routes.ts`** — the single source of truth for URL *strings*.
- **`router.tsx`** — wires those URLs to page components.
- **`routes/`** — the page and layout components themselves.

## `routes.ts`

A flat `ROUTES` const, one entry per page, grouped by area with a comment
(`// auth routes`, `// admin routes`, `// client routes`). Dynamic segments
use `{paramName}`, not react-router's `:paramName`:

```typescript
adminKnowledgeDetail: "/admin/knowledge/{id}",
client: "/client/{clientCode}",
clientChatSession: "/client/{clientCode}/chat/session/{sessionId}",
```

Anything that builds a URL — links, redirects, `service.ts` HTTP calls,
`router.tsx` — reads from `ROUTES`, never hardcodes a path string. This
keeps every path renameable from one place, and lets `service.ts` reuse the
same `{param}` placeholder syntax as `lib/http`'s `applyPathParams` (see
[lib.md](./lib.md)).

## `router.tsx`

Registers the `<Routes>` tree and translates `{param}` → `:param` for
react-router via `toReactRouterPath`. Structure:

- Every route **page and layout component must be lazy-loaded**:
  ```typescript
  const AdminLlmsPage = lazy(() =>
    import("./routes/admin/llms/page").then((m) => ({ default: m.AdminLlmsPage })));
  ```
  This is enforced project-wide — see [`web/CLAUDE.md`](../../web/CLAUDE.md)
  for the exact rule. Guards and shared chrome (`LoadingFallback`, `ROUTES`)
  stay as regular static imports.
- The whole tree is wrapped once in `<Suspense fallback={<LoadingFallback />}>`.
- Route **guards** wrap the routes they protect:
  ```tsx
  <Route element={<ApiKeyGuard />}>
    <Route element={<AdminLayout />}>
      <Route path={toReactRouterPath(ROUTES.admin)} element={<AdminPage />} />
      ...
    </Route>
  </Route>
  ```
  A guard is a component that renders `<Outlet />` when access is allowed,
  or `<Navigate />` otherwise. Guards belong to the feature that owns the
  auth concern (`features/api-key/components/api-key-guard.tsx`,
  `features/chat/components/chat-guard.tsx`), not to `routes/`.
- **Layouts** wrap a group of pages with shared chrome (sidebar, etc.) and
  render `<AnimatedOutlet />` where the matched page goes.

## `routes/` folder

The folder tree under `src/routes/` mirrors the URL tree. Each URL segment
is a folder; each folder has at most one `page.tsx` (the routed screen) and
optionally a `layout.tsx` (wraps its own children with shared chrome via
`<Outlet />`/`<AnimatedOutlet />`).

```
routes/
  admin/
    layout.tsx          # AdminLayout — sidebar + outlet for every /admin/* page
    page.tsx             # AdminPage — "/admin"
    llms/
      page.tsx            # AdminLlmsPage — "/admin/llms"
    knowledge/
      page.tsx
      detail/
        page.tsx           # "/admin/knowledge/{id}"
  client/
    (private)/
      layout.tsx          # ClientAdminLayout
      page.tsx
      users/page.tsx
    chat/
      layout.tsx
      page.tsx
      session/{sessionId}/page.tsx
```

Two folder-naming conventions to know:

- **`{paramName}/`** — a dynamic segment, matching the `{paramName}`
  placeholder used for the same route in `ROUTES`
  (`session/{sessionId}/page.tsx` ↔ `clientChatSession`).
- **`(groupName)/`** — a route group. Parentheses mean the segment is
  *not* part of the URL; it exists only to nest several pages under one
  shared `layout.tsx` without adding a path segment (`client/(private)/`
  groups the client-admin pages under `ClientAdminLayout`).

Exports from `page.tsx`/`layout.tsx` are **named**, not default (`export
function AdminLlmsPage()`), because `router.tsx`'s lazy loader maps the
named export to `default` itself.

## What a page does

A page component is thin: it composes feature hooks/components and the
shared page chrome from `src/components/page.tsx` (`<Page>`,
`<PageActions>`) for a consistent title/description/document-title header:

```tsx
export function AdminLlmsPage() {
  const { openDialog } = useDialog();
  return (
    <Page title="LLMs" description="..." documentTitle="LLMs">
      <PageActions>
        <Button onClick={() => openDialog({ component: CreateLlmDialog, props: {} })}>
          Create LLM
        </Button>
      </PageActions>
      <LlmTable />
    </Page>
  );
}
```

All actual data-fetching, mutations and domain UI live in the owning
feature (`src/features/<feature>/`) — see
[features/README.md](./features/README.md). Pages don't call `service.ts`
or hold query/mutation state directly.

## Adding a new route

1. Add the URL to `ROUTES` in `routes.ts`.
2. Create `routes/<path>/page.tsx` (and `layout.tsx` if it needs shared
   chrome not already provided by a parent layout).
3. Register it in `router.tsx` as a lazy import, inside the right
   guard/layout nesting.
