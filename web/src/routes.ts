export const ROUTES = {
  home: "/",

  // tenant management routes
  members: "/{tenantSlug}/members",
  inviteMember: "/{tenantSlug}/members/invite",
  apiKeys: "/{tenantSlug}/api-keys",
  settings: "/{tenantSlug}/settings",

  selectTenant: "/tenant-select",
  tenantHome: "/{tenantSlug}",
  knowledge: "/{tenantSlug}/knowledge",
  knowledgeDetail: "/{tenantSlug}/knowledge/{id}",
  agent: "/{tenantSlug}/agent",
  llms: "/{tenantSlug}/llms",
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