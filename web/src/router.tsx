import { lazy, Suspense } from "react";
import { BrowserRouter, Route, Routes } from "react-router";

import { LoadingFallback } from "./components/loading-fallback";
import { ROUTES } from "./routes";

import type { RoutePath } from "./routes";

const Layout = lazy(() =>
  import("./routes/layout").then((m) => ({ default: m.Layout })));
const HomePage = lazy(() =>
  import("./routes/page").then((m) => ({ default: m.HomePage })));
const AgentsPage = lazy(() =>
  import("./routes/agents/page").then((m) => ({ default: m.AgentsPage })));
const AgentPlaygroundPage = lazy(() =>
  import("./routes/agents/playground/page").then((m) => ({ default: m.AgentPlaygroundPage })));
const ApiKeysPage = lazy(() =>
  import("./routes/api-keys/page").then((m) => ({ default: m.ApiKeysPage })));
const AppsPage = lazy(() =>
  import("./routes/apps/page").then((m) => ({ default: m.AppsPage })));
const KnowledgePage = lazy(() =>
  import("./routes/knowledge/page").then((m) => ({ default: m.KnowledgePage })));
const KnowledgeDetailPage = lazy(() =>
  import("./routes/knowledge/detail/page").then((m) => ({ default: m.KnowledgeDetailPage })));
const LLMPage = lazy(() =>
  import("./routes/llms/page").then((m) => ({ default: m.LLMPage })));

const AppAdminLayout = lazy(() =>
  import("./routes/app/(private)/layout").then((m) => ({ default: m.AppAdminLayout })));
const AppPage = lazy(() =>
  import("./routes/app/(private)/page").then((m) => ({ default: m.AppPage })));
const AppSessionsPage = lazy(() =>
  import("./routes/app/(private)/sessions/page").then((m) => ({ default: m.AppSessionsPage })));
const AppSettingsPage = lazy(() =>
  import("./routes/app/(private)/settings/page").then((m) => ({ default: m.AppSettingsPage })));
const AppUsersPage = lazy(() =>
  import("./routes/app/(private)/users/page").then((m) => ({ default: m.AppUsersPage })));

const toReactRouterPath = (path: RoutePath) => {
  return path.replaceAll(/{([^}]+)}/g, (_, p1) => `:${p1}`);
}

export const Router = () => {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingFallback className="min-h-svh" />}>
        <Routes>
          <Route element={<Layout />}>
            <Route path={toReactRouterPath(ROUTES.home)} element={<HomePage />} />
            <Route path={toReactRouterPath(ROUTES.apps)} element={<AppsPage />} />
            <Route path={toReactRouterPath(ROUTES.agent)} element={<AgentsPage />} />
            <Route path={toReactRouterPath(ROUTES.agentPlayground)} element={<AgentPlaygroundPage />} />
            <Route path={toReactRouterPath(ROUTES.knowledge)} element={<KnowledgePage />} />
            <Route path={toReactRouterPath(ROUTES.knowledgeDetail)} element={<KnowledgeDetailPage />} />
            <Route path={toReactRouterPath(ROUTES.llms)} element={<LLMPage />} />
            <Route path={toReactRouterPath(ROUTES.apiKeys)} element={<ApiKeysPage />} />
          </Route>
          <Route element={<AppAdminLayout />}>
            <Route path={toReactRouterPath(ROUTES.app)} element={<AppPage />} />
            <Route path={toReactRouterPath(ROUTES.appUsers)} element={<AppUsersPage />} />
            <Route path={toReactRouterPath(ROUTES.appSessions)} element={<AppSessionsPage />} />
            <Route path={toReactRouterPath(ROUTES.appSettings)} element={<AppSettingsPage />} />
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
