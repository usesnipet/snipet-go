import type { SidebarNavSection } from "$lib/components/app-sidebar.types";

export function createClientSidebarSections(clientCode: string): SidebarNavSection[] {
	return [
		{
			title: "Overview",
			items: [
				{
					title: "Home",
					route: "/manage/(protected)/(client)/c/[code]",
					params: { code: clientCode },
					exact: true,
				},
				{
					title: "Sessions",
					route: "/manage/(protected)/(client)/c/[code]/session",
					params: { code: clientCode },
				},
			],
		},
	];
}
