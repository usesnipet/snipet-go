<script lang="ts" module>
  import type { RouteId } from "$app/types";

  type NavItem = {
    title: string;
    url: RouteId;
    items?: NavItem[];
    isActive?: boolean;
  }

  const data: Record<string, NavItem[]> = {
    "Main": [
      {
        title: "Getting Started",
        url: "/",
        items: [
          {
            title: "Installation",
            url: "/",
          },
          {
            title: "Project Structure",
            url: "/",
          },
        ],
      },
    ],
  };
</script>
<script lang="ts">
	import { resolve } from "$app/paths";

  import * as Collapsible from "$lib/components/ui/collapsible/index.js";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import CommandIcon from "@lucide/svelte/icons/command";
  import MinusIcon from "@lucide/svelte/icons/minus";
  import PlusIcon from "@lucide/svelte/icons/plus";
  import type { ComponentProps } from "svelte";
  let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();
</script>
<Sidebar.Root bind:ref variant="inset" {...restProps}>
  <Sidebar.Header>
    <Sidebar.Menu>
      <Sidebar.MenuItem>
        <Sidebar.MenuButton size="lg">
          {#snippet child({ props })}
            <a href="##" {...props}>
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
          {/snippet}
        </Sidebar.MenuButton>
      </Sidebar.MenuItem>
    </Sidebar.Menu>
  </Sidebar.Header>
  <Sidebar.Content>
    {#each Object.entries(data) as [key, items] (key)}
    <Sidebar.Group>
      <Sidebar.GroupLabel>{key}</Sidebar.GroupLabel>
      {#each items as item, index (item.title)}
        <Sidebar.Menu>
          <Collapsible.Root open={index === 1} class="group/collapsible">
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
              {#if item.items?.length}
                <Collapsible.Content>
                  <Sidebar.MenuSub>
                    {#each item.items as subItem (subItem.title)}
                      <Sidebar.MenuSubItem>
                        <Sidebar.MenuSubButton isActive={subItem.isActive}>
                          {#snippet child({ props })}
                            <a href={resolve(subItem.url)} {...props}
                              >{subItem.title}</a
                            >
                          {/snippet}
                        </Sidebar.MenuSubButton>
                      </Sidebar.MenuSubItem>
                    {/each}
                  </Sidebar.MenuSub>
                </Collapsible.Content>
              {/if}
            </Sidebar.MenuItem>
          </Collapsible.Root>
        </Sidebar.Menu>
      {/each}
      </Sidebar.Group>
      {/each}
  </Sidebar.Content>
  <Sidebar.Rail />
</Sidebar.Root>