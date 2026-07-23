import type { SidebarNavSection } from "$lib/components/app-sidebar.types";

export const systemSidebarSections: SidebarNavSection[] = [
	{
		title: "System",
		items: [
			{ title: "Home", route: "/manage/(protected)/(system)", exact: true  },
			{ title: "Agents", route: "/manage/(protected)/(system)/agent" },
			{ title: "Clients", route: "/manage/(protected)/(system)/client" },
			{ title: "Knowledge", route: "/manage/(protected)/(system)/knowledge" },
			{ title: "API Key", route: "/manage/(protected)/(system)/api-key" },
		],
	},
];
