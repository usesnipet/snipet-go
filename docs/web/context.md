# `src/context/`

Plain React Context providers for app-wide state that's read almost
everywhere, changes rarely, and doesn't need Zustand's ability to be read
outside a component. Today this is just `ThemeProvider`.

```tsx
type Theme = "dark" | "light";
const ThemeContext = React.createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = React.useState<Theme>(/* read localStorage / prefers-color-scheme */);
  // sync <html> class + localStorage on change
  return <ThemeContext.Provider value={{ theme, toggleTheme, setTheme }}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  const context = React.useContext(ThemeContext);
  if (context === undefined) throw new Error("useTheme must be used within a ThemeProvider");
  return context;
}
```

## Conventions

- Pair every context with a `use<Name>()` hook that throws if called
  outside its provider — never export the raw `Context` object or make
  callers use `useContext` directly.
- Mount the provider exactly once, in `root-providers.tsx` (see
  [README.md](./README.md#entry-point)) — providers here are app-wide by
  design, not something a feature or route mounts itself.

## Context vs. Zustand store

This app defaults to a Zustand `store.ts` for shared state (see
[features/store.md](./features/store.md)), and only reaches for
`src/context/` when **both** are true:

- the state is genuinely app-wide config, not one feature's data, and
- it only ever needs to be read from inside the React tree (unlike auth
  tokens or the dialog stack, which `lib/http`/`lib/dialog` need to read
  or mutate from outside React — see [lib.md](./lib.md)).

If either condition doesn't hold — the state belongs to one feature, or
something outside React needs to read/write it — use a Zustand store
instead of adding a new file here.
