<script lang="ts" module>
  import type { RouteId } from "$app/types";

  type NavItem = {
    title: string;
    url: RouteId;
  };

  type NavGroup = {
    title: string;
    items: NavItem[];
  };

  const items: Array<NavItem | NavGroup> = [
    {
      title: "Home",
      url: "/(protected)",
    },
    {
      title: "Knowledge",
      url: "/(protected)/knowledge",
    },
    {
      title: "API Key",
      url: "/(protected)/api-key",
    },
  ];

  function isNavGroup(item: NavItem | NavGroup): item is NavGroup {
    return "items" in item;
  }
</script>

<script lang="ts">
  import { page } from "$app/state";
  import { resolve } from "$app/paths";

  import * as Collapsible from "$lib/components/ui/collapsible/index.js";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import CommandIcon from "@lucide/svelte/icons/command";
  import MinusIcon from "@lucide/svelte/icons/minus";
  import PlusIcon from "@lucide/svelte/icons/plus";
  import type { ComponentProps } from "svelte";
  import ThemeToggle from "./theme-toggle.svelte";
	import Button from "./ui/button/button.svelte";
	import LogOutIcon from "@lucide/svelte/icons/log-out";
	import { logout } from "$lib/features/api-key/stores/api-key-auth.svelte";

  let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> =
    $props();

  let collapsibleOpen = $state<Record<string, boolean>>({});

  function isActive(url: RouteId) {
    if (url === "/(protected)/knowledge") {
      return page.route.id?.startsWith("/(protected)/knowledge") ?? false;
    }
    return page.route.id === url;
  }

  function isGroupActive(group: NavGroup) {
    return group.items.some((subItem) => isActive(subItem.url));
  }

  function isGroupOpen(title: string, group: NavGroup) {
    return collapsibleOpen[title] ?? isGroupActive(group);
  }
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
  <Sidebar.Header>
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
  </Sidebar.Header>
  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.Menu>
        {#each items as item (item.title)}
          {#if isNavGroup(item)}
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
                        <Sidebar.MenuSubButton isActive={isActive(subItem.url)}>
                          {#snippet child({ props })}
                            <a href={resolve(subItem.url)} {...props}>
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
              <Sidebar.MenuButton isActive={isActive(item.url)}>
                {#snippet child({ props })}
                  <a href={resolve(item.url)} {...props}>{item.title}</a>
                {/snippet}
              </Sidebar.MenuButton>
            </Sidebar.MenuItem>
          {/if}
        {/each}
      </Sidebar.Menu>
    </Sidebar.Group>
  </Sidebar.Content>
  <Sidebar.Footer>
    <Button variant="secondary" onclick={() => logout()}>
      <LogOutIcon />
      Logout
    </Button>
  </Sidebar.Footer>
  <Sidebar.Rail />
</Sidebar.Root>
