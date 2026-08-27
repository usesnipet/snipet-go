import { lazy, Suspense } from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router";

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
const MembersPage = lazy(() =>
  import("./routes/members/page").then((m) => ({ default: m.MembersPage })));
const InviteMembersPage = lazy(() =>
  import("./routes/members/invite/page").then((m) => ({ default: m.InviteMembersPage })));

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

const AuthLoginPage = lazy(() =>
  import("./routes/auth/login/page").then((m) => ({ default: m.AuthLoginPage })));
const AuthRegisterPage = lazy(() =>
  import("./routes/auth/register/page").then((m) => ({ default: m.AuthRegisterPage })));
const AuthActivatePage = lazy(() =>
  import("./routes/auth/activate/page").then((m) => ({ default: m.AuthActivatePage })));

const TenantSelectPage = lazy(() =>
  import("./routes/tenant-select/page").then((m) => ({ default: m.TenantSelectPage })));
const AcceptInvitePage = lazy(() =>
  import("./routes/invite/page").then((m) => ({ default: m.AcceptInvitePage })));

const toReactRouterPath = (path: RoutePath) => {
  return path.replaceAll(/{([^}]+)}/g, (_, p1) => `:${p1}`);
}

console.log(toReactRouterPath(ROUTES.knowledge));

export const Router = () => {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingFallback className="min-h-svh" />}>
        <Routes>
          <Route path={toReactRouterPath(ROUTES.home)} element={<Navigate to={toReactRouterPath(ROUTES.selectTenant)} replace />} />
          <Route element={<AuthGuard mode="public" />}>
            <Route path={toReactRouterPath(ROUTES.authLogin)} element={<AuthLoginPage />} />
            <Route path={toReactRouterPath(ROUTES.authRegister)} element={<AuthRegisterPage />} />
            <Route path={toReactRouterPath(ROUTES.authActivate)} element={<AuthActivatePage />} />
          </Route>
          <Route element={<AuthGuard />}>
            <Route path={toReactRouterPath(ROUTES.selectTenant)} element={<TenantSelectPage />} />
            <Route path={toReactRouterPath(ROUTES.acceptInvite)} element={<AcceptInvitePage />} />
            <Route element={<Layout />}>
              <Route path={toReactRouterPath(ROUTES.tenantHome)} element={<HomePage />} />
              <Route path={toReactRouterPath(ROUTES.apps)} element={<AppsPage />} />
              <Route path={toReactRouterPath(ROUTES.agent)} element={<AgentsPage />} />
              <Route path={toReactRouterPath(ROUTES.agentPlayground)} element={<AgentPlaygroundPage />} />
              <Route path={toReactRouterPath(ROUTES.knowledge)} element={<KnowledgePage />} />
              <Route path={toReactRouterPath(ROUTES.knowledgeDetail)} element={<KnowledgeDetailPage />}/>
              <Route path={toReactRouterPath(ROUTES.llms)} element={<LLMPage />} />
              <Route path={toReactRouterPath(ROUTES.members)} element={<MembersPage />} />
              <Route path={toReactRouterPath(ROUTES.inviteMember)} element={<InviteMembersPage />} />
              <Route path={toReactRouterPath(ROUTES.apiKeys)} element={<ApiKeysPage />} />
            </Route>
            <Route element={<AppAdminLayout />}>
              <Route path={toReactRouterPath(ROUTES.app)} element={<AppPage />} />
              <Route path={toReactRouterPath(ROUTES.appUsers)} element={<AppUsersPage />} />
              <Route path={toReactRouterPath(ROUTES.appSessions)} element={<AppSessionsPage />} />
              <Route path={toReactRouterPath(ROUTES.appSettings)} element={<AppSettingsPage />} />
            </Route>
          </Route>
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
