import { Button } from "$lib/components/ui/button/index.js";
import {
  Field, FieldDescription, FieldError, FieldGroup, FieldLabel, FieldLegend, FieldSet, FieldTitle
} from "$lib/components/ui/field/index.js";
import { Input } from "$lib/components/ui/input/index.js";
import { Textarea } from "$lib/components/ui/textarea/index.js";
import { setThemeContext } from "@sjsf/shadcn4-theme";
import {
  ButtonGroup, Checkbox, Select, SelectContent, SelectItem, SelectTrigger
} from "@sjsf/shadcn4-theme/new-york";

// https://x0k.dev/svelte-jsonschema-form/themes/shadcn4/#components
export function setShadcnThemeContext() {
	setThemeContext({
		components: {
			Button,
			FieldDescription,
			FieldError,
			FieldLabel,
			Field,
			FieldSet,
			FieldLegend,
			FieldGroup,
			FieldTitle,
			ButtonGroup,
			Checkbox,
			Input,
			Select,
			SelectTrigger,
			SelectContent,
			SelectItem,
			Textarea

			// Popover,
			// PopoverContent,
			// PopoverTrigger,
			// Command,
			// CommandInput,
			// CommandList,
			// CommandEmpty,
			// CommandGroup,
			// CommandItem,
			// Calendar,
			// RangeCalendar,
			// ToggleGroup,
			// ToggleGroupItem,
			// RadioGroup,
			// RadioGroupItem,
			// Slider,
			// Switch
		}
	});
}
