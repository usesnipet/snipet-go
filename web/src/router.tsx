import { lazy, Suspense } from "react";
import { BrowserRouter, Route, Routes } from "react-router";

import { LoadingFallback } from "./components/loading-fallback";
import { ApiKeyGuard } from "./features/api-key/components/api-key-guard";
import { ChatGuard } from "./features/chat/components/chat-guard";
import { ROUTES } from "./routes";

import type { RoutePath } from "./routes";

const AdminLayout = lazy(() =>
  import("./routes/admin/layout").then((m) => ({ default: m.AdminLayout })));
const AdminPage = lazy(() =>
  import("./routes/admin/page").then((m) => ({ default: m.AdminPage })));
const AdminAgentsPage = lazy(() =>
  import("./routes/admin/agents/page").then((m) => ({ default: m.AdminAgentsPage })));
const AdminApiKeysPage = lazy(() =>
  import("./routes/admin/api-keys/page").then((m) => ({ default: m.AdminApiKeysPage })));
const AdminClientsPage = lazy(() =>
  import("./routes/admin/clients/page").then((m) => ({ default: m.AdminClientsPage })));
const AdminKnowledgePage = lazy(() =>
  import("./routes/admin/knowledge/page").then((m) => ({ default: m.AdminKnowledgePage })));
const AdminKnowledgeDetailPage = lazy(() =>
  import("./routes/admin/knowledge/detail/page").then((m) => ({ default: m.AdminKnowledgeDetailPage })));
const AdminLlmsPage = lazy(() =>
  import("./routes/admin/llms/page").then((m) => ({ default: m.AdminLlmsPage })));

const ApiKeyLoginPage = lazy(() =>
  import("./routes/auth/api-key/page").then((m) => ({ default: m.ApiKeyLoginPage })));

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

const ClientChatLayout = lazy(() =>
  import("./routes/client/chat/layout").then((m) => ({ default: m.ClientChatLayout })));
const ClientChatPage = lazy(() =>
  import("./routes/client/chat/page").then((m) => ({ default: m.ClientChatPage })));
const ClientChatSessionPage = lazy(() =>
  import("./routes/client/chat/session/{sessionId}/page").then((m) => ({ default: m.ClientChatSessionPage })));
const ClientChatLoginAnonymousPage = lazy(() =>
  import("./routes/client/chat/login-anonymous/page").then((m) => ({ default: m.ClientChatLoginAnonymousPage })));

const HomePage = lazy(() =>
  import("./routes/page").then((m) => ({ default: m.HomePage })));

const toReactRouterPath = (path: RoutePath) => {
  return path.replaceAll(/{([^}]+)}/g, (_, p1) => `:${p1}`);
}

export const Router = () => {
  return (
    <BrowserRouter>
      <Suspense fallback={<LoadingFallback />}>
        <Routes>
          <Route path={toReactRouterPath(ROUTES.home)} element={<HomePage />} />
          <Route
            path={toReactRouterPath(ROUTES.clientChatLoginAnonymous)}
            element={<ClientChatLoginAnonymousPage />}
          />
          <Route element={<ChatGuard />}>
            <Route element={<ClientChatLayout />}>
              <Route path={toReactRouterPath(ROUTES.clientChat)} element={<ClientChatPage />} />
              <Route path={toReactRouterPath(ROUTES.clientChatSession)} element={<ClientChatSessionPage />} />
            </Route>
          </Route>
          <Route element={<ApiKeyGuard />}>
            <Route element={<ClientAdminLayout />}>
              <Route path={toReactRouterPath(ROUTES.client)} element={<ClientPage />} />
              <Route path={toReactRouterPath(ROUTES.clientUsers)} element={<ClientUsersPage />} />
              <Route path={toReactRouterPath(ROUTES.clientSessions)} element={<ClientSessionsPage />} />
              <Route path={toReactRouterPath(ROUTES.clientSettings)} element={<ClientSettingsPage />} />
            </Route>
            <Route element={<AdminLayout />}>
              <Route path={toReactRouterPath(ROUTES.admin)} element={<AdminPage />} />
              <Route path={toReactRouterPath(ROUTES.adminClients)} element={<AdminClientsPage />} />
              <Route path={toReactRouterPath(ROUTES.adminAgent)} element={<AdminAgentsPage />} />
              <Route path={toReactRouterPath(ROUTES.adminKnowledge)} element={<AdminKnowledgePage />} />
              <Route
                path={toReactRouterPath(ROUTES.adminKnowledgeDetail)}
                element={<AdminKnowledgeDetailPage />}
              />
              <Route path={toReactRouterPath(ROUTES.adminLlms)} element={<AdminLlmsPage />} />
              <Route path={toReactRouterPath(ROUTES.adminApiKeys)} element={<AdminApiKeysPage />} />
            </Route>
          </Route>
          <Route path={toReactRouterPath(ROUTES.authApiKey)} element={<ApiKeyLoginPage />} />
        </Routes>
      </Suspense>
    </BrowserRouter>
  )
}
