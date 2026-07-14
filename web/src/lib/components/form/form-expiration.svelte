<script lang="ts" generics="T extends Record<string, unknown>">
	import * as Popover from "$lib/components/ui/popover/index.js";
	import { Calendar } from "$lib/components/ui/calendar/index.js";
	import { Button, buttonVariants } from "$lib/components/ui/button/index.js";
	import {
		Field,
		FieldLabel,
		FieldError,
		FieldDescription,
	} from "$lib/components/ui/field/index.js";
	import { cn } from "$lib/utils.js";
	import { formFieldProxy, type FormPathLeaves, type SuperForm } from "sveltekit-superforms";
	import ChevronDownIcon from "@lucide/svelte/icons/chevron-down";
	import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
	import {
		CalendarDate,
		getLocalTimeZone,
		today,
		type DateValue,
	} from "@internationalized/date";
	import {
		CUSTOM_EXPIRATION_VALUE,
		DEFAULT_EXPIRATION_OPTIONS,
		resolveExpirationExpression,
		toExpirationFormValue,
		type ExpirationOption,
	} from "./expiration";

	type Props = {
		form: SuperForm<T>;
		field: FormPathLeaves<T>;
		label?: string;
		description?: string;
		placeholder?: string;
		options?: ExpirationOption[];
		allowCustom?: boolean;
		customLabel?: string;
	};

	let {
		form,
		field,
		label,
		description,
		placeholder = "Select expiration...",
		options = DEFAULT_EXPIRATION_OPTIONS,
		allowCustom = false,
		customLabel = "Custom date",
	}: Props = $props();

	const { value, constraints, errors } = $derived.by(() => formFieldProxy(form, field));
	const fieldErrors = $derived.by(() => {
		const list = $errors;
		if (!Array.isArray(list)) return undefined;
		return list.map((message) => ({ message }));
	});

	let open = $state(false);
	/** Which panel is visible inside the single popover. */
	let panel = $state<"options" | "calendar">("options");
	let selectValue = $state<string>("");
	let calendarValue = $state<DateValue | undefined>(undefined);

	const selectOptions = $derived.by(() => {
		const items = options.map((option) => ({
			label: option.label,
			value: option.expression,
		}));
		if (allowCustom) {
			items.push({ label: customLabel, value: CUSTOM_EXPIRATION_VALUE });
		}
		return items;
	});

	const isCustom = $derived(selectValue === CUSTOM_EXPIRATION_VALUE);

	const triggerContent = $derived.by(() => {
		if (isCustom && calendarValue) {
			return calendarValue.toDate(getLocalTimeZone()).toLocaleDateString(undefined, {
				day: "2-digit",
				month: "short",
				year: "numeric",
			});
		}
		return selectOptions.find((option) => option.value === selectValue)?.label ?? placeholder;
	});

	function dateToCalendarDate(date: Date): CalendarDate {
		return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
	}

	function applyExpression(expression: string) {
		const resolved = resolveExpirationExpression(expression);
		$value = toExpirationFormValue(resolved) as typeof $value;
		calendarValue = resolved ? dateToCalendarDate(resolved) : undefined;
	}

	function applyCustomDate(date: DateValue | undefined) {
		if (!date) {
			$value = "" as typeof $value;
			return;
		}
		$value = toExpirationFormValue(date.toDate(getLocalTimeZone())) as typeof $value;
	}

	function handleOptionSelect(next: string) {
		if (next === CUSTOM_EXPIRATION_VALUE) {
			selectValue = CUSTOM_EXPIRATION_VALUE;
			if (!calendarValue) {
				calendarValue = today(getLocalTimeZone());
			}
			panel = "calendar";
			return;
		}
		selectValue = next;
		applyExpression(next);
		open = false;
		panel = "options";
	}

	function handleCalendarSelect(next: DateValue | undefined) {
		selectValue = CUSTOM_EXPIRATION_VALUE;
		calendarValue = next;
		applyCustomDate(next);
		open = false;
		panel = "options";
	}

	function handleOpenChange(next: boolean) {
		open = next;
		if (next) {
			// Always open on the options list; calendar appears after choosing Custom.
			panel = "options";
		}
	}

	$effect(() => {
		const current = String($value ?? "");

		if (!current) {
			const neverOption = options.find(
				(option) => option.expression.trim().toLowerCase() === "never",
			);
			selectValue = neverOption?.expression ?? options[0]?.expression ?? "";
			calendarValue = undefined;
			return;
		}

		if (selectValue && selectValue !== CUSTOM_EXPIRATION_VALUE) return;

		const parsed = new Date(current);
		if (Number.isNaN(parsed.getTime())) return;

		if (allowCustom) {
			selectValue = CUSTOM_EXPIRATION_VALUE;
			calendarValue = dateToCalendarDate(parsed);
		}
	});
</script>

<Field>
	<div class="flex items-center">
		<FieldLabel for={field}>{label}</FieldLabel>
	</div>

	<Popover.Root bind:open onOpenChange={handleOpenChange}>
		<Popover.Trigger
			id={field}
			class={cn(buttonVariants({ variant: "outline" }), "w-full justify-between font-normal")}
			{...$constraints}
		>
			<span class="truncate">{triggerContent}</span>
			<ChevronDownIcon class="text-muted-foreground" />
		</Popover.Trigger>
		<Popover.Content
			class={cn("p-0", panel === "calendar" ? "w-auto" : "w-(--bits-popover-anchor-width)")}
			align="start"
		>
			{#if panel === "calendar"}
				<div class="flex flex-col gap-1 p-1">
					<Button
						variant="ghost"
						size="sm"
						class="justify-start"
						onclick={() => (panel = "options")}
					>
						<ChevronLeftIcon data-icon="inline-start" />
						Back
					</Button>
					<Calendar
						type="single"
						bind:value={calendarValue}
						captionLayout="dropdown"
						minValue={today(getLocalTimeZone())}
						onValueChange={handleCalendarSelect}
					/>
				</div>
			{:else}
				<div class="flex flex-col gap-0.5 p-1" role="listbox">
					{#each selectOptions as option (option.value)}
						<Button
							variant="ghost"
							class={cn(
								"h-8 w-full justify-start font-normal",
								selectValue === option.value && "bg-accent",
							)}
							onclick={() => handleOptionSelect(option.value)}
						>
							{option.label}
						</Button>
					{/each}
				</div>
			{/if}
		</Popover.Content>
	</Popover.Root>

	<FieldError errors={fieldErrors} />
	{#if description}
		<FieldDescription>{description}</FieldDescription>
	{/if}
</Field>
