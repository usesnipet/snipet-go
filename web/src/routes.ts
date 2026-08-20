export const ROUTES = {
  home: "/",

  // tenant management routes
  members: "/{tenantSlug}/members",
  inviteMember: "/{tenantSlug}/members/invite",
  apiKeys: "/{tenantSlug}/api-keys",
  settings: "/{tenantSlug}/settings",

  selectTenant: "/tenant-select",
  acceptInvite: "/invite",
  tenantHome: "/{tenantSlug}",
  knowledge: "/{tenantSlug}/knowledge",
  knowledgeDetail: "/{tenantSlug}/knowledge/{id}",
  agent: "/{tenantSlug}/agent",
  llms: "/{tenantSlug}/llms",
  apps: "/{tenantSlug}/apps",

  // auth routes
  authLogin: "/auth/login",
  authRegister: "/auth/register",
  authForgotPassword: "/auth/forgot-password",
  authResetPassword: "/auth/reset-password",
  authActivate: "/auth/activate",
  authResendActivation: "/auth/resend-activation",

  // app routes
  app: "/apps/{appCode}",
  appUsers: "/apps/{appCode}/users",
  appSessions: "/apps/{appCode}/sessions",
  appSettings: "/apps/{appCode}/settings",
} as const;

export type RoutePath = (typeof ROUTES)[keyof typeof ROUTES];