export const ROUTES = {
  home: "/",

  selectTenant: "/tenant-select",
  tenantHome: "/{tenantSlug}",
  knowledge: "/{tenantSlug}/knowledge",
  knowledgeDetail: "/{tenantSlug}/knowledge/{id}",
  agent: "/{tenantSlug}/agent",
  llms: "/{tenantSlug}/llms",
  apiKeys: "/{tenantSlug}/api-keys",
  settings: "/{tenantSlug}/settings",
  clients: "/{tenantSlug}/clients",

  // auth routes
  authLogin: "/auth/login",
  authRegister: "/auth/register",
  authForgotPassword: "/auth/forgot-password",
  authResetPassword: "/auth/reset-password",
  authActivate: "/auth/activate",
  authResendActivation: "/auth/resend-activation",

  // client routes
  client: "/client/{clientCode}",
  clientUsers: "/client/{clientCode}/users",
  clientSessions: "/client/{clientCode}/sessions",
  clientSettings: "/client/{clientCode}/settings",
} as const;

export type RoutePath = (typeof ROUTES)[keyof typeof ROUTES];