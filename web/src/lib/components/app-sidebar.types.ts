import { resolve } from "$app/paths";

import type { RouteId, RouteParams } from "$app/types";

export type SidebarNavLink = {
	title: string;
	route: RouteId;
	exact?: boolean;
	params?: RouteParams<RouteId>;
};

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

export function resolveNavLink(link: SidebarNavLink): string {
	if ("params" in link && link.params) {
		return resolve(link.route, link.params);
	}

	return resolve(link.route);
}
