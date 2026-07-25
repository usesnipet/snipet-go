import { BrowserRouter, Route, Routes } from "react-router";

import { ApiKeyGuard } from "./features/api-key/components/api-key-guard";
import { ROUTES } from "./routes";
import { AdminLayout } from "./routes/admin/layout";
import { AdminPage } from "./routes/admin/page";
import { ApiKeyLoginPage } from "./routes/auth/api-key/page";
import { HomePage } from "./routes/page";

export const Router = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path={ROUTES.home} element={<HomePage />} />
        <Route element={<ApiKeyGuard />}>
          <Route element={<AdminLayout />}>
            <Route path={ROUTES.admin} element={<AdminPage />} />
          </Route>
        </Route>
        <Route path={ROUTES.authApiKey} element={<ApiKeyLoginPage />} />
      </Routes>
    </BrowserRouter>
  )
}
