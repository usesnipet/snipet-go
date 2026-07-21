import type { SidebarNavSection } from "$lib/components/app-sidebar.types";

export const systemSidebarSections: SidebarNavSection[] = [
	{
		title: "System",
		items: [
			{ title: "Home", route: "/(protected)/(system)", exact: true  },
			{ title: "Agents", route: "/(protected)/(system)/agent" },
			{ title: "Clients", route: "/(protected)/(system)/client" },
			{ title: "Knowledge", route: "/(protected)/(system)/knowledge" },
			{ title: "API Key", route: "/(protected)/(system)/api-key" },
		],
	},
];
