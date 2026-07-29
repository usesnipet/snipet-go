import { BrowserRouter, Route, Routes } from "react-router";

import { ApiKeyGuard } from "./features/api-key/components/api-key-guard";
import { ChatGuard } from "./features/chat/components/chat-guard";
import { ROUTES } from "./routes";
import { AdminAgentsPage } from "./routes/admin/agents/page";
import { AdminApiKeysPage } from "./routes/admin/api-keys/page";
import { AdminClientsPage } from "./routes/admin/clients/page";
import { AdminLayout } from "./routes/admin/layout";
import { AdminLlmsPage } from "./routes/admin/llms/page";
import { AdminPage } from "./routes/admin/page";
import { ApiKeyLoginPage } from "./routes/auth/api-key/page";
import { ClientAdminLayout } from "./routes/client/(private)/layout";
import { ClientPage } from "./routes/client/(private)/page";
import { ClientSessionsPage } from "./routes/client/(private)/sessions/page";
import { ClientSettingsPage } from "./routes/client/(private)/settings/page";
import { ClientUsersPage } from "./routes/client/(private)/users/page";
import { ClientChatLayout } from "./routes/client/chat/layout";
import { ClientChatLoginAnonymousPage } from "./routes/client/chat/login-anonymous/page";
import { ClientChatPage } from "./routes/client/chat/page";
import { ClientChatSessionPage } from "./routes/client/chat/session/{sessionId}/page";
import { HomePage } from "./routes/page";

import type { RoutePath } from "./routes";
const toReactRouterPath = (path: RoutePath) => {
  return path.replaceAll(/{([^}]+)}/g, (_, p1) => `:${p1}`);
}

export const Router = () => {
  return (
    <BrowserRouter>
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
            <Route path={toReactRouterPath(ROUTES.adminLlms)} element={<AdminLlmsPage />} />
            <Route path={toReactRouterPath(ROUTES.adminApiKeys)} element={<AdminApiKeysPage />} />
          </Route>
        </Route>
        <Route path={toReactRouterPath(ROUTES.authApiKey)} element={<ApiKeyLoginPage />} />
      </Routes>
    </BrowserRouter>
  )
}
