export const ROUTES = {
  home: "/",

  apiKeys: "/api-keys",
  knowledge: "/knowledge",
  knowledgeDetail: "/knowledge/{id}",
  agent: "/agent",
  agentPlayground: "/agent/playground",
  llms: "/llms",
  apps: "/apps",

  // app routes
  app: "/apps/{appCode}",
  appUsers: "/apps/{appCode}/users",
  appSessions: "/apps/{appCode}/sessions",
  appSettings: "/apps/{appCode}/settings",
} as const;

export type RoutePath = (typeof ROUTES)[keyof typeof ROUTES];
