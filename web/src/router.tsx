import { BrowserRouter, Route, Routes } from "react-router";

import { ApiKeyGuard } from "./features/api-key/components/api-key-guard";
import { ROUTES } from "./routes";
import { AdminApiKeysPage } from "./routes/admin/api-keys/page";
import { AdminClientsPage } from "./routes/admin/clients/page";
import { AdminLayout } from "./routes/admin/layout";
import { AdminPage } from "./routes/admin/page";
import { ApiKeyLoginPage } from "./routes/auth/api-key/page";
import { ClientLayout } from "./routes/client/layout";
import { ClientPage } from "./routes/client/page";

const toReactRouterPath = (path: (typeof ROUTES)[keyof typeof ROUTES]) => {
  return path.replaceAll(/{([^}]+)}/g, (_, p1) => `:${p1}`);
}

export const Router = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<ClientLayout />}>
          <Route path={toReactRouterPath(ROUTES.client)} element={<ClientPage />} />
        </Route>
        <Route element={<ApiKeyGuard />}>
          <Route element={<AdminLayout />}>
            <Route path={ROUTES.admin} element={<AdminPage />} />
            <Route path={ROUTES.adminClients} element={<AdminClientsPage />} />
            <Route path={ROUTES.adminApiKeys} element={<AdminApiKeysPage />} />
          </Route>
        </Route>
        <Route path={ROUTES.authApiKey} element={<ApiKeyLoginPage />} />
      </Routes>
    </BrowserRouter>
  )
}
