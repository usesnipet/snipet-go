<script lang="ts">
	import { page } from "$app/state";
	import { resolve } from "$app/paths";

	import * as Collapsible from "$lib/components/ui/collapsible/index.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import type {
		SidebarNavGroup,
		SidebarNavLink,
		SidebarNavSection,
	} from "$lib/components/app-sidebar.types.js";
	import { isSidebarNavGroup, resolveNavLink } from "$lib/components/app-sidebar.types.js";
	import CommandIcon from "@lucide/svelte/icons/command";
	import MinusIcon from "@lucide/svelte/icons/minus";
	import PlusIcon from "@lucide/svelte/icons/plus";
	import LogOutIcon from "@lucide/svelte/icons/log-out";
	import type { ComponentProps, Snippet } from "svelte";
	import ThemeToggle from "./theme-toggle.svelte";
	import Button from "./ui/button/button.svelte";
	import { logout } from "$lib/features/api-key/stores/api-key-auth.svelte";

	type Props = ComponentProps<typeof Sidebar.Root> & {
		sections: SidebarNavSection[];
		header?: Snippet;
		footer?: Snippet;
	};

	let {
		ref = $bindable(null),
		sections,
		header,
		footer,
		...restProps
	}: Props = $props();

	let collapsibleOpen = $state<Record<string, boolean>>({});

	function isActive(link: SidebarNavLink) {
		const href = resolveNavLink(link);
		const { pathname } = page.url;

		if (link.exact) {
			return pathname === href;
		}

		return pathname === href || pathname.startsWith(`${href}/`);
	}

	function isGroupActive(group: SidebarNavGroup) {
		return group.items.some((subItem) => isActive(subItem));
	}

	function isGroupOpen(title: string, group: SidebarNavGroup) {
		return collapsibleOpen[title] ?? isGroupActive(group);
	}
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		{#if header}
			{@render header()}
		{:else}
			<Sidebar.Menu>
				<Sidebar.MenuItem>
					<Sidebar.MenuButton size="lg">
						{#snippet child({ props })}
							<div class="flex items-center gap-2">
								<a href={resolve("/")} {...props}>
									<div
										class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
									>
										<CommandIcon class="size-4" />
									</div>
									<div class="grid flex-1 text-start text-sm leading-tight">
										<span class="truncate font-medium">Snipet</span>
										<span class="truncate text-xs">v1.0.0</span>
									</div>
								</a>
								<ThemeToggle />
							</div>
						{/snippet}
					</Sidebar.MenuButton>
				</Sidebar.MenuItem>
			</Sidebar.Menu>
		{/if}
	</Sidebar.Header>
	<Sidebar.Content>
		{#each sections as section (section.title)}
			<Sidebar.Group>
				{#if section.title}
					<Sidebar.GroupLabel>{section.title}</Sidebar.GroupLabel>
				{/if}
				<Sidebar.Menu>
					{#each section.items as item (item.title)}
						{#if isSidebarNavGroup(item)}
							<Collapsible.Root
								open={isGroupOpen(item.title, item)}
								onOpenChange={(open) => (collapsibleOpen[item.title] = open)}
								class="group/collapsible"
							>
								<Sidebar.MenuItem>
									<Collapsible.Trigger>
										{#snippet child({ props })}
											<Sidebar.MenuButton {...props}>
												{item.title}
												<PlusIcon
													class="ms-auto group-data-[state=open]/collapsible:hidden"
												/>
												<MinusIcon
													class="ms-auto group-data-[state=closed]/collapsible:hidden"
												/>
											</Sidebar.MenuButton>
										{/snippet}
									</Collapsible.Trigger>
									<Collapsible.Content>
										<Sidebar.MenuSub>
											{#each item.items as subItem (subItem.title)}
												<Sidebar.MenuSubItem>
													<Sidebar.MenuSubButton isActive={isActive(subItem)}>
														{#snippet child({ props })}
															<a href={resolveNavLink(subItem)} {...props}>
																{subItem.title}
															</a>
														{/snippet}
													</Sidebar.MenuSubButton>
												</Sidebar.MenuSubItem>
											{/each}
										</Sidebar.MenuSub>
									</Collapsible.Content>
								</Sidebar.MenuItem>
							</Collapsible.Root>
						{:else}
							<Sidebar.MenuItem>
								<Sidebar.MenuButton isActive={isActive(item)}>
									{#snippet child({ props })}
										<a href={resolveNavLink(item)} {...props}>{item.title}</a>
									{/snippet}
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/if}
					{/each}
				</Sidebar.Menu>
			</Sidebar.Group>
		{/each}
	</Sidebar.Content>
	<Sidebar.Footer>
		{#if footer}
			{@render footer()}
		{:else}
			<Button variant="secondary" onclick={() => logout()}>
				<LogOutIcon />
				Logout
			</Button>
		{/if}
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>
