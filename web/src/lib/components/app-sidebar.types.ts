import { resolve } from "$app/paths";

import type { RouteId, RouteParams } from "$app/types";

export type SidebarNavLink = {
	[R in RouteId]: {
		title: string;
		route: R;
		exact?: boolean;
	} & (RouteParams<R> extends Record<string, never>
		? { params?: undefined }
		: { params: RouteParams<R> });
}[RouteId];

export type SidebarNavGroup = {
	title: string;
	items: SidebarNavLink[];
};

export type SidebarNavEntry = SidebarNavLink | SidebarNavGroup;

export type SidebarNavSection = {
	title: string;
	items: SidebarNavEntry[];
};

export function isSidebarNavGroup(entry: SidebarNavEntry): entry is SidebarNavGroup {
	return "items" in entry;
}

/** `resolve` overloads don't accept a `RouteId` union; widen for dynamic nav links. */
const resolveRoute = resolve as (route: RouteId, params?: Record<string, string>) => string;

export function resolveNavLink(link: SidebarNavLink): string {
	if (link.params) {
		return resolveRoute(link.route, link.params);
	}

	return resolveRoute(link.route);
}
