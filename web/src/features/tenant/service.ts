import { http } from "@/lib/http";

import {
  createTenantSchema, paginatedTenantSchema, tenantSchema, updateTenantSchema
} from "./schemas";

import type { CreateTenant, PaginatedTenant, Tenant, UpdateTenant } from "./schemas";
import type {
  ServiceDeleteOptions, ServiceGetOptions, ServicePostOptions, ServicePutOptions
} from "@/lib/services";

const TENANTS_URL = "/api/tenants";

const findMine = async (opts: ServiceGetOptions<PaginatedTenant>): Promise<PaginatedTenant> => {
  return http.get({
    url: `${TENANTS_URL}/me`,
    schemas: { response: paginatedTenantSchema },
    ...opts,
  });
};

const findByID = async (id: string, opts: ServiceGetOptions<Tenant>): Promise<Tenant> => {
  return http.get({
    url: `${TENANTS_URL}/{id}`,
    params: { id },
    schemas: { response: tenantSchema },
    ...opts,
  });
};

const findBySlug = async (slug: string, opts: ServiceGetOptions<Tenant>): Promise<Tenant> => {
  return http.get({
    url: `${TENANTS_URL}/slug/{slug}`,
    params: { slug },
    schemas: { response: tenantSchema },
    ...opts,
  });
};

const create = async (
  body: CreateTenant,
  opts: ServicePostOptions<CreateTenant, Tenant>,
): Promise<Tenant> => {
  return http.post({
    url: TENANTS_URL,
    body,
    schemas: {
      body: createTenantSchema,
      response: tenantSchema,
    },
    ...opts,
  });
};

const update = async (
  id: string,
  body: UpdateTenant,
  opts: ServicePutOptions<UpdateTenant, void>,
): Promise<void> => {
  return http.put({
    url: `${TENANTS_URL}/{id}`,
    params: { id },
    body,
    schemas: { body: updateTenantSchema },
    ...opts,
  });
};

const remove = async (id: string, opts: ServiceDeleteOptions<void>): Promise<void> => {
  return http.delete({
    url: `${TENANTS_URL}/{id}`,
    params: { id },
    ...opts,
  });
};

const leave = async (id: string, opts: ServicePostOptions<void, void>): Promise<void> => {
  return http.post({
    url: `${TENANTS_URL}/{id}/leave`,
    params: { id },
    ...opts,
  });
};

export const tenantService = {
  findMine,
  findByID,
  findBySlug,
  create,
  update,
  delete: remove,
  leave,
};
