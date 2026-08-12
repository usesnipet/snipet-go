export const ROUTES = {
  home: "/",
  knowledge: "/knowledge",
  knowledgeDetail: "/knowledge/{id}",
  agent: "/agent",
  llms: "/llms",
  apiKeys: "/api-keys",
  settings: "/settings",
  clients: "/clients",

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