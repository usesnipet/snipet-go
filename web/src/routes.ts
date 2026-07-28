export const ROUTES = {

  home: "/",
  // auth routes
  authApiKey: "/auth/api-key",

  // admin routes
  admin: "/admin",
  adminKnowledgeIndex: "/admin/knowledge/index",
  adminKnowledgeSource: "/admin/knowledge/source",
  adminAgent: "/admin/agents",
  adminApiKeys: "/admin/api-keys",
  adminSettings: "/admin/settings",
  adminClients: "/admin/clients",

  // client routes
  client: "/client/{clientCode}",
  clientChat: "/client/{clientCode}/chat",
  clientChatSession: "/client/{clientCode}/chat/session/{sessionId}",
  clientChatLoginAnonymous: "/client/{clientCode}/chat/login-anonymous",
  clientUsers: "/client/{clientCode}/users",
  clientSessions: "/client/{clientCode}/sessions",
  clientSettings: "/client/{clientCode}/settings",
} as const;

export type RoutePath = (typeof ROUTES)[keyof typeof ROUTES];