import { lazy, Suspense } from "react";
import { BrowserRouter, Route, Routes } from "react-router";

import { LoadingFallback } from "./components/loading-fallback";
import { AuthGuard } from "./features/auth/components/auth-guard";
import { ROUTES } from "./routes";

import type { RoutePath } from "./routes";

const Layout = lazy(() =>
  import("./routes/layout").then((m) => ({ default: m.Layout })));
const HomePage = lazy(() =>
  import("./routes/page").then((m) => ({ default: m.HomePage })));
const AgentsPage = lazy(() =>
  import("./routes/agents/page").then((m) => ({ default: m.AgentsPage })));
const ApiKeysPage = lazy(() =>
  import("./routes/api-keys/page").then((m) => ({ default: m.ApiKeysPage })));
const ClientsPage = lazy(() =>
  import("./routes/clients/page").then((m) => ({ default: m.ClientsPage })));
const KnowledgePage = lazy(() =>
  import("./routes/knowledge/page").then((m) => ({ default: m.KnowledgePage })));
const KnowledgeDetailPage = lazy(() =>
  import("./routes/knowledge/detail/page").then((m) => ({ default: m.KnowledgeDetailPage })));
const LLMPage = lazy(() =>
  import("./routes/llms/page").then((m) => ({ default: m.LLMPage })));

const ClientAdminLayout = lazy(() =>
  import("./routes/client/(private)/layout").then((m) => ({ default: m.ClientAdminLayout })));
const ClientPage = lazy(() =>
  import("./routes/client/(private)/page").then((m) => ({ default: m.ClientPage })));
const ClientSessionsPage = lazy(() =>
  import("./routes/client/(private)/sessions/page").then((m) => ({ default: m.ClientSessionsPage })));
const ClientSettingsPage = lazy(() =>
  import("./routes/client/(private)/settings/page").then((m) => ({ default: m.ClientSettingsPage })));
const ClientUsersPage = lazy(() =>
  import("./routes/client/(private)/users/page").then((m) => ({ default: m.ClientUsersPage })));

const AuthLoginPage = lazy(() =>
  import("./routes/auth/login/page").then((m) => ({ default: m.AuthLoginPage })));
const AuthRegisterPage = lazy(() =>
  import("./routes/auth/register/page").then((m) => ({ default: m.AuthRegisterPage })));

const toReactRouterPath = (path: RoutePath) => {
  return path.replaceAll(/{([^}]+)}/g, (_, p1) => `:${p1}`);
}

export const Router = () => {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingFallback />}>
        <Routes>
          <Route path={toReactRouterPath(ROUTES.authLogin)} element={<AuthLoginPage />} />
          <Route path={toReactRouterPath(ROUTES.authRegister)} element={<AuthRegisterPage />} />
          <Route element={<AuthGuard />}>
            <Route element={<ClientAdminLayout />}>
              <Route path={toReactRouterPath(ROUTES.client)} element={<ClientPage />} />
              <Route path={toReactRouterPath(ROUTES.clientUsers)} element={<ClientUsersPage />} />
              <Route path={toReactRouterPath(ROUTES.clientSessions)} element={<ClientSessionsPage />} />
              <Route path={toReactRouterPath(ROUTES.clientSettings)} element={<ClientSettingsPage />} />
            </Route>
            <Route element={<Layout />}>
              <Route path={toReactRouterPath(ROUTES.home)} element={<HomePage />} />
              <Route path={toReactRouterPath(ROUTES.clients)} element={<ClientsPage />} />
              <Route path={toReactRouterPath(ROUTES.agent)} element={<AgentsPage />} />
              <Route path={toReactRouterPath(ROUTES.knowledge)} element={<KnowledgePage />} />
              <Route
                path={toReactRouterPath(ROUTES.knowledgeDetail)}
                element={<KnowledgeDetailPage />}
              />
              <Route path={toReactRouterPath(ROUTES.llms)} element={<LLMPage />} />
              <Route path={toReactRouterPath(ROUTES.apiKeys)} element={<ApiKeysPage />} />
            </Route>
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
