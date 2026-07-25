import type { HttpGetOptions, HttpPostOptions, HttpPutOptions, HttpDeleteOptions } from "./http";

type WithoutUrl<T> = Omit<T, "url">;

export type ServiceGetOptions<T = unknown> =
  WithoutUrl<HttpGetOptions<T>>;

export type ServicePostOptions<B = unknown, T = unknown> =
  WithoutUrl<HttpPostOptions<B, T>>;

export type ServicePutOptions<B = unknown, T = unknown> =
  WithoutUrl<HttpPutOptions<B, T>>;

export type ServiceDeleteOptions<T = unknown> =
  WithoutUrl<HttpDeleteOptions<T>>;